package raft

import (
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

func (b *BallotBox) CommitAt(firstLogIndex, lastLogIndex int64, peer *PeerId) {

}

func (b *BallotBox) ClearPendingTasks() {

}

func (b *BallotBox) ResetPendingIndex(newPendingIndex int64) {

}

func (b *BallotBox) AppendPendingTask(conf, oldConf *Configuration, closure Closure) {

}

func (b *BallotBox) SetLastCommittedIndex(lastCommittedIndex int64) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

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
