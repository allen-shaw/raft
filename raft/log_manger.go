package raft

type LogManagerOptions struct {
	LogStorage *LogStorage
	ConfMgr    *ConfigurationManager
	FSMCaller  *FSMCaller
}

type LogManagerStatus struct {
	FirstIndex        int64
	LastIndex         int64
	DiskIndex         int64
	KnownAppliedIndex int64
}

type LogManager struct {
}
