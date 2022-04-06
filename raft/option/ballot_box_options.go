package option

import (
	"github.com/AllenShaw19/raft/raft"
	"github.com/AllenShaw19/raft/raft/closure"
)

type BallotBoxOptions struct {
	waiter       raft.FSMCaller
	closureQueue closure.ClosureQueue
}

func (b *BallotBoxOptions) Waiter() raft.FSMCaller {
	return b.waiter
}

func (b *BallotBoxOptions) SetWaiter(waiter raft.FSMCaller) {
	b.waiter = waiter
}

func (b *BallotBoxOptions) ClosureQueue() closure.ClosureQueue {
	return b.closureQueue
}

func (b *BallotBoxOptions) SetClosureQueue(closureQueue closure.ClosureQueue) {
	b.closureQueue = closureQueue
}
