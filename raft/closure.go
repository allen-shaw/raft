package raft

type Closure interface {
	// Run Called when task is done.
	// @param status the task status.
	Run(status *Status)
}
