package raft

type Closure interface {
	run(status Status)
}
