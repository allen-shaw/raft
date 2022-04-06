package timer

type RaftTimerFactory interface {
	GetElectionTimer(shared bool, name string) Timer

	GetVoteTimer(shared bool, name string) Timer

	GetStepDownTimer(shared bool, name string) Timer

	GetSnapshotTimer(shared bool, name string) Timer

	CreateTimer(name string) Timer

	//    Scheduler getRaftScheduler(final boolean shared, final int workerNum, final String name);
	//    Scheduler createScheduler(final int workerNum, final String name);
}
