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
	// @param task task to apply
	Apply(task *entity.Task)

	// ReadIndex
	// [Thread-safe and wait-free]
	// Starts a linearizable read-only query request with request context(optional,
	// such as request id etc.) and closure.  The closure will be called when the
	// request is completed, and user can read data from state machine if the result
	// status is OK.
	// @param requestContext the context of request
	// @param done           callback
	ReadIndex(requestContext []byte, done closure.ReadIndexClosure)

	// ListPeers List peers of this raft group, only leader returns.
	// [NOTE] when list_peers concurrency with {@link #addPeer(PeerId, Closure)}/{@link #removePeer(PeerId, Closure)},
	// maybe return peers is staled.  Because {@link #addPeer(PeerId, Closure)}/{@link #removePeer(PeerId, Closure)}
	// immediately modify configuration in memory</strong>
	// @return the peer list
	ListPeers() []*entity.PeerId

	// ListAlivePeers List all alive peers of this raft group, only leader returns.</p>
	// [NOTE] list_alive_peers is just a transient data (snapshot)
	// and a short-term loss of response by the follower will cause it to
	// temporarily not exist in this list.</strong>
	// @return the alive peer list
	ListAlivePeers() []*entity.PeerId

	// ListLearners List all learners of this raft group, only leader returns.</p>
	// [NOTE] when listLearners concurrency with {@link #addLearners(List, Closure)}/{@link #removeLearners(List, Closure)}/{@link #resetLearners(List, Closure)},
	// maybe return peers is staled.  Because {@link #addLearners(List, Closure)}/{@link #removeLearners(List, Closure)}/{@link #resetLearners(List, Closure)}
	// immediately modify configuration in memory</strong>
	// @return the learners set
	ListLearners() []*entity.PeerId

	// ListAliveLearners List all alive learners of this raft group, only leader returns.</p>
	// [NOTE] when listAliveLearners concurrency with {@link #addLearners(List, Closure)}/{@link #removeLearners(List, Closure)}/{@link #resetLearners(List, Closure)},
	// maybe return peers is staled.  Because {@link #addLearners(List, Closure)}/{@link #removeLearners(List, Closure)}/{@link #resetLearners(List, Closure)}
	// immediately modify configuration in memory</strong>
	// @return the  alive learners set
	ListAliveLearners() []*entity.PeerId

	// AddPeer Add a new peer to the raft group. done.run() would be invoked after this
	// operation finishes, describing the detailed result.
	// @param peer peer to add
	// @param done callback
	AddPeer(peer *entity.PeerId, done Closure)

	// RemovePeer Remove the peer from the raft group. done.run() would be invoked after
	// operation finishes, describing the detailed result.
	// @param peer peer to remove
	// @param done callback
	RemovePeer(peer *entity.PeerId, done Closure)

	// ChangePeers Change the configuration of the raft group to |newPeers| , done.run()
	// would be invoked after this operation finishes, describing the detailed result.
	// @param newPeers new peers to change
	// @param done     callback
	ChangePeers(newPeers *conf.Configuration, done Closure)

	// ResetPeers Reset the configuration of this node individually, without any replication
	// to other peers before this node becomes the leader. This function is
	// supposed to be invoked when the majority of the replication group are
	// dead and you'd like to revive the service in the consideration of availability.
	// Notice that neither consistency nor consensus are guaranteed in this
	// case, BE CAREFUL when dealing with this method.
	// @param newPeers new peers
	ResetPeers(newPeers *conf.Configuration) *Status

	// AddLearners Add some new learners to the raft group. done.run() will be invoked after this
	// operation finishes, describing the detailed result.
	// @param learners learners to add
	// @param done     callback
	AddLearners(learners []*entity.PeerId, done Closure)

	// RemoveLearners Remove some learners from the raft group. done.run() will be invoked after this
	// operation finishes, describing the detailed result.
	// @param learners learners to remove
	// @param done     callback
	RemoveLearners(learners []*entity.PeerId, done Closure)

	// ResetLearners Reset learners in the raft group. done.run() will be invoked after this
	// operation finishes, describing the detailed result.
	// @param learners learners to set
	// @param done     callback
	ResetLearners(learners []*entity.PeerId, done Closure)

	// Snapshot Start a snapshot immediately if possible. done.run() would be invoked when
	// the snapshot finishes, describing the detailed result.
	// @param done callback
	Snapshot(done Closure)

	// ResetElectionTimeoutMs Reset the election_timeout for the every node.
	// @param electionTimeoutMs the timeout millis of election
	ResetElectionTimeoutMs(electionTimeoutMs int)

	// TransferLeadershipTo Try transferring leadership to |peer|. If peer is ANY_PEER, a proper follower
	// will be chosen as the leader for the next term.
	// Returns 0 on success, -1 otherwise.
	// @param peer the target peer of new leader
	// @return operation status
	TransferLeadershipTo(peer entity.PeerId) *Status

	// ReadCommittedUserLog Read the first committed user log from the given index.
	// Return OK on success and user_log is assigned with the very data. Be awared
	// that the user_log may be not the exact log at the given index, but the
	// first available user log from the given index to lastCommittedIndex.
	// Otherwise, appropriate errors are returned:
	// 	- return ELOGDELETED when the log has been deleted;
	//  - return ENOMOREUSERLOG when we can't get a user log even reaching lastCommittedIndex.
	// [NOTE] in consideration of safety, we use lastAppliedIndex instead of lastCommittedIndex
	// in code implementation.
	// @param index log index
	// @return user log entry
	// @return LogNotFoundError  the user log is deleted at index.
	// @return LogIndexOutOfBoundsError  the special index is out of bounds.
	ReadCommittedUserLog(index int64) (*entity.UserLog, error)

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
