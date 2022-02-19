package raft

type FSMCaller struct {
}

func (f *FSMCaller) OnCommitted(committedIndex int64) error {

	return nil
}
