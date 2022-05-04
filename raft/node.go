package raft

import (
	"math/rand"
)

var (
	// init from gflag
	raftMaxElectionDelayMs int
)

type NodeTimer struct {
	task *RepeatedTimerTask
	node *NodeImpl
}

func NewNodeTimer() *NodeTimer {
	timer := &NodeTimer{
		task: NewRepeatedTimerTask(),
		node: nil,
	}
	return timer
}

func (t *NodeTimer) Init(node *NodeImpl, timeoutMs int) error {
	t.task.Init(timeoutMs)
	t.node = node
	return nil
}

func (t *NodeTimer) onDestroy() {
	if t.node != nil {
		t.node = nil
	}
}

type ElectionTimer struct {
	*NodeTimer
}

func (t *ElectionTimer) Run() {
	t.node.HandleElectionTimeout()
}
func (t *ElectionTimer) AdjustTimeoutMs(timeoutMs int) int {
	return randomTimeout(timeoutMs)
}

type VoteTimer struct {
	*NodeTimer
}

func (t *VoteTimer) Run() {
	t.node.HandleVoteTimeout()
}
func (t *VoteTimer) AdjustTimeoutMs(timeoutMs int) int {
	return randomTimeout(timeoutMs)
}

type SnapshotTimer struct {
	*NodeTimer
	firstSchedule bool
}

func NewSnapshotTimer() *SnapshotTimer {
	timer := &SnapshotTimer{}
	timer.NodeTimer = NewNodeTimer()
	timer.firstSchedule = true
	return timer
}

func (t *SnapshotTimer) Run() {
	t.node.HandleSnapshotTimeout()
}
func (t *SnapshotTimer) AdjustTimeoutMs(timeoutMs int) int {
	if !t.firstSchedule {
		return timeoutMs
	}
	if timeoutMs > 0 {
		timeoutMs = rand.Intn(timeoutMs) + 1
	}
	t.firstSchedule = false
	return timeoutMs
}

func randomTimeout(timeoutMs int) int {
	delta := timeoutMs
	if timeoutMs > raftMaxElectionDelayMs {
		delta = raftMaxElectionDelayMs
	}
	return rand.Intn(delta) + timeoutMs
}

type State int

const (
	StateLeader        State = 1
	StateTransferring        = 2
	StateCandidate           = 3
	StateFollower            = 4
	StateError               = 5
	StateUninitialized       = 6
	StateShutting            = 7
	StateShutdown            = 8
	StateEnd                 = 9
)

func (s State) String() string {
	str := []string{"LEADER", "TRANSFERRING", "CANDIDATE", "FOLLOWER",
		"ERROR", "UNINITIALIZED", "SHUTTING", "SHUTDOWN"}
	if s < StateEnd {
		return str[s-1]
	}
	return "UNKNOWN"
}

type NodeImpl struct {
	state       State
	currentTerm int64
	leaderId    *PeerId
	votedId     *PeerId
}

func NewNodeImpl(groupId GroupId, peerId *PeerId) *NodeImpl {

}

func GetNodeId() *NodeId {
	return NewNodeId()
}

// timer func
func (n *NodeImpl) HandleElectionTimeout() {}

func (n *NodeImpl) HandleVoteTimeout() {

}

func (n *NodeImpl) HandleSnapshotTimeout() {

}

type Stage int

const (
	StageNoe Stage = iota
	StageCatchingUp
	StageJoint
	StageStable
)

type ConfigurationCtx struct {
	node     *NodeImpl
	stage    Stage
	nChanges int
	version  int64
	newPeers map[PeerId]struct{}
}

type VoteBallotCtx struct {
	timer
}

func NewVoteBallotCtx() *VoteBallotCtx {
}
