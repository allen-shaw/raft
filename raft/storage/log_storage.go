package storage

import "github.com/AllenShaw19/raft/raft/entity"

type LogStorage interface {
	// GetFirstLogIndex Returns first log index in log.
	GetFirstLogIndex() int64

	// GetLastLogIndex Returns last log index in log.
	GetLastLogIndex() int64

	// GetEntry Get logEntry by index.
	GetEntry(index int64) entity.LogEntry

	// AppendEntry Append entries to log.
	AppendEntry(entry *entity.LogEntry) bool

	// AppendEntries Append entries to log, return append success number
	AppendEntries(entries []*entity.LogEntry) int

	// TruncatePrefix Delete logs from storage's head
	// [first_log_index, first_index_kept) will be discarded.
	TruncatePrefix(firstIndexKept int64) bool

	// TruncateSuffix Delete uncommitted logs from storage's tail,
	// (last_index_kept, last_log_index] will be discarded.
	TruncateSuffix(lastIndexKept int64) bool

	// Reset Drop all the existing logs and reset next log index to |next_log_index|.
	// This function is called after installing snapshot from leader.
	Reset(nextLogIndex int64) bool
}
