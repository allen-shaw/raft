package raft

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/AllenShaw19/raft/log"
	"github.com/AllenShaw19/raft/utils"
)

const (
	RaftSegmentOpenPattern   = "log_inprogress_%020d"
	RaftSegmentClosedPattern = "log_%020d_%020d"
	RAFT_SEGMENT_META_FILE   = "log_meta"

	EntryHeaderSize uint64 = 24
)

type ChecksumType int

const (
	ChecksumMurmurhash32 ChecksumType = iota
	ChecksumCrc32
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
		log.Error("Check failed: Create on a closed segment at first_index=%v in %s", s.firstIndex, s.path)
		return fmt.Errorf("fail to create on a closed segment")
	}
	path := filepath.Join(s.path, fmt.Sprintf(RaftSegmentOpenPattern, s.firstIndex))
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error("Fail to open file %s", path)
		return err
	}
	log.Info("Created new segment %s", path)
	s.file = f
	return nil
}

func (s *Segment) Load(m *ConfigurationManager) error {
	var err error
	path := s.path
	if s.isOpen {
		path = filepath.Join(path, fmt.Sprintf(RaftSegmentOpenPattern, s.firstIndex))
	} else {
		path = filepath.Join(path, fmt.Sprintf(RaftSegmentClosedPattern, s.firstIndex, atomic.LoadInt64(&s.lastIndex)))
	}

	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		log.Error("open file fail, path=%s, err=%v", path, err)
		return err
	}
	s.file = f
	stat, err := f.Stat()
	if err != nil {
		log.Error("get file stat fail, path=%s, err=%v", path, err)
		return err
	}

	// load entry index
	fileSize := stat.Size()
	entryOff := int64(0)
	actualLastIndex := s.firstIndex - 1
	for i := s.firstIndex; entryOff < fileSize; i++ {
		header := &entryHeader{}
		err = s.loadEntry(entryOff, header, nil, EntryHeaderSize)
		if err != nil {
			break
		}
		skipLen := EntryHeaderSize + header.dataLen
		if int64(skipLen)+entryOff > fileSize {
			break
		}
		if EntryType(header.entryType) == EntryType_ENTRY_TYPE_CONFIGURATION {
			data := bytes.NewBuffer(make([]byte, 0))
			if err := s.loadEntry(entryOff, nil, data, skipLen); err != nil {
				break
			}
			entry := NewLogEntry()
			entry.ID.Index = i
			entry.ID.Term = header.term
			status := ParseConfigurationMeta(data, entry)
			if status.OK() {
				confEntry := NewConfigurationEntry(entry)
				m.Add(confEntry)
			} else {
				log.Error("fail to parse configuration meta, path: %s entry_off %d", s.path, entryOff)
				err = errors.New("parse configuration meta fail")
				break
			}
		}
		s.offsetAndTerm = append(s.offsetAndTerm, offsetAndTerm{entryOff, header.term})
		actualLastIndex++
		entryOff += int64(skipLen)
	}

	lastIndex := atomic.LoadInt64(&s.lastIndex)
	if err == nil && !s.isOpen {
		if actualLastIndex < lastIndex {
			
		}
	}

	if err != nil {
		return err
	}

	if s.isOpen {
		atomic.StoreInt64(&s.lastIndex, actualLastIndex)
	}

	// truncate last uncompleted entry
	if entryOff != fileSize {
		log.Info("")
		err = ftruncateUninterrupted()
	}

	ret, err := s.file.Seek(entryOff, os.SEEK_SET)

	s.bytes = entryOff
	return err
}

func (s *Segment) Append(entry *LogEntry) {}

func (s *Segment) Get(index int64) *LogEntry {
	return nil
}

func (s *Segment) GetTerm(index int64) (int64, error) {
	meta, err := s.getMeta(index)
	if err != nil {
		return 0, err
	}
	return meta.term, nil
}

func (s *Segment) Close(sync bool) {}

func (s *Segment) Sync(sync bool) error {
	if s.lastIndex > s.firstIndex && sync {
		s.file.Sync()
	}
	return nil
}

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
	if !s.isOpen {
		return fmt.Sprintf(RaftSegmentClosedPattern, s.firstIndex, atomic.LoadInt64(&s.lastIndex))
	}
	return fmt.Sprintf(RaftSegmentOpenPattern, s.firstIndex)
}

func (s *Segment) loadEntry(offset int64, head *entryHeader, data *bytes.Buffer, sizeHint uint64) error {
	var buf bytes.Buffer
	toRead := maxUint64(sizeHint, EntryHeaderSize)

	n, err := filePread(&buf, s.file, offset, toRead)
	if err != nil {
		log.Error("read file %s fail, offset %d, to read size %d, err %v", s.path, offset, toRead, err)
		return err
	}
	if n != toRead {
		log.Error("read file %s fail, read len %d, to read %d", s.path, n, toRead)
		return errors.New("read file fail")
	}

	headerBuf := bytes.NewBuffer(buf.Next(int(EntryHeaderSize)))
	unpacker := utils.NewRawUnpacker(headerBuf)

	term := unpacker.Unpack64()
	metaField := unpacker.Unpack32()
	dataLen := unpacker.Unpack32()
	dataChecksum := unpacker.Unpack32()
	headerChecksum := unpacker.Unpack32()

	tmp := entryHeader{
		term:         int64(term),
		entryType:    int(metaField >> 24),
		checksumType: int((metaField << 8) >> 24),
		dataLen:      uint64(dataLen),
		dataChecksum: dataChecksum,
	}

	if !verifyChecksum(tmp.checksumType, headerBuf.Bytes()[:EntryHeaderSize-4], headerChecksum) {
		log.Error("Found corrupted header at offset=%d, header=%s, path: %s", offset, tmp.String(), s.path)
		return errors.New("verify header checksum fail")
	}
	if head != nil {
		*head = tmp
	}

	if data != nil {
		if buf.Len() < int(EntryHeaderSize+uint64(dataLen)) {
			toRead := EntryHeaderSize + uint64(dataLen) - uint64(buf.Len())
			n, err := filePread(&buf, s.file, offset+int64(buf.Len()), toRead)
			if err != nil {
				log.Error("read file %s fail, offset %d, to read size %d, err %v", s.path, offset+int64(buf.Len()), toRead, err)
				return err
			}
			if n != toRead {
				log.Error("read file %s fail, read len %d, to read %d", s.path, n, toRead)
				return errors.New("read file fail")
			}
		} else if buf.Len() > int(EntryHeaderSize)+int(dataLen) {
			buf.Truncate(int(EntryHeaderSize) + int(dataLen))
		}
		if buf.Len() != int(EntryHeaderSize)+int(dataLen) {
			log.Error("read file %s fail, read buffer len=%d, not equal data_size=%d", s.path, buf.Len(), uint32(EntryHeaderSize)+dataLen)
			return errors.New("read file fail")
		}
		buf.Next(int(EntryHeaderSize))
		if !verifyChecksum(tmp.checksumType, buf.Bytes(), tmp.dataChecksum) {
			log.Error("Found corrupted data at offset=%d, header=%v, path=%s", offset+int64(EntryHeaderSize), tmp, s.path)
			return errors.New("buffer checksum invalid")
		}
		data.Reset()
		n, err := data.ReadFrom(&buf)
		if err != nil {
			log.Error("data read from buf fail, err %v", err)
			return err
		}
		if n != int64(buf.Len()) {
			log.Error("swap buffer fail, buf_len=%d, read data_size=%d", buf.Len(), n)
			return errors.New("data read from buf fail")
		}
	}

	return nil
}

func (s *Segment) getMeta(index int64) (*logMeta, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if index > atomic.LoadInt64(&s.lastIndex) || index < s.firstIndex {
		log.Error("invalid index=%v, lastIndex=%v, firstIndex=%v",
			index, atomic.LoadInt64(&s.lastIndex), s.firstIndex)
		return nil, errors.New("invalid index")
	} else if s.lastIndex == s.firstIndex-1 {
		log.Error("lastIndex=%v, firstIndex=%v", atomic.LoadInt64(&s.lastIndex), s.firstIndex)
		return nil, errors.New("invalid lastIndex And firstIndex")
	}

	metaIndex := index - s.firstIndex
	entryCursor := s.offsetAndTerm[metaIndex].offset
	nextCursor := s.bytes
	if index < atomic.LoadInt64(&s.lastIndex) {
		nextCursor = s.offsetAndTerm[metaIndex+1].offset
	}

	if entryCursor >= nextCursor {
		log.Fatal("entryCursor %v >= nextCursor %v", entryCursor, nextCursor)
	}

	meta := &logMeta{
		offset: entryCursor,
		term:   s.offsetAndTerm[metaIndex].term,
		length: uint64(nextCursor - nextCursor),
	}
	return meta, nil
}

func (s *Segment) truncateMetaAndGetLast(last int64) int {

	return 0
}

func (s *Segment) String() string {
	return ""
}

// SegmentLogStorage implement LogStorage
type SegmentLogStorage struct {
}

// util function
func verifyChecksum(checksumType int, data []byte, value uint32) bool {
	switch ChecksumType(checksumType) {
	case ChecksumMurmurhash32:
		return (value == murmurhash32(data))
	case ChecksumCrc32:
		return (value == crc32(data))
	default:
		log.Error("Unknown checksum type=%v", checksumType)
		return false
	}
}

func getChecksum(checksumType int, data []byte) (uint32, error) {
	switch ChecksumType(checksumType) {
	case ChecksumMurmurhash32:
		return murmurhash32(data), nil
	case ChecksumCrc32:
		return crc32(data), nil
	default:
		log.Error("Unknown checksum type=%v", checksumType)
		return 0, errUnknownChecksumType
	}
}
