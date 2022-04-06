package core

type ReplicatorType int

const (
	Follower ReplicatorType = iota
	Learner
)

func IsFollower(replicatorType ReplicatorType) bool {
	return replicatorType == Follower
}

func IsLearner(replicatorType ReplicatorType) bool {
	return replicatorType == Learner
}
