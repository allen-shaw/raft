package codec

import "github.com/AllenShaw19/raft/raft/entity"

type LogEntryDecoder interface {
	Decode(bs []byte) *entity.LogEntry
}
