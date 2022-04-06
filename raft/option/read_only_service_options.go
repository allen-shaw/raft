package option

import (
	"github.com/AllenShaw19/raft/raft"
)

type ReadOnlyServiceOptions struct {
	raftOptions *RaftOptions
	node        raft.Node
	fsmCaller   raft.FSMCaller
}

func (r *ReadOnlyServiceOptions) RaftOptions() *RaftOptions {
	return r.raftOptions
}

func (r *ReadOnlyServiceOptions) SetRaftOptions(raftOptions *RaftOptions) {
	r.raftOptions = raftOptions
}

func (r *ReadOnlyServiceOptions) Node() raft.Node {
	return r.node
}

func (r *ReadOnlyServiceOptions) SetNode(node raft.Node) {
	r.node = node
}

func (r *ReadOnlyServiceOptions) FsmCaller() raft.FSMCaller {
	return r.fsmCaller
}

func (r *ReadOnlyServiceOptions) SetFsmCaller(fsmCaller raft.FSMCaller) {
	r.fsmCaller = fsmCaller
}
