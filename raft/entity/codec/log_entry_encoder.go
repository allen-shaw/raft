package codec

import "github.com/AllenShaw19/raft/raft/entity"

type LogEntryEncoder interface {
	Encode(log *entity.LogEntry) []byte
}
