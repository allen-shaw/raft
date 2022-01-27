package raft

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/AllenShaw19/raft/log"
)

const (
	RAFT_SEGMENT_OPEN_PATTERN    = "log_inprogress_%020d"
	BRAFT_SEGMENT_CLOSED_PATTERN = "log_%020d_%020d"
	BRAFT_SEGMENT_META_FILE      = "log_meta"
)

type ChecksumType int

const (
	CHECKSUM_MURMURHASH32 ChecksumType = iota
	CHECKSUM_CRC32
)

type entryHeader struct {
	term         int64
	entryType    int
	checksumType int
	dataLen      uint64
	dataChecksum uint32
}

func (h *entryHeader) String() string {
	return fmt.Sprintf("{term=%v, type=%v, data_len=%v, checksum_type=%v, data_checksum=%v",
		h.term, h.entryType, h.dataLen, h.checksumType, h.dataChecksum)
}

type logMeta struct {
	offset int64
	length uint64
	term   int64
}

type offsetAndTerm struct {
	offset int64
	term   int64
}

type Segment struct {
	path          string
	bytes         int64
	mutex         sync.Mutex
	file          *os.File
	isOpen        bool
	firstIndex    int64
	lastIndex     int64 // should be atomic
	checksumType  int
	offsetAndTerm []offsetAndTerm
}

func NewSegment(path string, firstIndex, lastIndex int64, checksumType int) *Segment {
	s := &Segment{}
	s.path = path
	s.bytes = 0
	s.file = nil
	s.isOpen = false
	s.firstIndex = firstIndex
	if lastIndex == 0 {
		lastIndex = firstIndex - 1
	}
	s.lastIndex = lastIndex
	s.checksumType = checksumType
	return s
}

func (s *Segment) Create() error {
	if !s.isOpen {
		log.Error("Check failed: Create on a closed segment at first_index=%s in %s", s.firstIndex, s.path)
		return fmt.Errorf("fail to create on a closed segment")
	}
	path := filepath.Join(s.path, fmt.Sprintf(RAFT_SEGMENT_OPEN_PATTERN, s.firstIndex))
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error("Fail to open file %s", path)
		return err
	}
	log.Info("Created new segment %s", path)
	s.file = f
	return nil
}

func (s *Segment) Load(m *ConfigurationManager) {
	
}

func (s *Segment) Append(entry *LogEntry) {}

func (s *Segment) Get(index int64) *LogEntry {
	return nil
}

func (s *Segment) GetTerm(index int64) int64 {}

func (s *Segment) Close(sync bool) {}

func (s *Segment) Sync(sync bool) {}

func (s *Segment) Unlink() {}

func (s *Segment) Truncate(lastIndexKept int64) {}

func (s *Segment) IsOpen() bool {
	return s.isOpen
}

func (s *Segment) Bytes() int64 {
	return s.bytes
}

func (s *Segment) FirstIndex() int64 {
	return s.firstIndex
}

func (s *Segment) LastIndex() int64 {
	return s.lastIndex
}

func (s *Segment) FileName() string {

}

func (s *Segment) loadEntry(offset int64, head *entryHeader, body *bytes.Buffer, sizeHint uint64) {

}

func (s *Segment) getMeta(index int64) (*logMeta, error) {

}

func (s *Segment) truncateMetaAndGetLast(last int64) int {

}

func (s *Segment) String() string {

}

// SegmentLogStorage implement LogStorage
type SegmentLogStorage struct {
}

// util function
func verifyChecksum(checksumType int, data []byte, value uint32) bool {
	switch ChecksumType(checksumType) {
	case CHECKSUM_MURMURHASH32:
		return (value == murmurhash32(data))
	case CHECKSUM_CRC32:
		return (value == crc32(data))
	default:
		log.Error("Unknown checksum type=%v", checksumType)
		return false
	}
}

func getChecksum(checksumType int, data []byte) (uint32, error) {
	switch ChecksumType(checksumType) {
	case CHECKSUM_MURMURHASH32:
		return murmurhash32(data), nil
	case CHECKSUM_CRC32:
		return crc32(data), nil
	default:
		log.Error("Unknown checksum type=%v", checksumType)
		return 0, errUnknownChecksumType
	}
}
