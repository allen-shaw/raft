package core

const (
	// ElectionPriorityDisable Priority -1 represents this node disabled the priority election function.
	ElectionPriorityDisable int = -1

	// ElectionPriorityNotElected Priority 0 is a special value so that a node will never participate in election.
	ElectionPriorityNotElected = 0

	// ElectionPriorityMinValue Priority 1 is a minimum value for priority election.
	ElectionPriorityMinValue = 1
)
