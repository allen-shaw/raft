package core

import (
	"github.com/AllenShaw19/raft/log"
	"github.com/AllenShaw19/raft/raft"
	"github.com/AllenShaw19/raft/raft/closure"
	"github.com/AllenShaw19/raft/raft/entity"
	"github.com/AllenShaw19/raft/raft/option"
	"github.com/AllenShaw19/raft/raft/util"
	"sync"
	"sync/atomic"
	_error "github.com/AllenShaw19/raft/raft/error"
)

type BallotBox struct {
	waiter             raft.FSMCaller
	closureQueue       closure.ClosureQueue
	lock               sync.RWMutex
	lastCommittedIndex int64 // atomic
	pendingIndex       int64
	pendingMetaQueue   util.SegmentList[*entity.Ballot]
}

func NewBallotBox() *BallotBox {

}
func (b *BallotBox) Init(opts *option.BallotBoxOptions) bool {
	if opts.Waiter() == nil || opts.ClosureQueue() == nil {
		log.Error("waiter or closure queue is null.")
		return false
	}
	b.waiter = opts.Waiter()
	b.closureQueue = opts.ClosureQueue()
	return true
}

func (b *BallotBox) LastCommittedIndex() int64 {
	return atomic.LoadInt64(&b.lastCommittedIndex)
}

// CommitAt Called by leader, otherwise the behavior is undefined
// Set logs in [first_log_index, last_log_index] are stable at |peer|.
func (b *BallotBox) CommitAt(firstLogIndex int64, lastLogIndex int64, peer *entity.PeerId) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.pendingIndex == 0 {
		return _error.ErrInval
	}
	if lastLogIndex < b.pendingIndex {
		return nil
	}
	if lastLogIndex >= b.pendingIndex + b.pendingMetaQueue
}
