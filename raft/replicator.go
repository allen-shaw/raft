package raft

import "time"

type ReplicatorStatus struct {
	LastRpcSendTimestamp int64 // atomic
}

type ReplicatorOptions struct {
	DynamicHeartbeatTimeoutMs *int
	ElectionTimeoutMs         *int
	GroupId                   GroupID
	ServerId                  PeerId
	PeerId                    PeerId
	LogManager                *LogManager
	BallotBox                 *BallotBox
	Node                      *NodeImpl
	Term                      int64
	SnapshotStorage           *SnapshotStorage
	SnapshotThrottle          *SnapshotThrottle
	ReplicatorStatus          *ReplicatorStatus
}

type ReplicatorId uint64

// CatchupClosure start
type CatchupClosure struct {
	maxMargin   int64
	timer       time.Timer
	hasTimer    bool
	errorWasSet bool
}

func (cc *CatchupClosure) Run() {

}

func (cc *CatchupClosure) run() {

}

//Replicator start
type Replicator struct {
}
