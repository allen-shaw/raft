package raft

type LogStorage interface {
	Init(manager *ConfigurationManager) error
	FirstLogIndex() int64
	LastLogIndex() int64
}
