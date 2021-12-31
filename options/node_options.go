package options

import (
	"github.com/AllenShaw19/raft"
	"github.com/AllenShaw19/raft/conf"
)

type NodeOptions struct {
	ElectionTimeoutMs    int // 选举时间，follower在electionTimeoutMs时间内没有收到心跳，就变成candidate
	SnapshotIntervalSecs int // 快照的间隔时间
	initialConf          *conf.Configuration
	fsm                  raft.StateMachine // 业务状态机
	logUri               string            // raft log的存储路径 ${type}://${parameters}
	raftMetaUri          string            // raft 元数据存储路径
	snapshotUri          string            // 快照存储路径

}

func (opt *NodeOptions) SetRaftMetaUri(raftMetaUri string) {
	opt.raftMetaUri = raftMetaUri
}

func (opt *NodeOptions) SetSnapshotUri(snapshotUri string) {
	opt.snapshotUri = snapshotUri
}

func (opt *NodeOptions) SetLogUri(logUri string) {
	opt.logUri = logUri
}

func (opt *NodeOptions) SetInitialConf(conf *conf.Configuration) {
	opt.initialConf = conf
}

func (opt *NodeOptions) SetFsm(fsm raft.StateMachine) {
	opt.fsm = fsm
}

func DefaultNodeOptions() *NodeOptions {
	opts := &NodeOptions{}
	return opts
}
