package raft

import (
	"github.com/AllenShaw19/raft/entity"
	"github.com/AllenShaw19/raft/options"
	"github.com/AllenShaw19/raft/rpc"
)

type GroupService struct {
	started     bool // should atomic load
	serverID    *entity.PeerID
	nodeOptions *options.NodeOptions
	groupID     string
	node        Node
	rpcServer   rpc.Server
}

func (s *GroupService) Start() {

}