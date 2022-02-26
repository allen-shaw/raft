package raft

import (
	"errors"
	"github.com/AllenShaw19/raft/log"
	"github.com/AllenShaw19/raft/utils"
	"sync"
	"sync/atomic"
)

type MemoryData []*LogEntry

type MemoryLogStorage struct {
	path          string
	firstLogIndex int64 // atomic<int64>
	lastLogIndex  int64 // atomic<int64>
	logEntryData  MemoryData
	mutex         sync.Mutex
}

func NewMemoryLogStorage(path string) *MemoryLogStorage {
	return &MemoryLogStorage{
		path:          path,
		firstLogIndex: 1,
		lastLogIndex:  0,
	}
}

func (s *MemoryLogStorage) Init(manager *ConfigurationManager) error {
	atomic.StoreInt64(&s.firstLogIndex, 1)
	atomic.StoreInt64(&s.lastLogIndex, 0)
	s.logEntryData = make(MemoryData, 0)
	return nil
}

// FirstLogIndex first log index in log
func (s *MemoryLogStorage) FirstLogIndex() int64 {
	return atomic.LoadInt64(&s.firstLogIndex)
}

// LastLogIndex last log index in log
func (s *MemoryLogStorage) LastLogIndex() int64 {
	return atomic.LoadInt64(&s.lastLogIndex)
}

// GetEntry get log entry by index
func (s *MemoryLogStorage) GetEntry(index int64) (*LogEntry, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if index < atomic.LoadInt64(&s.firstLogIndex) || index > atomic.LoadInt64(&s.lastLogIndex) {
		return nil, errors.New("")
	}
	temp := s.logEntryData[index-atomic.LoadInt64(&s.firstLogIndex)]
	if temp.ID.Index != index {
		log.Fatal("get_entry entry index not equal. "+
			"logEntry index: %v, required_index %v", temp.ID.Index, index)
	}
	return temp, nil
}

// GetTerm get log entry's term by index
func (s *MemoryLogStorage) GetTerm(index int64) (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if index < atomic.LoadInt64(&s.firstLogIndex) || index > atomic.LoadInt64(&s.lastLogIndex) {
		return 0, errors.New("")
	}
	temp := s.logEntryData[index-atomic.LoadInt64(&s.firstLogIndex)]
	if temp.ID.Index != index {
		log.Fatal("get_entry entry index not equal. "+
			"logEntry index: %v, required_index %v", temp.ID.Index, index)
	}
	return temp.ID.Term, nil
}

// AppendEntry append entries to log
func (s *MemoryLogStorage) AppendEntry(entry *LogEntry) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if entry.ID.Index != atomic.LoadInt64(&s.lastLogIndex)+1 {
		log.Error("input_entry index=%v, last_log_index=%v, first_log_index=%v",
			entry.ID.Index, s.lastLogIndex, s.firstLogIndex)
		return errRange
	}
	s.logEntryData = append(s.logEntryData, entry)
	atomic.AddInt64(&s.lastLogIndex, 1)
	return nil
}

// AppendEntries append entries to log and update IOMetric, return append success number
func (s *MemoryLogStorage) AppendEntries(entries []*LogEntry, metric *IOMetric) (int, error) {
	if len(entries) == 0 {
		return 0, nil
	}
	for i, entry := range entries {
		err := s.AppendEntry(entry)
		if err != nil {
			return i, err
		}
	}
	return len(entries), nil
}

// TruncatePrefix delete logs from storage's head, [first_log_index, first_index_kept) will be discarded
func (s *MemoryLogStorage) TruncatePrefix(firstIndexKept int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, entry := range s.logEntryData {
		if entry.ID.Index >= firstIndexKept {
			s.logEntryData = s.logEntryData[i:]
			break
		}
	}
	atomic.StoreInt64(&s.firstLogIndex, firstIndexKept)
	if atomic.LoadInt64(&s.firstLogIndex) > atomic.LoadInt64(&s.lastLogIndex) {
		atomic.StoreInt64(&s.lastLogIndex, firstIndexKept-1)
	}
	return nil
}

// TruncateSuffix delete uncommitted logs from storage's tail, (last_index_kept, last_log_index] will be discarded
func (s *MemoryLogStorage) TruncateSuffix(lastIndexKept int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i := len(s.logEntryData) - 1; i >= 0; i-- {
		entry := s.logEntryData[i]
		if entry.ID.Index <= lastIndexKept {
			s.logEntryData = s.logEntryData[:i]
			break
		}
	}
	atomic.StoreInt64(&s.lastLogIndex, lastIndexKept)
	if atomic.LoadInt64(&s.firstLogIndex) > atomic.LoadInt64(&s.lastLogIndex) {
		atomic.StoreInt64(&s.firstLogIndex, lastIndexKept+1)
	}

	return nil
}

// Reset Drop all the existing logs and reset next log index to |next_log_index|.
// This function is called after installing snapshot from leader
func (s *MemoryLogStorage) Reset(nextLogIndex int64) error {
	if nextLogIndex <= 0 {
		log.Error("Invalid next_log_index=%v", nextLogIndex)
		return errInvalid
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.logEntryData = make(MemoryData, 0)
	atomic.StoreInt64(&s.firstLogIndex, nextLogIndex)
	atomic.StoreInt64(&s.lastLogIndex, nextLogIndex)

	return nil
}

// NewInstance Create an instance of this kind of LogStorage with the parameters encoded in |uri|
// Return the address referenced to the instance on success, NULL otherwise.
func (s *MemoryLogStorage) NewInstance(uri string) LogStorage {
	return NewMemoryLogStorage(uri)
}

// GCInstance GC an instance of this kind of LogStorage with the parameters encoded  in |uri|
func (s *MemoryLogStorage) GCInstance(uri string) utils.Status {
	return utils.Status{}
}
