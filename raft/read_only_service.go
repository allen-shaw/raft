package raft

import (
	"github.com/AllenShaw19/raft/raft/option"
)

type ReadOnlyService interface {
	Init(opts *option.ReadOnlyServiceOptions) bool
}
