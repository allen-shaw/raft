package raft

type state struct {
	Code int
	Msg  string
}

type Status struct {
	s state
}
