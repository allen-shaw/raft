package raft

//A Status encapsulates the result of an operation. It may indicate success,
//or it may indicate an error with an associated error message. It's suitable
//for passing status of functions with richer information than just error_code
//in exception-forbidden code. This utility is inspired by leveldb::Status.
//
//Multiple threads can invoke const methods on a Status without
//external synchronization, but if any of the threads may call a
//non-const method, all threads accessing the same Status must use
//external synchronization.
//
//Since failed status needs to allocate memory, you should be careful when
//failed status is frequent.
type Status struct {
	state *state
}

func OK() *Status {
	return &Status{}
}

func NewStatus(code int, msg string) *Status {
	state := &state{code: code, msg: msg}
	return &Status{state: state}
}

func (s *Status) Reset() {
	s.state = nil
}

func (s *Status) IsOK() bool {
	return s.state == nil || s.state.code == 0
}

func (s *Status) SetCode(code int) {
	if s.state == nil {
		s.state = &state{code: code}
	} else {
		s.state.code = code
	}
}

func (s *Status) GetCode() int {
	if s.state == nil {
		return 0
	}
	return s.state.code
}

//Status internal state.
type state struct {
	code int
	msg  string
}
