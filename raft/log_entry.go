package raft

import (
	"bytes"
	"fmt"
	"github.com/AllenShaw19/raft/log"

	"github.com/AllenShaw19/raft/utils"
	"google.golang.org/protobuf/proto"
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
	Peers    []*PeerId
	OldPeers []*PeerId
	Data     bytes.Buffer
}

func NewLogEntry() *LogEntry {
	return &LogEntry{
		Type:     EntryType_ENTRY_TYPE_UNKNOWN,
		Peers:    nil,
		OldPeers: nil,
	}
}

func ParseConfigurationMeta(data *bytes.Buffer, entry *LogEntry) error {
	var meta ConfigurationPBMeta

	err := proto.Unmarshal(data.Bytes(), &meta)
	if err != nil {
		log.Error("meta unmarshal fail %v", err)
		return err
	}
	for _, peer := range meta.Peers {
		entry.Peers = append(entry.Peers, NewPeerId(peer))
	}
	for _, peer := range meta.OldPeers {
		entry.OldPeers = append(entry.OldPeers, NewPeerId(peer))
	}
	return nil
}

func SerializeConfigurationMeta(entry *LogEntry) (*bytes.Buffer, utils.Status) {
	var (
		status utils.Status
		meta   ConfigurationPBMeta
	)

	for _, peer := range entry.Peers {
		meta.Peers = append(meta.Peers, peer.String())
	}
	for _, peer := range entry.OldPeers {
		meta.OldPeers = append(meta.OldPeers, peer.String())
	}
	datas, err := proto.Marshal(&meta)
	data := bytes.NewBuffer(datas)
	if err != nil {
		status.SetError(int32(RaftError_EINVAL), "Fail to serialize ConfigurationPBMeta")
	}
	return data, status
}
