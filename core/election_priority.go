package core

type ElectionPriority int

const (
	Disabled   ElectionPriority = -1
	NotElected ElectionPriority = 0
	MinValue   ElectionPriority = 1
)
