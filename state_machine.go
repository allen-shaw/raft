package raft

import (
	"github.com/AllenShaw19/raft/conf"
	"github.com/AllenShaw19/raft/entity"
	"github.com/AllenShaw19/raft/error"
	"github.com/AllenShaw19/raft/storage/snapshot"
)

type StateMachine interface {
	onApply(task *entity.Task)
	onShutdown()
	onSnapshotSave(writer snapshot.Writer, done Closure)
	onSnapshotLoad(reader snapshot.Reader) bool
	onLeaderStart(term int64)
	onLeaderStop(status Status)
	onError(e error.RaftError)
	onConfigCommitted(conf *conf.Configuration)
	onStopFollowing(ctx *entity.LeaderChangeContext)
	onStartFollowing(ctx *entity.LeaderChangeContext)
}
