package entity

import (
	"github.com/AllenShaw19/raft/raft/core"
	"github.com/AllenShaw19/raft/raft/util"
	"strconv"
	"strings"
)

type PeerId struct {
	endpoint *util.Endpoint
	idx      int
	// Node's local priority value, if node don't support priority election, this value is -1
	priority int
	str      string
}

func NewPeerId(endpoint *util.Endpoint, idx, priority int) *PeerId {
	return &PeerId{endpoint: endpoint, idx: idx, priority: priority}
}

func (peer *PeerId) String() string {
	if peer.str == "" {
		var builder strings.Builder
		if peer.idx != 0 {
			builder.WriteString(":")
			builder.WriteString(strconv.Itoa(peer.idx))
		}
		if peer.priority != core.ElectionPriorityDisable {
			if peer.idx == 0 {
				builder.WriteString(":")
			}
			builder.WriteString(":")
			builder.WriteString(strconv.Itoa(peer.priority))
		}

		peer.str = builder.String()
	}
	return peer.str
}

func (peer *PeerId) Copy() *PeerId {
	return NewPeerId(peer.endpoint.Copy(), peer.idx, peer.priority)
}

func (peer *PeerId) Endpoint() *util.Endpoint {
	return peer.endpoint
}

func (peer *PeerId) SetEndpoint(endpoint *util.Endpoint) {
	peer.endpoint = endpoint
}

func (peer *PeerId) Idx() int {
	return peer.idx
}

func (peer *PeerId) SetIdx(idx int) {
	peer.idx = idx
}

func (peer *PeerId) Priority() int {
	return peer.priority
}

func (peer *PeerId) SetPriority(priority int) {
	peer.priority = priority
}
