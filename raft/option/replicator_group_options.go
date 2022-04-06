package option

import (
	"github.com/AllenShaw19/raft/raft"
	"github.com/AllenShaw19/raft/raft/core"
	"github.com/AllenShaw19/raft/raft/rpc"
	"github.com/AllenShaw19/raft/raft/storage"
)

type ReplicatorGroupOptions struct {
	heartbeatTimeoutMs   int
	electionTimeoutMs    int
	logManager           storage.LogManager
	ballotBox            *core.BallotBox
	node                 raft.Node
	snapshotStorage      storage.SnapshotStorage
	raftRpcClientService rpc.RaftClientService
	raftOptions          *RaftOptions
	//timerManager       	Scheduler
}

func (r *ReplicatorGroupOptions) HeartbeatTimeoutMs() int {
	return r.heartbeatTimeoutMs
}

func (r *ReplicatorGroupOptions) SetHeartbeatTimeoutMs(heartbeatTimeoutMs int) {
	r.heartbeatTimeoutMs = heartbeatTimeoutMs
}

func (r *ReplicatorGroupOptions) ElectionTimeoutMs() int {
	return r.electionTimeoutMs
}

func (r *ReplicatorGroupOptions) SetElectionTimeoutMs(electionTimeoutMs int) {
	r.electionTimeoutMs = electionTimeoutMs
}

func (r *ReplicatorGroupOptions) LogManager() storage.LogManager {
	return r.logManager
}

func (r *ReplicatorGroupOptions) SetLogManager(logManager storage.LogManager) {
	r.logManager = logManager
}

func (r *ReplicatorGroupOptions) BallotBox() *core.BallotBox {
	return r.ballotBox
}

func (r *ReplicatorGroupOptions) SetBallotBox(ballotBox *core.BallotBox) {
	r.ballotBox = ballotBox
}

func (r *ReplicatorGroupOptions) Node() raft.Node {
	return r.node
}

func (r *ReplicatorGroupOptions) SetNode(node raft.Node) {
	r.node = node
}

func (r *ReplicatorGroupOptions) SnapshotStorage() storage.SnapshotStorage {
	return r.snapshotStorage
}

func (r *ReplicatorGroupOptions) SetSnapshotStorage(snapshotStorage storage.SnapshotStorage) {
	r.snapshotStorage = snapshotStorage
}

func (r *ReplicatorGroupOptions) RaftRpcClientService() rpc.RaftClientService {
	return r.raftRpcClientService
}

func (r *ReplicatorGroupOptions) SetRaftRpcClientService(raftRpcClientService rpc.RaftClientService) {
	r.raftRpcClientService = raftRpcClientService
}

func (r *ReplicatorGroupOptions) RaftOptions() *RaftOptions {
	return r.raftOptions
}

func (r *ReplicatorGroupOptions) SetRaftOptions(raftOptions *RaftOptions) {
	r.raftOptions = raftOptions
}
