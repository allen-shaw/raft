package option

import "github.com/AllenShaw19/raft/raft/conf"

type NodeOptions struct {
	RpcOptions
	raftOptions *RaftOptions

	enableMetrics        bool
	electionPriority     int
	snapshotIntervalSecs int
	initialConf          *conf.Configuration
}

func (n *NodeOptions) ElectionPriority() int {
	return n.electionPriority
}

func (n *NodeOptions) SetElectionPriority(electionPriority int) {
	n.electionPriority = electionPriority
}

func (n *NodeOptions) IsEnableMetrics() bool {
	return n.enableMetrics
}

func (n *NodeOptions) SetEnableMetrics(enableMetrics bool) {
	n.enableMetrics = enableMetrics
}

func (n *NodeOptions) RaftOptions() *RaftOptions {
	return n.raftOptions
}

func (n *NodeOptions) InitialConf() *conf.Configuration {
	return n.initialConf
}

func (n *NodeOptions) SetInitialConf(initialConf *conf.Configuration) {
	n.initialConf = initialConf
}

func (n *NodeOptions) SnapshotIntervalSecs() int {
	return n.snapshotIntervalSecs
}

func (n *NodeOptions) SetSnapshotIntervalSecs(snapshotIntervalSecs int) {
	n.snapshotIntervalSecs = snapshotIntervalSecs
}
