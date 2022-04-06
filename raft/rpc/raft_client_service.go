package rpc

import "github.com/AllenShaw19/raft/raft/option"

type RaftClientService interface {
	Init(options *option.NodeOptions) bool
}
