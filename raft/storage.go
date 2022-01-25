package raft

import (
	"fmt"
	"strings"

	"github.com/AllenShaw19/raft/log"
	"github.com/AllenShaw19/raft/utils"
)

type logStorageExtension struct {
	NewInstance func(string) LogStorage
	Instance    LogStorage
}

var (
	logStorageExtensions map[string]*logStorageExtension
)

type IOMetric struct {
	StartTimeUs       int64
	QueueTimeUs       int64
	OpenSegmentTimeUs int64
	AppendEntryTimeUs int64
	SyncSegmentTimeUs int64
}

func (m *IOMetric) String() string {
	return fmt.Sprintf("queue_time_us: %v, open_segment_time_us: %v, append_entry_time_us: %v, sync_segment_time_us: %v",
		m.QueueTimeUs, m.OpenSegmentTimeUs, m.AppendEntryTimeUs, m.SyncSegmentTimeUs)
}

type LogStorage interface {
	Init(manager *ConfigurationManager) error
	FirstLogIndex() int64
	LastLogIndex() int64
	GetEntry(index int64) *LogEntry
	GetTerm(index int64) int64
	AppendEntry(entry *LogEntry) error
	AppendEntries(entries []*LogEntry, metric *IOMetric) (int, error)

	// delete logs from storage's head, [first_log_index, first_index_kept) will be discarded
	TruncatePrefix(firstIndexKept int64) error
	// delete uncommitted logs from storage's tail, (last_index_kept, last_log_index] will be discarded
	TruncateSuffix(lastIndexKept int64) error

	// Drop all the existing logs and reset next log index to |next_log_index|.
	// This function is called after installing snapshot from leader
	Reset(nextLogIndex int64) error

	GcInstance(uri string) utils.Status
	Destroy(uri string) utils.Status
}

func CreateLogStorage(uri string) LogStorage {
	protocol, parameter, err := parseUri(uri)
	if err != nil {
		log.Error("Parse log storage uri='%s', error: %v", uri, err)
		return nil
	}
	if ext, ok := logStorageExtensions[protocol]; ok {
		return ext.NewInstance(parameter)
	}
	log.Error("Fail to find log storage type %s, uri=%s", protocol, uri)
	return nil
}

func DestroyLogStorage(uri string) utils.Status {
	var status utils.Status
	protocol, parameter, err := parseUri(uri)
	if err != nil {
		log.Error("Parse log storage uri='%s', error: %v", uri, err)
		status.SetError(int32(RaftError_EINVAL), "Invalid log storage uri")
		return status
	}
	if ext, ok := logStorageExtensions[protocol]; ok {
		return ext.Instance.GcInstance(parameter)
	}
	log.Error("Fail to find log storage type %s, uri=%s", protocol, uri)
	status.SetError(int32(RaftError_EINVAL), "Invalid log storage uri")
	return status
}

func parseUri(uri string) (protocol string, parameter string, err error) {
	// ${protocol}://${parameters}
	s := strings.Split(uri, "://")
	if len(s) != 2 {
		return "", "", fmt.Errorf("invaild log storage uri")
	}
	return s[0], strings.TrimSpace(s[1]), nil
}
