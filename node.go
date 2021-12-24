package raft

import "github.com/AllenShaw19/raft/entity"

type Node interface {
	apply(task *entity.Task)
	getLeaderId() *entity.PeerId
	shutdown(done Closure)
	join()
	snapshot(done Closure)
}
