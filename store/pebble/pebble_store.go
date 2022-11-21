package pebble

import (
	"errors"
	"github.com/AllenShaw19/raft/raft"
	"github.com/armon/go-metrics"
	"github.com/cockroachdb/pebble"
	"time"
)

var (
	dbLogsPrefix = []byte("logs")
	dbConfPrefix = []byte("conf")

	ErrKeyNotFound = errors.New("not found")
)

type PebbleStore struct {
	db   *pebble.DB
	path string
}

type Options struct {
	Path          string
	PebbleOptions *pebble.Options
}

func NewPebbleStore(path string) (*PebbleStore, error) {
	return New(Options{Path: path})
}

func New(options Options) (*PebbleStore, error) {
	db, err := pebble.Open("demo", &pebble.Options{})
	if err != nil {
		return nil, err
	}

	store := &PebbleStore{
		db:   db,
		path: options.Path,
	}

	err = store.initialize()
	if err != nil {
		store.Close()
		return nil, err
	}
	return store, nil
}

func (s *PebbleStore) initialize() error {
	return nil
}

func (s *PebbleStore) Close() error {
	return s.db.Close()
}

func (s *PebbleStore) FirstIndex() (uint64, error) {
	o := &pebble.IterOptions{}
	iter := s.db.NewIter(o)
	ok := iter.First()
	if !ok {
		return 0, nil
	}
	first := iter.Key()
	return bytesToUint64(first), nil
}

func (s *PebbleStore) LastIndex() (uint64, error) {
	o := &pebble.IterOptions{}
	iter := s.db.NewIter(o)
	ok := iter.Last()
	if !ok {
		return 0, nil
	}
	last := iter.Key()
	return bytesToUint64(last), nil
}

func (s *PebbleStore) GetLog(idx uint64, log *raft.Log) error {
	defer metrics.MeasureSince([]string{"raft", "pebble", "getLog"}, time.Now())

	key := genLogKey(idx)
	val, closer, err := s.db.Get(key)
	if err == pebble.ErrNotFound {
		return raft.ErrLogNotFound
	}
	if err != nil {
		return err
	}
	defer closer.Close()
	return decodeMsgPack(val, log)
}

func (s *PebbleStore) StoreLog(log *raft.Log) error {
	return s.StoreLogs([]*raft.Log{log})
}

func (s *PebbleStore) StoreLogs(logs []*raft.Log) error {
	now := time.Now()

	b := s.db.NewBatch()
	batchSize := 0

	for _, log := range logs {
		key := genLogKey(log.Index)
		val, err := encodeMsgPack(log)
		if err != nil {
			return err
		}
		logLen := val.Len()
		err = b.Set(key, val.Bytes(), nil)
		if err != nil {
			return err
		}

		batchSize += logLen
		metrics.AddSample([]string{"raft", "pebble", "logSize"}, float32(logLen))
	}

	metrics.AddSample([]string{"raft", "pebble", "logsPerBatch"}, float32(len(logs)))
	metrics.AddSample([]string{"raft", "pebble", "logBatchSize"}, float32(batchSize))

	defer func() {
		metrics.AddSample([]string{"raft", "pebble", "writeCapacity"}, (float32(1_000_000_000)/float32(time.Since(now).Nanoseconds()))*float32(len(logs)))
		metrics.MeasureSince([]string{"raft", "pebble", "storeLogs"}, now)
	}()

	return b.Commit(&pebble.WriteOptions{Sync: true})
}

func (s *PebbleStore) DeleteRange(min, max uint64) error {
	startKey := genLogKey(min)
	endKey := genLogKey(max)
	return s.db.DeleteRange(startKey, endKey, &pebble.WriteOptions{Sync: true})
}

func (s *PebbleStore) Set(key []byte, val []byte) error {
	k := genConfKey(key)
	return s.db.Set(k, val, &pebble.WriteOptions{Sync: true})
}

func (s *PebbleStore) Get(key []byte) ([]byte, error) {
	k := genConfKey(key)
	val, closer, err := s.db.Get(k)
	if err == pebble.ErrNotFound {
		return nil, ErrKeyNotFound
	}
	if err != nil {
		return nil, err
	}
	defer closer.Close()
	return val, nil
}

func (s *PebbleStore) SetUint64(key []byte, val uint64) error {
	return s.Set(key, uint64ToBytes(val))
}

func (s *PebbleStore) GetUint64(key []byte) (uint64, error) {
	val, err := s.Get(key)
	if err != nil {
		return 0, err
	}
	return bytesToUint64(val), nil
}

func (s *PebbleStore) Sync() error {
	return s.db.Flush()
}
