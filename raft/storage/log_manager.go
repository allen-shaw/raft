package storage

import (
	"github.com/AllenShaw19/raft/raft"
	"github.com/AllenShaw19/raft/raft/conf"
	"github.com/AllenShaw19/raft/raft/entity"
)

// Closure to run in stable state.
type stableClosure struct {
	firstLogIndex int64
	entries       []*entity.LogEntry
	nEntries      int
}

func (s *stableClosure) GetFirstLogIndex() int64 {
	return s.firstLogIndex
}

func (s *stableClosure) SetFirstLogIndex(firstLogIndex int64) {
	s.firstLogIndex = firstLogIndex
}

func (s *stableClosure) GetEntries() []*entity.LogEntry {
	return s.entries
}

func (s *stableClosure) SetEntries(entries []*entity.LogEntry) {
	s.entries = entries
	s.nEntries = len(entries)
}

func newStableClosure(entries []*entity.LogEntry) *stableClosure {
	sc := &stableClosure{}
	sc.SetEntries(entries)
	return sc
}

// Listen on last log index change event, but it's not reliable,
// the user should not count on this listener to receive all changed events.
type lastLogIndexListener interface {
	// OnLastLogIndexChanged Called when last log index is changed.
	OnLastLogIndexChanged(lastLogIndex int64)
}

// New log notifier callback.
type newLogCallback interface {
	// OnNewLog Called while new log come in.
	// @param arg the waiter pass-in argument
	// @param errorCode error code
	OnNewLog(arg interface{}, errorCode int)
}

type LogManager interface {
	// Join Wait the log manager to be shut down.
	// return InterruptedError if the current thread is interrupted while waiting
	Join() error

	// HasAvailableCapacityToAppendEntries
	// Given specified <tt>requiredCapacity</tt> determines if that amount of space
	// is available to append these entries. Returns true when available.
	// Returns true when available.
	HasAvailableCapacityToAppendEntries(requiredCapacity int) bool

	// AppendEntries Append log entry vector and wait until it's stable (NOT COMMITTED!)
	// @param entries log entries
	// @param done    callback
	AppendEntries(entries []*entity.LogEntry, closure *stableClosure)

	// SetSnapshot Notify the log manager about the latest snapshot, which indicates the
	// logs which can be safely truncated.
	SetSnapshot(meta entity.SnapshotMeta)

	// ClearBufferedLogs We don't delete all the logs before last snapshot to avoid installing
	// snapshot on slow replica. Call this method to drop all the logs before
	// last snapshot immediately.
	ClearBufferedLogs()

	// GetEntry Get the log entry at index.
	// @param the index of log entry
	// @return the log entry with {@code index}
	GetEntry(index int64) *entity.LogEntry

	// GetTerm Get the log term at index.
	// @param  the index of log entry
	// @return the term of log entry
	GetTerm(index int64) int64

	// GetFirstLogIndex Get the first log index of log
	GetFirstLogIndex() int64

	// GetLastLogIndex Get the last log index of log
	// @param isFlush whether to flush from disk.
	GetLastLogIndex(isFlush bool) int64

	// GetLastLogId Return the id the last log.
	// @param isFlush whether to flush all pending task.
	GetLastLogId(isFlush bool) *entity.LogEntry

	// GetConfiguration Get the configuration at index.
	GetConfiguration(index int64) *conf.ConfigurationEntry

	// CheckAndSetConfiguration Check if |current| should be updated to the latest configuration
	// Returns the latest configuration, otherwise null.
	CheckAndSetConfiguration(current *conf.ConfigurationEntry) *conf.ConfigurationEntry

	// Wait until there are more logs since |last_log_index| and |on_new_log|
	// would be called after there are new logs or error occurs, return the waiter id.
	// @param expectedLastLogIndex  expected last index of log
	// @param cb                    callback
	// @param arg                   the waiter pass-in argument
	Wait(expectedLastLogIndex int64, cb newLogCallback, arg interface{})

	// RemoveWaiter Remove a waiter.
	// @param id waiter id
	// @return true on success
	RemoveWaiter(id int64) bool

	// SetAppliedId Set the applied id, indicating that the log before applied_id (included)
	// can be dropped from memory logs.
	SetAppliedId(appliedId *entity.LogId)

	// CheckConsistency Check log consistency, returns the status
	// @return status
	CheckConsistency() *raft.Status
}
