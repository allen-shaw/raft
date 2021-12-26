package options

import (
	"github.com/AllenShaw19/raft"
	"github.com/AllenShaw19/raft/conf"
)

type NodeOptions struct {
	ElectionTimeoutMs      int64
	ElectionPriority       int
	DecayPriorityGap       int
	LeaderLeaseTimeRatio   int
	SnapshotIntervalSecs   int64
	SnapshotLogIndexMargin int
	CatchupMargin          int

	InitialConf            *conf.Configuration
	Fsm                    raft.StateMachine
	LogUri                 string
	RaftMetaUri            string
	SnapshotUri            string
	FilterBeforeCopyRemote bool
	DisableCli             bool
	SharedTimerPool        bool
	TimerPoolSize          int
	EnableMetrics          bool
}

func DefaultNodeOptions() *NodeOptions {

}
