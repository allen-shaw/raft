package core

import (
	"github.com/AllenShaw19/raft/log"
	"github.com/AllenShaw19/raft/raft"
	"github.com/AllenShaw19/raft/raft/closure"
	"github.com/AllenShaw19/raft/raft/conf"
	"github.com/AllenShaw19/raft/raft/entity"
	"github.com/AllenShaw19/raft/raft/option"
	"github.com/AllenShaw19/raft/raft/rpc"
	"github.com/AllenShaw19/raft/raft/storage"
	"github.com/AllenShaw19/raft/raft/util"
	"github.com/AllenShaw19/raft/raft/util/disruptor"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
)

var (
	GlobalNumNodes uint32 = 0 // atomic
)

// stopTransferArg
type stopTransferArg struct {
	node *Node
	term int64
	peer *entity.PeerId
}

func NewStopTransferArg(node *Node, term int64, peer *entity.PeerId) *stopTransferArg {
	arg := &stopTransferArg{}
	arg.node = node
	arg.term = term
	arg.peer = peer
	return arg
}

// Node service event.
type logEntryAndClosure struct {
}

type logEntryAndClosureFactory struct {
}

// Event handler.
type logEntryAndClosureHandler struct {
}

type Node struct {
	// Internal states
	mutex               sync.RWMutex
	state               State // atomic
	wg                  sync.WaitGroup
	currTerm            int64
	lastLeaderTimestamp int64 // atomic
	leaderId            *entity.PeerId
	votedId             *entity.PeerId
	voteCtx             *entity.Ballot // const
	prevVoteCtx         *entity.Ballot // const
	conf                *conf.ConfigurationEntry
	stopTransferArg     *stopTransferArg

	// Raft group and node options and identifier
	groupId     string // const
	options     *option.NodeOptions
	raftOptions *option.RaftOptions
	serverId    *entity.PeerId // const

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
	shutdownContinuations []raft.Closure // const
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
	//applyQueue     disruptor.RingBuffer

	// Metrics
	metrics *NodeMetrics

	nodeId *entity.NodeId

	// ReplicatorStateListeners
	replicatorStateListeners []ReplicatorStateListener //const

	// Node's target leader election priority value
	targetPriority int // atomic

	// The number of elections time out for current node
	electionTimeoutCounter int // atomic
}

func NewNode(groupId string, serverId *entity.PeerId) *Node {
	if groupId != "" {
		if err := util.VerifyGroupId(groupId); err != nil {
			return nil
		}
	}
	node := &Node{}
	node.groupId = groupId
	if serverId != nil {
		node.serverId = serverId.Copy()
	}
	node.state = StateUninitialized
	node.currTerm = 0
	node.updateLastLeaderTimestamp(util.MonotonicMs())
	node.confCtx = NewConfigurationCtx(node)
	node.wakingCandidate = nil

	num := atomic.AddUint32(&GlobalNumNodes, 1)
	log.Info("The number of active nodes increment to", zap.Uint32("num", num))
	return node
}

func (n *Node) Init(opt *option.NodeOptions) bool {
	n.options = opt
	n.raftOptions = opt.RaftOptions()
	n.metrics = NewNodeMetrics(opt.IsEnableMetrics())
	n.serverId.SetPriority(opt.ElectionPriority())
	n.electionTimeoutCounter = 0

	// TODO(allen): do some check

	// Init timers

	n.configManager = conf.NewConfigurationManager()

	// Init applyDisruptor

	n.fsmCaller = NewFSMCaller()
	if !n.initLogStorage() {
		log.Error("Node initLogStorage failed.", zap.String("NodeId", n.GetNodeId().String()))
		return false
	}
	if !n.initMetaStorage() {
		log.Error("Node initMetaStorage failed.", zap.String("NodeId", n.GetNodeId().String()))
		return false
	}
	if !n.initFSMCaller(entity.NewLogId(0, 0)) {
		log.Error("Node initFSMCaller failed.", zap.String("NodeId", n.GetNodeId().String()))
		return false
	}
	n.ballotBox = NewBallotBox()
	ballotBoxOpts := &option.BallotBoxOptions{}
	ballotBoxOpts.SetWaiter(n.fsmCaller)
	ballotBoxOpts.SetClosureQueue(n.closureQueue)
	if !n.ballotBox.Init(ballotBoxOpts) {
		log.Error("Node init ballotBox failed.", zap.String("NodeId", n.GetNodeId().String()))
		return false
	}

	if !n.initSnapshotStorage() {
		log.Error("Node initSnapshotStorage failed.", zap.String("NodeId", n.GetNodeId().String()))
		return false
	}

	st := n.logManager.CheckConsistency()
	if !st.IsOK() {
		log.Error("Node is initialized with inconsistent log",
			zap.String("NodeId", n.GetNodeId().String()),
			zap.String("status", st.String()))
		return false
	}

	n.conf = conf.NewConfigurationEntry()
	n.conf.SetId(entity.NewLogId(0, 0))

	// if have log using conf in log, else using conf in options
	if n.logManager.GetLastLogIndex() > 0 {
		n.checkAndSetConfiguration(false)
	} else {
		n.conf.SetConf(n.options.InitialConf())
		// initially set to max(priority of all nodes)
		n.targetPriority = n.getMaxPriorityOfNodes(n.conf.Conf().Peers())
	}

	if !n.conf.IsEmpty() {
		util.RequireTrue(n.conf.IsValid(), "Invalid conf: %s", n.conf)
	} else {
		log.Info("Init node with empty conf.", zap.String("serverId", n.serverId.String()))
	}

	// TODO RPC service and ReplicatorGroup is in cycle dependent, refactor it
	n.replicatorGroup = &ReplicatorGroup{}
	n.rpcService = NewDefaultRaftClientService(n.replicatorGroup)
	rgOpts := &option.ReplicatorGroupOptions{}

	//rgOpts.setHeartbeatTimeoutMs(heartbeatTimeout(this.options.getElectionTimeoutMs()));
	//rgOpts.setElectionTimeoutMs(this.options.getElectionTimeoutMs());
	//rgOpts.setLogManager(this.logManager);
	//rgOpts.setBallotBox(this.ballotBox);
	//rgOpts.setNode(this);
	//rgOpts.setRaftRpcClientService(this.rpcService);
	//rgOpts.setSnapshotStorage(this.snapshotExecutor != null ? this.snapshotExecutor.getSnapshotStorage() : null);
	//rgOpts.setRaftOptions(this.raftOptions);
	//rgOpts.setTimerManager(this.timerManager);

	n.options.SetMetricRegistry(n.metrics.Metrics())

	if !n.rpcService.Init(n.options) {
		log.Error("Fail to init rpc service.")
		return false
	}

	n.replicatorGroup.Init(entity.NewNodeId(n.groupId, n.serverId), rgOpts)
	n.readOnlyService = NewReadOnlyService()
	rosOpts := &option.ReadOnlyServiceOptions{}
	rosOpts.SetFsmCaller(n.fsmCaller)
	rosOpts.SetNode(n)
	rosOpts.SetRaftOptions(n.raftOptions)

	if !n.readOnlyService.Init(rosOpts) {
		log.Error("Fail to init readOnlyService.")
		return false
	}

	// set state to follower
	n.state = StateFollower

	log.Info("Node init",
		zap.String("nodeId", n.GetNodeId().String()),
		zap.Int64("term", n.currTerm),
		zap.String("lastLogId", n.logManager.GetLastLogId(false).String()),
		zap.String("conf", n.conf.Conf().String()),
		zap.String("oldConf", n.conf.OldConf().String()))

	if n.snapshotExecutor != nil && n.options.SnapshotIntervalSecs() > 0 {
		log.Debug("Node start snapshot timer", zap.String("nodeId", n.GetNodeId().String()), zap.Int64("term", n.currTerm))
		n.snapshotTimer.Start()
	}

	if !n.conf.IsEmpty() {
		n.stepDown(n.currTerm, false, raft.NewNilStatus())
	}

	if !raft.GetNodeManagerInstance().Add(n) {
		log.Error("NodeManager add node failed.", zap.String("nodeId", n.GetNodeId().String()))
		return false
	}

	// Now the raft node is started , have to acquire the writeLock to avoid race conditions
	n.mutex.Lock()
	defer n.mutex.Unlock()
	if n.conf.IsStable() && n.conf.Conf().Size() == 1 && n.conf.Conf().Contains(n.serverId) {
		// The group contains only this server which must be the LEADER, trigger the timer immediately.
		n.electSelf()
	}
}

func (n *Node) GetLeaderId() *entity.PeerId {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	return n.leaderId
}

func (n *Node) GetNodeId() *entity.NodeId {
	if n.nodeId == nil {
		n.nodeId = entity.NewNodeId(n.groupId, n.serverId)
	}
	return n.nodeId
}

func (n *Node) GetNodeMetrics() *NodeMetrics {
	return n.metrics
}

func (n *Node) GetGroupId() string {
	return n.groupId
}

func (n *Node) GetOptions() *option.NodeOptions {
	return n.options
}

func (n *Node) GetRaftOptions() *option.RaftOptions {
	return n.raftOptions
}

func (n *Node) IsLeader(blocking bool) bool {
	if !blocking {
		return n.GetNodeState() == StateLeader
	}
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.GetNodeState() == StateLeader
}

func (n *Node) Shutdown(done raft.Closure) {

}

func (n *Node) GetNodeState() State {
	return State(atomic.LoadInt32((*int32)(&n.state)))
}

// private function
func (n *Node) updateLastLeaderTimestamp(lastLeaderTimestamp int64) {
	n.lastLeaderTimestamp = lastLeaderTimestamp
}

func (n *Node) initLogStorage() bool {

}

func (n *Node) initMetaStorage() bool {

}

func (n *Node) initSnapshotStorage() bool {

}

func (n *Node) initFSMCaller(bootstrapId *entity.LogId) bool {

}

func (n *Node) checkAndSetConfiguration(inLock bool) {

}

func (n *Node) getMaxPriorityOfNodes(peers []*entity.PeerId) int {

}

func (n *Node) Join() error {
	//TODO implement me
	panic("implement me")
}

func (n *Node) Apply(task *entity.Task) {
	//TODO implement me
	panic("implement me")
}

func (n *Node) ReadIndex(requestContext []byte, done closure.ReadIndexClosure) {
	//TODO implement me
	panic("implement me")
}

func (n *Node) ListPeers() []*entity.PeerId {
	//TODO implement me
	panic("implement me")
}

func (n *Node) ListAlivePeers() []*entity.PeerId {
	//TODO implement me
	panic("implement me")
}

func (n *Node) ListLearners() []*entity.PeerId {
	//TODO implement me
	panic("implement me")
}

func (n *Node) ListAliveLearners() []*entity.PeerId {
	//TODO implement me
	panic("implement me")
}

func (n *Node) AddPeer(peer *entity.PeerId, done raft.Closure) {
	//TODO implement me
	panic("implement me")
}

func (n *Node) RemovePeer(peer *entity.PeerId, done raft.Closure) {
	//TODO implement me
	panic("implement me")
}

func (n *Node) ChangePeers(newPeers *conf.Configuration, done raft.Closure) {
	//TODO implement me
	panic("implement me")
}

func (n *Node) ResetPeers(newPeers *conf.Configuration) *raft.Status {
	//TODO implement me
	panic("implement me")
}

func (n *Node) AddLearners(learners []*entity.PeerId, done raft.Closure) {
	//TODO implement me
	panic("implement me")
}

func (n *Node) RemoveLearners(learners []*entity.PeerId, done raft.Closure) {
	//TODO implement me
	panic("implement me")
}

func (n *Node) ResetLearners(learners []*entity.PeerId, done raft.Closure) {
	//TODO implement me
	panic("implement me")
}

func (n *Node) Snapshot(done raft.Closure) {
	//TODO implement me
	panic("implement me")
}

func (n *Node) ResetElectionTimeoutMs(electionTimeoutMs int) {
	//TODO implement me
	panic("implement me")
}

func (n *Node) TransferLeadershipTo(peer entity.PeerId) *raft.Status {
	//TODO implement me
	panic("implement me")
}

func (n *Node) ReadCommittedUserLog(index int64) (*entity.UserLog, error) {
	//TODO implement me
	panic("implement me")
}

func (n *Node) AddReplicatorStateListener(replicatorStateListener ReplicatorStateListener) {
	//TODO implement me
	panic("implement me")
}

func (n *Node) RemoveReplicatorStateListener(replicatorStateListener ReplicatorStateListener) {
	//TODO implement me
	panic("implement me")
}

func (n *Node) ClearReplicatorStateListeners() {
	//TODO implement me
	panic("implement me")
}

func (n *Node) GetReplicatorStatueListeners() []ReplicatorStateListener {
	//TODO implement me
	panic("implement me")
}

func (n *Node) GetNodeTargetPriority() int {
	//TODO implement me
	panic("implement me")
}

func (n *Node) stepDown(term int64, wakeupCandidate bool, status *raft.Status) {

}

func (n *Node) electSelf() {

}
