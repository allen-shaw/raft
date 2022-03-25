package raft

import (
	"github.com/AllenShaw19/raft/raft/closure"
	"github.com/AllenShaw19/raft/raft/conf"
	"github.com/AllenShaw19/raft/raft/core"
	"github.com/AllenShaw19/raft/raft/entity"
	"github.com/AllenShaw19/raft/raft/option"
)

// Node A raft replica node.
type Node interface {
	// GetLeaderId Get the leader peer id for redirect, null if absent.
	GetLeaderId() *entity.PeerId

	// GetNodeId Get current node id.
	GetNodeId() *entity.NodeId

	// GetNodeMetrics Get the node metrics,
	// only valid when node option {@link NodeOptions#IsEnableMetrics()} is true.
	GetNodeMetrics() *core.NodeMetrics

	// GetGroupId Get the raft group id.
	GetGroupId() string

	// GetOptions Get the node options.
	GetOptions() *option.NodeOptions

	// GetRaftOptions Get the raft options
	GetRaftOptions() *option.RaftOptions

	// IsLeader Returns true when the node is leader.
	// @param blocking if true, will be blocked until the node finish it's state change
	IsLeader(blocking bool) bool

	// Shutdown local replica node.
	// @param done callback
	Shutdown(done Closure)

	// Join Block the thread until the node is successfully stopped.
	// return InterruptedError if the current thread is interrupted while waiting
	Join() error

	// Apply task to the replicated-state-machine
	// [Thread-safe and wait-free]
	// About the ownership:
	// |task.data|: for the performance consideration, we will take away the
	//               content. If you want to keep the content, copy it before call
	//               this function
	// |task.done|: If the data is successfully committed to the raft group. We
	//              will pass the ownership to #{@link StateMachine#onApply(Iterator)}.
	//              Otherwise we will specify the error and call it.
	//
	// @param task task to apply
	Apply(task *entity.Task)

	// ReadIndex
	// [Thread-safe and wait-free]
	// Starts a linearizable read-only query request with request context(optional,
	// such as request id etc.) and closure.  The closure will be called when the
	// request is completed, and user can read data from state machine if the result
	// status is OK.
	//
	// @param requestContext the context of request
	// @param done           callback
	ReadIndex(requestContext []byte, done closure.ReadIndexClosure)

	ListPeers() []*entity.PeerId

	ListAlivePeers() []*entity.PeerId

	ListLearners() []*entity.PeerId

	ListAliveLearners() []*entity.PeerId

	AddPeer(peer *entity.PeerId, done Closure)

	RemovePeer(peer *entity.PeerId, done Closure)

	ChangePeers(newPeers *conf.Configuration, done Closure)

	ResetPeers(newPeers *conf.Configuration) *Status

	AddLearners(learners []*entity.PeerId, done Closure)

	RemoveLearners(learners []*entity.PeerId, done Closure)

	ResetLearners(learners []*entity.PeerId, done Closure)

	Snapshot(done Closure)

	ResetElectionTimeoutMs(electionTimeoutMs int)

	TransferLeadershipTo(peer entity.PeerId) *Status

	ReadCommittedUserLog(index int64) *entity.UserLog

	// AddReplicatorStateListener users can implement the ReplicatorStateListener interface by themselves.
	// So users can do their own logical operator in this listener when replicator created, destroyed or had some errors.
	// @param replicatorStateListener added ReplicatorStateListener which is implemented by users.
	AddReplicatorStateListener(replicatorStateListener core.ReplicatorStateListener)

	// RemoveReplicatorStateListener
	// End User can remove their implement the ReplicatorStateListener interface by themselves.
	// @param replicatorStateListener need to remove the ReplicatorStateListener which has been added by users
	RemoveReplicatorStateListener(replicatorStateListener core.ReplicatorStateListener)

	// ClearReplicatorStateListeners
	// Remove all the ReplicatorStateListeners which have been added by users.
	ClearReplicatorStateListeners()

	// GetReplicatorStatueListeners Get the ReplicatorStateListeners which have been added by users.
	// @return node's replicatorStatueListeners which have been added by users.
	GetReplicatorStatueListeners() []core.ReplicatorStateListener

	// GetNodeTargetPriority Get the node's target election priority value.
	// @return node's target election priority value.
	GetNodeTargetPriority() int

	// GetNodeState Get the node's state.
	// @return node's state.
	GetNodeState() core.State
}
