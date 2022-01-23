package raft

import (
	"bytes"
	"fmt"
	"github.com/AllenShaw19/raft/utils"
)

type LogID struct {
	Index int64
	Term  int64
}

func (id *LogID) String() string {
	return fmt.Sprintf("(index=%d,term=%d)", id.Index, id.Term)
}

func CompareLogID(lhs *LogID, rhs *LogID) int {
	if lhs.Term == rhs.Term {
		if lhs.Index == rhs.Term {
			return 0
		}
		if lhs.Index < rhs.Index {
			return -1
		}
		return +1
	}

	if lhs.Term < rhs.Term {
		return -1
	}
	return +1
}

type LogEntry struct {
	Type     EntryType
	ID       LogID
	Peers    []PeerId
	OldPeers []PeerId
	Data     bytes.Buffer
}

func ParseConfigurationMeta(data *bytes.Buffer, entry *LogEntry) utils.Status {
	var (
		status utils.Status
		meta   ConfigurationPBMeta
	)
	meta.
}

//func SerializeConfigurationMeta(entry *LogEntry, data *bytes.Buffer) utils.Status {
//
//}
