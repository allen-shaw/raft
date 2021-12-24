package entity

import (
	"github.com/AllenShaw19/raft/core"
	"github.com/AllenShaw19/raft/util"
)

type PeerId struct {
	Idx      int
	Addr     util.EndPoint
	Priority core.ElectionPriority
}
