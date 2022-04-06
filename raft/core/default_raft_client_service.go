package core

import (
	"github.com/AllenShaw19/raft/raft"
	"github.com/AllenShaw19/raft/raft/option"
)

type DefaultRaftClientService struct {
	rgGroup raft.ReplicatorGroup
}

func (s *DefaultRaftClientService) Init(options *option.NodeOptions) bool {
	//TODO implement me
	panic("implement me")
}

func NewDefaultRaftClientService(rgGroup raft.ReplicatorGroup) *DefaultRaftClientService {
	s := &DefaultRaftClientService{}
	s.rgGroup = rgGroup
	return s
}
