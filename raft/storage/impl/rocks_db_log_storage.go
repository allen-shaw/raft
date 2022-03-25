package impl

import (
	"github.com/AllenShaw19/raft/raft/entity/codec"
	"github.com/AllenShaw19/raft/raft/util"
	"github.com/cockroachdb/pebble"
	"sync"
)

type RocksDBLogStorage struct {
	//    private final List<ColumnFamilyOptions> cfOptions     = new ArrayList<>();
	//    private ColumnFamilyHandle              defaultHandle;
	//    private ColumnFamilyHandle              confHandle;
	//    private ReadOptions                     totalOrderReadOptions;
	//    private DebugStatistics                 statistics;
	//    private final ReadWriteLock             readWriteLock = new ReentrantReadWriteLock();
	//    private final Lock                      readLock      = this.readWriteLock.readLock();
	//    private final Lock                      writeLock     = this.readWriteLock.writeLock();
	path           string
	sync           bool
	openStatistics bool
	db             *pebble.DB
	dbOptions      *pebble.Options
	writeOptions   *pebble.WriteOptions
	statistics     *util.DebugStatistics
	mutex          sync.Mutex

	firstLogIndex        int64 // atomic
	hasLoadFirstLogIndex bool  // atomic

	LogEntryEncoder codec.LogEntryEncoder
	LogEntryDecoder codec.LogEntryDecoder
}
