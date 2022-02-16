package raft

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"strings"

	"github.com/AllenShaw19/raft/log"
	"github.com/AllenShaw19/raft/utils"
)

type logStorageExtension struct {
	NewInstance func(string) LogStorage
	Instance    LogStorage
}
type metaStorageExtension struct {
	NewInstance func(string) MetaStorage
	Instance    MetaStorage
}
type snapshotStorageExtension struct {
	NewInstance func(string) SnapshotStorage
	Instance    SnapshotStorage
}

var (
	logStorageExtensions      map[string]*logStorageExtension
	metaStorageExtensions     map[string]*metaStorageExtension
	snapshotStorageExtensions map[string]*snapshotStorageExtension
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

// LogStorage Interface
type LogStorage interface {
	Init(manager *ConfigurationManager) error
	FirstLogIndex() int64
	LastLogIndex() int64
	GetEntry(index int64) *LogEntry
	GetTerm(index int64) int64
	AppendEntry(entry *LogEntry) error
	AppendEntries(entries []*LogEntry, metric *IOMetric) (int, error)

	// TruncatePrefix delete logs from storage's head, [first_log_index, first_index_kept) will be discarded
	TruncatePrefix(firstIndexKept int64) error
	// TruncateSuffix delete uncommitted logs from storage's tail, (last_index_kept, last_log_index] will be discarded
	TruncateSuffix(lastIndexKept int64) error

	// Reset Drop all the existing logs and reset next log index to |next_log_index|.
	// This function is called after installing snapshot from leader
	Reset(nextLogIndex int64) error

	GCInstance(uri string) utils.Status
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
		return ext.Instance.GCInstance(parameter)
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

// MetaStorage Interface
type MetaStorage interface {
	Init() error
	SetTermAndVotedFor(group VersionedGroupId, term int64, peerID *PeerId) error
	GetTermAndVotedFor(group VersionedGroupId) (term int64, peerID *PeerId, err error)
	NewInstance(uri string) *MetaStorage
	Create(uri string) *MetaStorage
	GCInstance(uri string, vgid VersionedGroupId) error
	Destroy(uri string, vgid VersionedGroupId) error
}

// Snapshot interface begin
type Snapshot interface {
	GetPath() string
	ListFiles() []string
	GetFileMeta(filename string, fileMeta proto.Message)
}

type SnapshotWriter interface {
	Snapshot
	SaveMeta(meta *SnapshotMeta) error
	AddFile(filename string, meta proto.Message) error
	RemoveFile(filename string) error
}

type SnapshotReader interface {
	Snapshot
	LoadMeta() (*SnapshotMeta, error)
	GenerateUriForCopy() string
}

type SnapshotCopier interface {
	Cancel()
	Join()
	GetReader() SnapshotReader
}

// SnapshotStorage begin
type SnapshotStorage interface {
	Init() error
	CreateWriter() SnapshotWriter
	CloseWriter(writer SnapshotWriter) error
	Open() SnapshotReader
	CloseReader(reader SnapshotReader) error
	CopyFrom(uri string) SnapshotReader
	StartToCopyFrom(uri string) SnapshotCopier
	CloseCopier(copier SnapshotCopier) error
	NewInstance(uri string) SnapshotStorage
	Create(uri string) SnapshotStorage
	GCInstance(uri string) error
	Destroy(uri string) error
}
