package raft

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/AllenShaw19/raft/log"
	"github.com/AllenShaw19/raft/utils"
)

var (
	FlagsRaftSyncSegments bool
)

const (
	RaftSegmentOpenPattern   = "log_inprogress_%020d"
	RaftSegmentClosedPattern = "log_%020d_%020d"
	RaftSegmentMetaFile      = "log_meta"

	EntryHeaderSize uint64 = 24
)

type ChecksumType int

const (
	ChecksumMurmurhash32 ChecksumType = iota
	ChecksumCrc32
)

type entryHeader struct {
	term         int64
	entryType    EntryType
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

// Segment begin here
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

// Create create open segment
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
			err = ParseConfigurationMeta(data, entry)
			if err != nil {
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
			log.Error("data lost in a full segment, path: %s, first_index: %d, expect_last_index: %d, actual_last_index: %d",
				s.path, s.firstIndex, s.lastIndex, actualLastIndex)
			return errors.New("data lost in a full segment")
		} else {
			log.Error("found garbage in a full segment, path: %s, first_index: %d, expect_last_index: %d, actual_last_index: %d",
				s.path, s.firstIndex, s.lastIndex, actualLastIndex)
			return errors.New("found garbage in a full segment")
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
		log.Info("truncate last uncompleted write entry, path: %s, first_index: %d, old_size: %d, new_size: %d",
			s.path, s.firstIndex, fileSize, entryOff)
		err = s.file.Truncate(entryOff)
	}

	_, _ = s.file.Seek(entryOff, io.SeekStart)

	s.bytes = entryOff
	return err
}

func (s *Segment) Append(entry *LogEntry) error {
	if entry == nil || !s.isOpen {
		return errInvalid
	}
	if entry.ID.Index != atomic.LoadInt64(&s.lastIndex)+1 {
		log.Error("entry->index=%d lastIndex=%d firstIndex=%d", entry.ID.Index, s.lastIndex, s.firstIndex)
		return errRange
	}

	data := bytes.NewBuffer(make([]byte, 0))
	switch entry.Type {
	case EntryType_ENTRY_TYPE_DATA:
		{
			n, err := data.Write(entry.Data.Bytes())
			if err != nil {
				log.Error("data write fail, err %v, entry data %v", err, entry.Data.Bytes())
				return err
			}
			if n != entry.Data.Len() {
				log.Error("data write fail, write len=%v, data len=%v", n, entry.Data.Len())
				return errors.New("data write fail")
			}
		}
	case EntryType_ENTRY_TYPE_NO_OP:
	// do nothing
	case EntryType_ENTRY_TYPE_CONFIGURATION:
		{
			var status utils.Status
			data, status = SerializeConfigurationMeta(entry)
			if !status.OK() {
				log.Error("Fail to serialize ConfigurationPBMeta, path: %s", s.path)
				return errors.New(status.Error())
			}
		}
	default:
		log.Fatal("unknown entry type: %v, path %s", entry.Type, s.path)
		return errors.New("unknown entry type")
	}

	if uint64(data.Len()) >= (uint64(1) << 56) {
		log.Error("invalid data len %d, large than max len %d", data.Len(), uint64(1)<<56)
		return errors.New("invalid data len, too large")
	}

	header := bytes.NewBuffer(make([]byte, 0, EntryHeaderSize))
	metaField := uint32(entry.Type)<<24 | uint32(s.checksumType<<16)
	packer := utils.NewRawPacker(header)
	err := packer.Pack64(uint64(entry.ID.Term))
	if err != nil {
		return err
	}
	err = packer.Pack32(metaField)
	if err != nil {
		return err
	}
	err = packer.Pack32(uint32(data.Len()))
	if err != nil {
		return err
	}
	sum, err := getChecksum(s.checksumType, data.Bytes())
	if err != nil {
		return err
	}
	err = packer.Pack32(sum)
	if err != nil {
		return err
	}
	headerBuf := header.Bytes()
	sum, err = getChecksum(s.checksumType, headerBuf[:len(headerBuf)-4])
	if err != nil {
		return err
	}
	err = packer.Pack32(sum)
	if err != nil {
		return err
	}

	toWrite := header.Len() + data.Len()
	n, err := s.file.Write(header.Bytes())
	if err != nil || n != header.Len() {
		log.Error("file [%s] write header fail, err %v, write len %d", s.path, err, n)
		return fmt.Errorf("file write header fail %v", err)
	}
	n, err = s.file.Write(data.Bytes())
	if err != nil {
		log.Error("file [%s] write data fail, err %v, write len %d", s.path, err, n)
		return fmt.Errorf("file write data fail %v", err)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.offsetAndTerm = append(s.offsetAndTerm, offsetAndTerm{s.bytes, entry.ID.Term})
	atomic.AddInt64(&s.lastIndex, 1)
	s.bytes += int64(toWrite)

	return nil
}

func (s *Segment) Get(index int64) (*LogEntry, error) {
	meta, err := s.getMeta(index)
	if err != nil {
		log.Error("get meta fail, err: %v", err)
		return nil, err
	}

	header := &entryHeader{}
	data := bytes.NewBuffer(make([]byte, 0))
	err = s.loadEntry(meta.offset, header, data, meta.length)
	if err != nil {
		log.Error("get meta fail, err: %v", err)
		return nil, err
	}
	if meta.term != header.term {
		log.Error("meta term=%d, header term=%d", meta.term, header.term)
		return nil, errors.New("invalid segment term")
	}
	entry := NewLogEntry()
	switch header.entryType {
	case EntryType_ENTRY_TYPE_DATA:
		entry.Data = *data
	case EntryType_ENTRY_TYPE_NO_OP:
		if data.Len() != 0 {
			log.Error("data of NO_OP must be empty")
			return nil, errors.New("NO_OP data not empty")
		}
	case EntryType_ENTRY_TYPE_CONFIGURATION:
		err := ParseConfigurationMeta(data, entry)
		if err != nil {
			log.Warn("fail to parse ConfigurationPBMeta, path: %s, err %v", s.path, err)
			return nil, err
		}
	default:
		log.Error("unknown entry type %v, path: %s", header.entryType, s.path)
		return nil, errors.New("unknown entry type")
	}
	entry.ID.Index = index
	entry.ID.Term = header.term
	entry.Type = header.entryType

	return entry, nil
}

func (s *Segment) GetTerm(index int64) (int64, error) {
	meta, err := s.getMeta(index)
	if err != nil {
		return 0, err
	}
	return meta.term, nil
}

func (s *Segment) Close(sync bool) error {
	if !s.isOpen {
		log.Fatal("segment is not open")
	}
	oldPath := filepath.Join(s.path, fmt.Sprintf(RaftSegmentOpenPattern, s.firstIndex))
	newPath := filepath.Join(s.path, fmt.Sprintf(RaftSegmentClosedPattern, s.firstIndex, atomic.LoadInt64(&s.lastIndex)))

	log.Info("close a full segment. Current first_index: %v, last_index: %v, raft_sync_segments: %v, will_sync: %v, path: %v",
		s.firstIndex, s.lastIndex, FlagsRaftSyncSegments, sync, newPath)

	err := s.Sync(FlagsRaftSyncSegments && sync)
	if err == nil {
		s.isOpen = false
		err = os.Rename(oldPath, newPath)
		if err != nil {
			log.Error("fail to rename %s to %s", oldPath, newPath)
		} else {
			log.Info("rename %s to %s", oldPath, newPath)
		}
	}
	return err
}

func (s *Segment) Sync(sync bool) error {
	if s.lastIndex > s.firstIndex && sync {
		return s.file.Sync()
	}
	return nil
}

func (s *Segment) RunUnlink() {

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

// TODO: 入参可以优化一下 loadEntry(offset int64, sizeHint uint64) (head *entryHeader, data *bytes.Buffer, err error)
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
		entryType:    EntryType(metaField >> 24),
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

// SegmentMap SegmentLogStorage implement LogStorage
type SegmentMap map[int64]*Segment

type SegmentLogStorage struct {
	path          string
	firstLogIndex int64 // atomic
	lastLogIndex  int64 // atomic
	mutex         sync.Mutex
	segments      SegmentMap
	openedSegment *Segment
	checksumType  ChecksumType
	enableSync    bool
}

func NewSegmentLogStorage(path string, enableSync bool) *SegmentLogStorage {
	sls := &SegmentLogStorage{
		path:          path,
		firstLogIndex: 1,
		lastLogIndex:  0,
		checksumType:  ChecksumMurmurhash32,
		segments:      make(SegmentMap, 0),
		enableSync:    enableSync,
	}
	return sls
}

func (s *SegmentLogStorage) Init(manager *ConfigurationManager) {

}

func (s *SegmentLogStorage) FirstLogIndex() int64 {
	return 0
}

func (s *SegmentLogStorage) LastLogIndex() int64 {
	return 0
}
func (s *SegmentLogStorage) GetEntry(index int64) *LogEntry {
	return nil

}
func (s *SegmentLogStorage) GetTerm(index int64) int64 {
	return 0
}
func (s *SegmentLogStorage) AppendEntry(entry *LogEntry) error {
	return nil

}

func (s *SegmentLogStorage) TruncatePrefix(firstIndexKept int64) error {
	return nil

}

func (s *SegmentLogStorage) TruncateSuffix(lastIndexKept int64) error {
	return nil

}

func (s *SegmentLogStorage) Reset(nextLogIndex int64) error {
	return nil

}

func (s *SegmentLogStorage) NewInstance(uri string) LogStorage {
	return nil
}

func (s *SegmentLogStorage) GCInstance(uri string) error {
	return nil
}

func (s *SegmentLogStorage) Segments() SegmentMap {
	return nil

}

func (s *SegmentLogStorage) ListFiles() []string {

	return nil
}

func (s *SegmentLogStorage) Sync() {

}

func (s *SegmentLogStorage) openSegment() *Segment {

	return nil
}

func (s *SegmentLogStorage) saveMeta(logIndex int64) error {

	return nil
}

func (s *SegmentLogStorage) loadMeta() error {

	return nil
}

func (s *SegmentLogStorage) listSegments(isEmpty bool) error {

	return nil
}

func (s *SegmentLogStorage) loadSegments(manager *ConfigurationManager) error {

	return nil
}

func (s *SegmentLogStorage) getSegment(logIndex int64) (*Segment, error) {

	return nil, nil
}

func (s *SegmentLogStorage) popSegments(firstIndexKept int64) []*Segment {

	return nil
}

// popped []*Segment, lastSegment *Segment
func (s *SegmentLogStorage) popSegmentsFromBack(lastIndexKept int64) ([]*Segment, *Segment) {

	return nil, nil
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
