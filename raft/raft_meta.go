package raft

import (
	"github.com/cockroachdb/pebble"
	"sync"
)

// MixedMetaStorage start
type MixedMetaStorage struct {
}

// FileBasedSingleMetaStorage start
type FileBasedSingleMetaStorage struct {
}

// KVBasedMergedMetaStorage start
type KVBasedMergedMetaStorage struct {
}

// WriteTask kvBasedMergedMetaStorageImpl start
type WriteTask struct {
	Term     int64
	VotedFor PeerId
	VgID     VersionedGroupId
	Done     Closure
}

type KVBasedMergedMetaStorageImpl struct {
	mutex    sync.Mutex
	isInited bool
	path     string
	db       *pebble.DB
}

func NewKVBasedMergedMetaStorageImpl(path string) *KVBasedMergedMetaStorageImpl {
	return &KVBasedMergedMetaStorageImpl{
		path:     path,
		isInited: false,
	}
}
func (s *KVBasedMergedMetaStorageImpl) Init() error {

	// todo implement me

	return nil
}
