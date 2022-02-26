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
	"time"

	"github.com/AllenShaw19/raft/log"
	"github.com/AllenShaw19/raft/utils"
)

var (
	FlagsRaftSyncSegments            bool
	FlagsRaftMaxSegmentSize          int64
	FlagsRaftSync                    bool
	FlagsRaftSyncMeta                bool
	FlagsRaftTraceAppendEntryLatency bool
)

const (
	SegmentOpenPattern   = "log_inprogress_%020d"
	SegmentClosedPattern = "log_%020d_%020d"
	SegmentMetaFile      = "log_meta"

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
	checksumType ChecksumType
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
	checksumType  ChecksumType
	offsetAndTerm []offsetAndTerm
}

func NewSegment(path string, firstIndex, lastIndex int64, checksumType ChecksumType) *Segment {
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
	path := filepath.Join(s.path, fmt.Sprintf(SegmentOpenPattern, s.firstIndex))
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
		path = filepath.Join(path, fmt.Sprintf(SegmentOpenPattern, s.firstIndex))
	} else {
		path = filepath.Join(path, fmt.Sprintf(SegmentClosedPattern, s.firstIndex, atomic.LoadInt64(&s.lastIndex)))
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
	headerBuf := header.Bytes()

	packer := utils.NewRawPacker(header)
	packer.Pack64(uint64(entry.ID.Term))
	packer.Pack32(metaField)
	packer.Pack32(uint32(data.Len()))
	packer.Pack32(getChecksum(s.checksumType, data.Bytes()))
	packer.Pack32(getChecksum(s.checksumType, headerBuf[:len(headerBuf)-4]))

	if err := packer.Error(); err != nil {
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
	oldPath := filepath.Join(s.path, fmt.Sprintf(SegmentOpenPattern, s.firstIndex))
	newPath := filepath.Join(s.path, fmt.Sprintf(SegmentClosedPattern, s.firstIndex, atomic.LoadInt64(&s.lastIndex)))

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

func (s *Segment) runUnlink(filepath string) {
	start := time.Now()
	err := os.Remove(filepath)
	durationMs := time.Since(start).Milliseconds()
	if err != nil {
		log.Error("unlink %s err %v time: %d ms", filepath, err, durationMs)
		return
	}
	log.Info("unlink %s success, time: %d ms", filepath, durationMs)
}

func (s *Segment) Unlink() error {
	path := s.path
	if s.isOpen {
		path = filepath.Join(path, fmt.Sprintf(SegmentOpenPattern, s.firstIndex))
	} else {
		path = filepath.Join(path, fmt.Sprintf(SegmentClosedPattern, s.firstIndex, atomic.LoadInt64(&s.lastIndex)))
	}
	tmpPath := path + ".tmp"
	err := os.Rename(path, tmpPath)
	if err != nil {
		log.Error("Fail to rename %s to %s", path, tmpPath)
		return err
	}
	// start goroutine to unlink
	go s.runUnlink(tmpPath)
	log.Info("Unlinked segment %s", path)
	return nil
}

func (s *Segment) Truncate(lastIndexKept int64) error {
	s.mutex.Lock()
	if lastIndexKept >= s.lastIndex {
		s.mutex.Unlock()
		return nil
	}

	firstTruncateInOffset := lastIndexKept + 1 - s.firstIndex
	truncateSize := s.offsetAndTerm[firstTruncateInOffset].offset
	log.Info("Truncating %s, firstIndex %v lastIndex from %v to %v, truncate size to %v",
		s.path, s.firstIndex, s.lastIndex, lastIndexKept, truncateSize)
	s.mutex.Unlock()

	if !s.isOpen {
		oldPath := filepath.Join(s.path, fmt.Sprintf(SegmentClosedPattern, s.firstIndex, atomic.LoadInt64(&s.lastIndex)))
		newPath := filepath.Join(s.path, fmt.Sprintf(SegmentOpenPattern, s.firstIndex))
		err := os.Rename(oldPath, newPath)
		if err != nil {
			log.Error("Fail to rename %s to %s, err %v", oldPath, newPath, err)
			return err
		}
		log.Info("Renamed %s to %s success", oldPath, newPath)
		s.isOpen = true
	}

	err := s.file.Truncate(truncateSize)
	if err != nil {
		log.Error("file truncate fail, err %v", err)
		return err
	}

	retOff, err := s.file.Seek(truncateSize, io.SeekStart)
	if err != nil {
		log.Error("fail to seek file path %s to size=%v, err %v, ret %d", s.path, truncateSize, err, retOff)
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.offsetAndTerm = s.offsetAndTerm[:firstTruncateInOffset]
	atomic.StoreInt64(&s.lastIndex, lastIndexKept)
	s.bytes = truncateSize
	return nil
}

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
		return fmt.Sprintf(SegmentClosedPattern, s.firstIndex, atomic.LoadInt64(&s.lastIndex))
	}
	return fmt.Sprintf(SegmentOpenPattern, s.firstIndex)
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
		checksumType: ChecksumType((metaField << 8) >> 24),
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

func (s *SegmentLogStorage) Init(manager *ConfigurationManager) error {
	if FlagsRaftMaxSegmentSize < 0 {
		log.Fatal("FlagsRaftMaxSegmentSize %d must be greater than or equal to 0", FlagsRaftMaxSegmentSize)
		return errors.New("invalid raft max segment size")
	}
	dirPath := filepath.Dir(s.path)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Error("Fail to create %s : %v", dirPath, err)
		return err
	}

	// default use crc32 checksum
	s.checksumType = ChecksumCrc32
	log.Info("Use crc32c as the checksum type of appending entries")

	isEmpty := false

	defer func() {
		if isEmpty {
			atomic.StoreInt64(&s.firstLogIndex, 1)
			atomic.StoreInt64(&s.lastLogIndex, 0)
			err = s.saveMeta(1)
		}
	}()

	err = s.loadMeta()
	if err != nil && err == errNoEntry {
		log.Warn("%s is empty", s.path)
		isEmpty = true
	} else if err != nil {
		return err
	}

	err = s.listSegments(isEmpty)
	if err != nil {
		return err
	}

	err = s.loadSegments(manager)
	if err != nil {
		return err
	}

	return err
}

func (s *SegmentLogStorage) FirstLogIndex() int64 {
	return atomic.LoadInt64(&s.firstLogIndex)
}

func (s *SegmentLogStorage) LastLogIndex() int64 {
	return atomic.LoadInt64(&s.lastLogIndex)
}

func (s *SegmentLogStorage) GetEntry(index int64) (*LogEntry, error) {
	segment, err := s.getSegment(index)
	if err != nil {
		log.Error("get segment fail, index %d, err %v", index, err)
		return nil, err
	}
	return segment.Get(index)
}

func (s *SegmentLogStorage) GetTerm(index int64) (int64, error) {
	segment, err := s.getSegment(index)
	if err != nil {
		log.Error("get segment fail, index %d, err %v", index, err)
		return 0, err
	}
	return segment.GetTerm(index)
}
func (s *SegmentLogStorage) AppendEntry(entry *LogEntry) error {
	segment := s.openSegment()
	if segment == nil {
		return errIO
	}
	err := segment.Append(entry)
	if err != nil /*&& !errors.Is(err, EEXIST)*/ { // FIXME:处理文件已存在的错误
		return err
	}

	term, err := s.GetTerm(entry.ID.Index)
	if err != nil {
		return err
	}
	if /** errors.Is(err, EEXIST) && */ entry.ID.Term != term {
		return errInvalid
	}
	atomic.AddInt64(&s.lastLogIndex, 1)
	return segment.Sync(s.enableSync)
}

func (s *SegmentLogStorage) AppendEntries(entries []*LogEntry, metric *IOMetric) (int, error) {
	if len(entries) == 0 {
		return 0, nil
	}
	if lastLogIndex := atomic.LoadInt64(&s.lastLogIndex); lastLogIndex+1 != entries[0].ID.Index {
		log.Fatal("There's gap between appending entries and last_log_index, path %s", s.path)
		return -1, errors.New("")
	}

	var lastSegment *Segment
	for i, entry := range entries {
		segment := s.openSegment()
		if FlagsRaftTraceAppendEntryLatency && metric != nil {
			//	TODO: 增加metric
		}
		if segment == nil {
			return i, errors.New("nil segment")
		}
		err := segment.Append(entry)
		if err != nil {
			return i, err
		}
		if FlagsRaftTraceAppendEntryLatency && metric != nil {
			//	TODO: 增加metric
		}
		atomic.AddInt64(&s.lastLogIndex, 1)
		lastSegment = segment
	}

	err := lastSegment.Sync(s.enableSync)
	if err != nil {
		return len(entries), err
	}

	if FlagsRaftTraceAppendEntryLatency && metric != nil {
		//	TODO: 增加metric
	}

	return len(entries), nil
}

func (s *SegmentLogStorage) TruncatePrefix(firstIndexKept int64) error {
	if atomic.LoadInt64(&s.firstLogIndex) >= firstIndexKept {
		log.Info("nothing is going to happen since first_log_index=%v >= first_index_kept=%v", s.firstLogIndex, firstIndexKept)
		return nil
	}
	if err := s.saveMeta(firstIndexKept); err != nil {
		log.Error("Fail to save meta, path: %s", s.path)
		return err
	}
	poppedSegments := s.popSegments(firstIndexKept)
	for i := range poppedSegments {
		poppedSegments[i].Unlink()
		poppedSegments[i] = nil
	}
	return nil
}

func (s *SegmentLogStorage) TruncateSuffix(lastIndexKept int64) error {
	popped, lastSegment := s.popSegmentsFromBack(lastIndexKept)

	truncateLastSegment := false
	var err error
	if lastSegment != nil {
		if s.firstLogIndex <= s.lastLogIndex {
			truncateLastSegment = true
		} else {
			s.mutex.Lock()
			popped = append(popped, lastSegment)
			delete(s.segments, lastSegment.FirstIndex())
			if s.openedSegment != nil {
				if s.openedSegment != lastSegment {
					log.Fatal("Check failed: openedSegment != lastSegment")
				}
				s.openedSegment = nil
			}
			s.mutex.Unlock()
		}
	}

	for i := range popped {
		err = popped[i].Unlink()
		if err != nil {
			return err
		}
		popped[i] = nil
	}

	if truncateLastSegment {
		closed := !lastSegment.IsOpen()
		err = lastSegment.Truncate(lastIndexKept)
		if err == nil && closed && lastSegment.IsOpen() {
			s.mutex.Lock()
			if s.openedSegment != nil {
				log.Fatal("check failed：openedSegment not nil")
			}
			delete(s.segments, lastSegment.FirstIndex())
			s.openedSegment, lastSegment = lastSegment, s.openedSegment
			s.mutex.Unlock()
		}
	}

	return err
}

func (s *SegmentLogStorage) Reset(nextLogIndex int64) error {
	if nextLogIndex <= 0 {
		log.Error("invalid next_log_index=%d, path: %s", nextLogIndex, s.path)
		return errInvalid
	}

	popped := make([]*Segment, 0, len(s.segments))
	s.mutex.Lock()
	for _, segment := range s.segments {
		popped = append(popped, segment)
	}
	s.segments = make(SegmentMap)
	if s.openedSegment != nil {
		popped = append(popped, s.openedSegment)
		s.openedSegment = nil
	}
	atomic.StoreInt64(&s.firstLogIndex, nextLogIndex)
	atomic.StoreInt64(&s.lastLogIndex, nextLogIndex-1)
	s.mutex.Unlock()

	if err := s.saveMeta(nextLogIndex); err != nil {
		log.Error("Fail to save meta, path: %s", s.path)
		return err
	}

	for i := range popped {
		popped[i].Unlink()
		popped[i] = nil
	}
	return nil
}

func (s *SegmentLogStorage) NewInstance(uri string) LogStorage {
	return NewSegmentLogStorage(uri, true)
}

func (s *SegmentLogStorage) GCInstance(uri string) utils.Status {
	status := utils.Status{}
	if err := gcDir(uri); err != nil {
		log.Error("Failed to gc log storage from path %s, err %v", uri, err)
		status.SetError(int32(RaftError_EINVAL), "Failed to gc log storage")
		return status
	}
	log.Info("Succeed to gc log storage from path %s", uri)
	return status
}

func (s *SegmentLogStorage) Segments() SegmentMap {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.segments
}

func (s *SegmentLogStorage) ListFiles() []string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	segFiles := make([]string, 0, len(s.segments)+1)
	for _, segment := range s.segments {
		segFiles = append(segFiles, segment.FileName())
	}
	if s.openedSegment != nil {
		segFiles = append(segFiles, s.openedSegment.FileName())
	}

	return segFiles
}

func (s *SegmentLogStorage) Sync() {
	s.mutex.Lock()
	segments := make([]*Segment, 0, len(s.segments))
	for _, segment := range s.segments {
		segments = append(segments, segment)
	}
	s.mutex.Unlock()

	for _, segment := range segments {
		segment.Sync(true)
	}
}

func (s *SegmentLogStorage) openSegment() *Segment {
	var prevOpenSegment *Segment

	s.mutex.Lock()
	if s.openedSegment == nil {
		s.openedSegment = NewSegment(s.path, s.LastLogIndex()+1, s.LastLogIndex(), s.checksumType)
		if err := s.openedSegment.Create(); err != nil {
			s.openedSegment = nil
			s.mutex.Unlock()
			return nil
		}
	}
	if s.openedSegment.Bytes() > FlagsRaftMaxSegmentSize {
		s.segments[s.openedSegment.FirstIndex()] = s.openedSegment
		prevOpenSegment, s.openedSegment = s.openedSegment, prevOpenSegment // swap
	}
	s.mutex.Unlock()

	if prevOpenSegment != nil {
		if err := prevOpenSegment.Close(s.enableSync); err == nil {
			s.mutex.Lock()
			s.openedSegment = NewSegment(s.path, s.LastLogIndex()+1, s.LastLogIndex(), s.checksumType)
			if err := s.openedSegment.Create(); err == nil {
				// success
				s.mutex.Unlock()
				return s.openedSegment
			}
			s.mutex.Unlock()
		}
		log.Error("Fail to close old open_segment or create new open_segment, path: %s", s.path)
		// Failed, revert former changes
		s.mutex.Lock()
		delete(s.segments, prevOpenSegment.FirstIndex())
		s.openedSegment, prevOpenSegment = prevOpenSegment, s.openedSegment
		s.mutex.Unlock()
		return nil
	}
	return s.openedSegment
}

func (s *SegmentLogStorage) saveMeta(logIndex int64) error {
	start := time.Now()
	metaPath := filepath.Join(s.path, SegmentMetaFile)

	meta := &LogPBMeta{FirstLogIndex: logIndex}
	pbFile := ProtoBufFile{path: metaPath}

	err := pbFile.Save(meta, syncMeta())
	if err != nil {
		log.Error("Fail to save meta to %s", metaPath)
	}
	log.Info("log save_meta %s first_log_index: %v, time: %dms", metaPath, logIndex, time.Since(start).Milliseconds())
	return err
}

func (s *SegmentLogStorage) loadMeta() error {
	start := time.Now()
	metaPath := filepath.Join(s.path, SegmentMetaFile)

	meta := &LogPBMeta{}
	pbFile := ProtoBufFile{path: metaPath}
	err := pbFile.Load(meta)
	if err != nil {
		log.Error("Fail to load meta from %s", metaPath)
		return err
	}
	atomic.StoreInt64(&s.firstLogIndex, meta.FirstLogIndex)

	log.Info("log load_meta %s first_log_index: %v, time: %dms", metaPath, meta.FirstLogIndex, time.Since(start).Milliseconds())
	return nil
}

func (s *SegmentLogStorage) listSegments(isEmpty bool) error {
	// TODO: to implements
	return nil
}

func (s *SegmentLogStorage) loadSegments(manager *ConfigurationManager) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, segment := range s.segments {
		log.Info("load closed segment, path: %s, first index: %v, last index: %v",
			s.path, segment.FirstIndex(), segment.LastIndex())
		err := segment.Load(manager)
		if err != nil {
			return err
		}
		atomic.StoreInt64(&s.lastLogIndex, segment.LastIndex())
	}

	if s.openedSegment != nil {
		log.Info("oad open segment, path: %s, first index %v", s.path, s.openedSegment.FirstIndex())
		err := s.openedSegment.Load(manager)
		if err != nil {
			return nil
		}
		if atomic.LoadInt64(&s.firstLogIndex) > s.openedSegment.LastIndex() {
			log.Error("open segment need discard, path: %s, first_log_index %v, first_index %v, last_index %v",
				s.path, atomic.LoadInt64(&s.firstLogIndex), s.openedSegment.FirstIndex(), s.openedSegment.LastIndex())
			s.openedSegment.Unlink()
			s.openedSegment = nil
		} else {
			atomic.StoreInt64(&s.lastLogIndex, s.openSegment().LastIndex())
		}
	}

	if s.lastLogIndex == 0 {
		s.lastLogIndex = s.firstLogIndex - 1
	}

	return nil
}

func (s *SegmentLogStorage) getSegment(index int64) (*Segment, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	firstIndex := s.FirstLogIndex()
	lastIndex := s.LastLogIndex()

	if firstIndex == lastIndex+1 {
		return nil, errors.New("no segments")
	}

	if index < firstIndex || index > lastIndex+1 {
		if index > lastIndex {
			log.Error("Attempted to access entry %v outside of log, first_log_index %v, last_log_index %v",
				index, firstIndex, lastIndex)
		}
		return nil, errors.New("no segments")
	} else if index == lastIndex+1 {
		return nil, errors.New("no segments")
	}

	if s.openedSegment != nil && index >= s.openedSegment.FirstIndex() {
		return s.openedSegment, nil
	} else {
		log.CHECK(len(s.segments) != 0, "!s.segment.empty()")

		// FIXME: 需要map 有序，支持upper_bound(index)
		//return
	}

	return nil, nil
}

func (s *SegmentLogStorage) popSegments(firstIndexKept int64) []*Segment {
	popped := make([]*Segment, 0)
	s.mutex.Lock()
	defer s.mutex.Unlock()

	atomic.StoreInt64(&s.firstLogIndex, firstIndexKept)
	for key, segment := range s.segments {
		if segment.LastIndex() < firstIndexKept {
			popped = append(popped, segment)
			delete(s.segments, key)
		} else {
			return popped
		}
	}

	if s.openedSegment != nil {
		if s.openedSegment.LastIndex() < firstIndexKept {
			popped = append(popped, s.openedSegment)
			s.openedSegment = nil
			atomic.StoreInt64(&s.lastLogIndex, firstIndexKept-1)
		} else {
			log.CHECK(s.openedSegment.FirstIndex() <= firstIndexKept, "s.openedSegment.FirstIndex() <= firstIndexKept")
		}
	} else {
		atomic.StoreInt64(&s.lastLogIndex, firstIndexKept-1)
	}

	return popped
}

// popped []*Segment, lastSegment *Segment
func (s *SegmentLogStorage) popSegmentsFromBack(lastIndexKept int64) ([]*Segment, *Segment) {
	popped := make([]*Segment, 0)
	var lastSegment *Segment
	s.mutex.Lock()
	defer s.mutex.Unlock()

	atomic.StoreInt64(&s.lastLogIndex, lastIndexKept)
	if s.openedSegment != nil {
		if s.openedSegment.FirstIndex() <= lastIndexKept {
			lastSegment = s.openedSegment
			return popped, lastSegment
		}
		popped = append(popped, s.openedSegment)
		s.openedSegment = nil
	}

	// FIXME: segmentMap需要有序，而且支持反向查找，segmentMap需要重写
	//for key, segment := range s.segments {
	//	if segment.FirstIndex() <= lastIndexKept {
	//		break
	//	}
	//	popped = append(popped, segment)
	//	delete(s.segments, key)
	//}
	return nil, nil
}

// util function
func verifyChecksum(checksumType ChecksumType, data []byte, value uint32) bool {
	switch checksumType {
	case ChecksumMurmurhash32:
		return value == murmurhash32(data)
	case ChecksumCrc32:
		return value == crc32(data)
	default:
		log.Error("Unknown checksum type=%v", checksumType)
		return false
	}
}

func getChecksum(checksumType ChecksumType, data []byte) uint32 {
	switch checksumType {
	case ChecksumMurmurhash32:
		return murmurhash32(data)
	case ChecksumCrc32:
		return crc32(data)
	default:
		log.Error("Unknown checksum type=%v", checksumType)
		return 0
	}
}

func syncMeta() bool {
	return FlagsRaftSync || FlagsRaftSyncMeta
}
