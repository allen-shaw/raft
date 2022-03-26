package core

import (
	"github.com/AllenShaw19/raft/raft"
	"github.com/AllenShaw19/raft/raft/closure"
	"github.com/AllenShaw19/raft/raft/conf"
	"github.com/AllenShaw19/raft/raft/entity"
	"github.com/AllenShaw19/raft/raft/option"
	"github.com/AllenShaw19/raft/raft/rpc"
	"github.com/AllenShaw19/raft/raft/storage"
	"github.com/AllenShaw19/raft/raft/util"
	disruptor "github.com/smartystreets-prototypes/go-disruptor"
	"sync"
)

const (
	GlobalNumNodes = 0 // atomic
)

type Node struct {
	// Internal states
	mutex               nodeReadWriteLock
	state               State
	wg                  sync.WaitGroup
	currTerm            int64
	lastLeaderTimestamp int64
	leaderId            *entity.PeerId
	votedId             *entity.PeerId
	voteCtx             *entity.Ballot
	prevVoteCtx         *entity.Ballot
	conf                *conf.ConfigurationEntry
	stopTransferArg     *stopTransferArg

	// Raft group and node options and identifier
	groupId     string
	options     *option.NodeOptions
	raftOptions *option.RaftOptions
	serverId    *entity.PeerId

	// Other services
	confCtx               *configurationCtx
	logStorage            storage.LogStorage
	metaStorage           storage.RaftMetaStorage
	closureQueue          closure.ClosureQueue
	configManager         *conf.ConfigurationManager
	logManager            storage.LogManager
	fsmCaller             raft.FSMCaller
	ballotBox             *BallotBox
	snapshotExecutor      storage.SnapshotExecutor
	replicatorGroup       raft.ReplicatorGroup
	shutdownContinuations []raft.Closure
	rpcService            rpc.RaftClientService
	readOnlyService       raft.ReadOnlyService

	// Timers
	timerManager    Scheduler
	electionTimer   *util.RepeatedTimer
	voteTimer       *util.RepeatedTimer
	stepDownTimer   *util.RepeatedTimer
	snapshotTimer   *util.RepeatedTimer
	transferTimer   *util.RepeatedTimer
	wakingCandidate *util.ThreadId

	// Disruptor to run node service
	applyDisruptor disruptor.Disruptor
	applyQueue     disruptor.RingBuffer

	// Metrics
	metrics *NodeMetrics

	nodeId *entity.NodeId

	// ReplicatorStateListeners
	replicatorStateListeners []ReplicatorStateListener
	// Node's target leader election priority value
	targetPriority int // atomic
	// The number of elections time out for current node
	electionTimeoutCounter int
}

type stopTransferArg struct {
}

// Configuration commit context.
type configurationCtx struct {
}

type nodeReadWriteLock struct {
}

// Node service event.
type logEntryAndClosure struct {
}

type logEntryAndClosureFactory struct {
}

// Event handler.
type logEntryAndClosureHandler struct {
}
