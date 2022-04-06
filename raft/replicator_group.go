package raft

import (
	"github.com/AllenShaw19/raft/raft/core"
	"github.com/AllenShaw19/raft/raft/entity"
	"github.com/AllenShaw19/raft/raft/option"
	"github.com/AllenShaw19/raft/raft/rpc"
)

// ReplicatorGroup Replicators in a raft group.
type ReplicatorGroup interface {

	// Init the replicator group.
	// @param nodeId node id
	// @param opts   options of replicator group
	// @return true if init success
	Init(nodeId *entity.NodeId, opts *option.ReplicatorGroupOptions) bool

	// AddReplicatorForFollower Adds a replicator for follower({@link ReplicatorType#Follower}).
	// @param peer target peer
	// @return true on success
	AddReplicatorForFollower(peer *entity.PeerId) bool

	// AddReplicatorSync Add a replicator attached with |peer|
	// will be a notification when the replicator catches up according to the arguments.
	// NOTE: when calling this function, the replicators starts to work
	// immediately, and might call Node#stepDown which might have race with
	// the caller, you should deal with this situation.
	// @param peer           target peer
	// @param replicatorType replicator type
	// @return true on success
	AddReplicatorSync(peer *entity.PeerId, replicatorType core.ReplicatorType) bool

	// AddReplicator Try to add a replicator attached with |peer|
	// will be a notification when the replicator catches up according to the arguments.
	// NOTE: when calling this function, the replicators starts to work
	// immediately, and might call Node#stepDown which might have race with
	// the caller, you should deal with this situation.
	// @param peer           target peer
	// @param replicatorType replicator type
	// @param sync           synchronous
	// @return true on success
	AddReplicator(peer *entity.PeerId, replicatorType core.ReplicatorType, sync bool) bool

	// Send heartbeat to a peer.
	// @param peer    target peer
	// @param closure callback
	SendHeartbeat(peer *entity.PeerId, closure rpc.RpcResponseClosure[AppendEntriesResponse])
}
