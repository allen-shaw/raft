package utils

type State struct {
	Code    int
	Message string
}

type Status struct {
	state *State
}
