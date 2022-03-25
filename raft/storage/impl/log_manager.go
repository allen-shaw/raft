package impl

import (
	"github.com/AllenShaw19/raft/raft"
	"github.com/AllenShaw19/raft/raft/conf"
	"github.com/AllenShaw19/raft/raft/entity"
	"github.com/AllenShaw19/raft/raft/storage"
	"github.com/AllenShaw19/raft/raft/util"
	"sync"
)

type LogManager struct {
	logStorage     storage.LogStorage
	configManager  *conf.ConfigurationManager
	fsmCaller      raft.FSMCaller
	mutex          sync.Mutex
	stopped        bool // atomic
	hasError       bool // atomic
	nextWaitId     int64
	diskId         *entity.LogId
	appliedId      *entity.LogId
	logsInMemory   *util.SegmentList[*entity.LogEntry]
	firstLogIndex  int64         // atomic
	lastLogIndex   int64         // atomic
	lastSnapshotId *entity.LogId // atomic
}
