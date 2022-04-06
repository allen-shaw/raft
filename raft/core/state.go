package core

type State int32

const (
	StateLeader        State = iota // It's a leader
	StateTransferring               // It's transferring leadership
	StateCandidate                  //  It's a candidate
	StateFollower                   // It's a follower
	StateError                      // It's in error
	StateUninitialized              // It's uninitialized
	StateShutting                   // It's shutting down
	StateShutdown                   // It's shutdown already
	StateEnd                        // State end
)

func (s State) IsActive() bool {
	return s < StateError
}
