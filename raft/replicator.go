package raft

import "time"

type ReplicatorStatus struct {
	LastRpcSendTimestamp int64 // atomic
}

type ReplicatorOptions struct {
	DynamicHeartbeatTimeoutMs *int
	ElectionTimeoutMs         *int
	GroupId                   GroupId
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

//
type St int

const (
	StIdle St = 0
	StBlocking
	StAppendingEntries
	StInstallingSnapshot
)

type Stat struct {
	st               St
	firstLogIndex    int64
	lastLogIncluded  int64
	lastLogIndex     int64
	LastTermIncluded int64
}

type FlyingAppendEntriesRpc struct {
}

//Replicator start
type Replicator struct {
	nextIndex               int64
	flyingAppendEntriesSize int64
	consecutiveErrorTimes   int
	hasSucceeded            bool
	timeoutNowIndex         int64
	heartbeatCounter        int64
	appendEntriesCounter    int64
	InstallSnapshotCounter  int64
	readonlyIndex           int64
	st                      Stat
	AppendEntriesInfly      []FlyingAppendEntriesRpc
}

func Start(options *ReplicatorOptions) (ReplicatorId, error) {
	//if options.LogManager == nil || options.
	return 0, nil
}

func Stop(id ReplicatorId) error {

}

func Join(id Replicator) error {

}

func WaitForCaughtUp(id ReplicatorId, maxMargin int64, dueTime time.Time, done *CatchupClosure) {

}
