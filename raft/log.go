package raft

import (
	"bytes"
	"fmt"
	"os"
	"sync"
)

type entryHeader struct {
	term         int64
	entryType    int
	checksumType int
	dataLen      uint32
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

func (s *Segment) Create() {}

func (s *Segment) Load(m *ConfigurationManager) {}

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

// SegmentLogStorage implement LogStorage
type SegmentLogStorage struct {
}
