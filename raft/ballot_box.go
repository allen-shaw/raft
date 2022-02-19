package raft

import (
	"errors"
	"fmt"
	"github.com/AllenShaw19/raft/log"
	"sync"
	"sync/atomic"
)

type BallotBoxOptions struct {
	Waiter       *FSMCaller
	ClosureQueue *ClosureQueue
}

type BallotBoxStatus struct {
	CommittedIndex   int64
	PendingIndex     int64
	PendingQueueSize int64
}

type BallotBox struct {
	waiter             *FSMCaller
	closureQueue       *ClosureQueue
	mutex              sync.Mutex
	lastCommittedIndex int64 // atomic
	pendingIndex       int64
	pendingMetaQueue   []Ballot
}

func (b *BallotBox) Init(options *BallotBoxOptions) error {
	if options.Waiter == nil || options.ClosureQueue == nil {
		log.Error("waiter or closure queue is null")
		return errInvalid
	}
	b.waiter = options.Waiter
	b.closureQueue = options.ClosureQueue
	return nil
}

func (b *BallotBox) CommitAt(firstLogIndex, lastLogIndex int64, peer *PeerId) error {
	b.mutex.Lock()
	if b.pendingIndex == 0 {
		b.mutex.Unlock()
		return errInvalid
	}
	if lastLogIndex < b.pendingIndex {
		b.mutex.Unlock()
		return nil
	}
	if lastLogIndex >= b.pendingIndex+int64(len(b.pendingMetaQueue)) {
		b.mutex.Unlock()
		return errRange
	}

	var lastCommittedIndex int64
	startAt := maxInt64(b.pendingIndex, firstLogIndex)
	var posHint PosHint
	for logIndex := startAt; logIndex <= lastLogIndex; logIndex++ {
		bl := b.pendingMetaQueue[logIndex-b.pendingIndex]
		posHint = bl.GrantWithPosHint(peer, posHint)
		if bl.Granted() {
			lastCommittedIndex = logIndex
		}
	}

	if lastCommittedIndex == 0 {
		b.mutex.Unlock()
		return nil
	}

	popIndex := lastCommittedIndex - b.pendingIndex
	b.pendingMetaQueue = b.pendingMetaQueue[popIndex:]
	b.pendingIndex = lastCommittedIndex + 1
	atomic.StoreInt64(&b.lastCommittedIndex, lastCommittedIndex)
	b.mutex.Unlock()

	b.waiter.OnCommitted(lastCommittedIndex)
	return nil
}

func (b *BallotBox) ClearPendingTasks() error {
	b.mutex.Lock()
	b.pendingMetaQueue = make([]Ballot, 0)
	b.pendingIndex = 0
	b.mutex.Unlock()
	b.closureQueue.Clear()
	return nil
}

func (b *BallotBox) ResetPendingIndex(newPendingIndex int64) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if !(b.pendingIndex == 0 && len(b.pendingMetaQueue) == 0) {
		log.Error("pendingIndex=%d, pendingMetaQueueSize=%d", b.pendingIndex, len(b.pendingMetaQueue))
		return errors.New("ballot box has pending ballot")
	}
	if newPendingIndex <= atomic.LoadInt64(&b.lastCommittedIndex) {
		log.Error("newPendingIndex=%d, lastCommittedIndex=%d", newPendingIndex, b.lastCommittedIndex)
		return errors.New("invalid new pending index")
	}

	b.pendingIndex = newPendingIndex
	b.closureQueue.ResetFirstIndex(newPendingIndex)
	return nil
}

func (b *BallotBox) AppendPendingTask(conf, oldConf *Configuration, closure Closure) error {
	bl := Ballot{}
	if err := bl.Init(conf, oldConf); err != nil {
		log.Error("Fail to init ballot, err: %v", err)
		return err
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.pendingMetaQueue = append(b.pendingMetaQueue, bl)
	b.closureQueue.AppendPendingClosure(closure)
	return nil
}

func (b *BallotBox) SetLastCommittedIndex(lastCommittedIndex int64) error {
	if b.pendingIndex != 0 || len(b.pendingMetaQueue) != 0 {
		if lastCommittedIndex >= b.pendingIndex {
			log.Error("node changes to leader, pending_index=%d, parameter last_committed_index=%d", b.pendingIndex, lastCommittedIndex)
		}
		return errors.New("ballot has pending index")
	}

	if lastCommittedIndex < atomic.LoadInt64(&b.lastCommittedIndex) {
		return errInvalid
	}

	if lastCommittedIndex > atomic.LoadInt64(&b.lastCommittedIndex) {
		atomic.StoreInt64(&b.lastCommittedIndex, lastCommittedIndex)
		b.waiter.OnCommitted(lastCommittedIndex)
	}

	return nil
}

func (b *BallotBox) LastCommittedIndex() int64 {
	return atomic.LoadInt64(&b.lastCommittedIndex)
}

func (b *BallotBox) Describe() string {
	b.mutex.Lock()
	committedIndex := b.lastCommittedIndex
	pendingIndex := int64(0)
	pendingQueueSize := 0
	if b.pendingIndex != 0 {
		pendingIndex = b.pendingIndex
		pendingQueueSize = len(b.pendingMetaQueue)
	}
	defer b.mutex.Unlock()

	s := fmt.Sprintf("last_committed_index: %d\n", committedIndex)
	if pendingQueueSize != 0 {
		s += fmt.Sprintf("pending_index: %d\n", pendingIndex)
		s += fmt.Sprintf("pending_queue_size: %d\n", pendingQueueSize)
	}
	return s
}

func (b *BallotBox) GetStatus() BallotBoxStatus {
	status := BallotBoxStatus{}

	b.mutex.Lock()
	defer b.mutex.Unlock()
	status.CommittedIndex = b.lastCommittedIndex
	if len(b.pendingMetaQueue) != 0 {
		status.PendingIndex = b.pendingIndex
		status.PendingQueueSize = int64(len(b.pendingMetaQueue))
	}
	return status
}
