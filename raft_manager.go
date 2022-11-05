package raft

import (
	"context"
	"github.com/AllenShaw19/raft/raft"
	"github.com/AllenShaw19/raft/transport"
	"sync"
)

type NodeManager struct {
	lock   sync.RWMutex
	nodes  []*Node
	tm     *transport.Manager
	server *Server
}

func NewNodeManager() *NodeManager {
	m := &NodeManager{}
	m.nodes = make([]*Node, 0)
	return m
}

func (m *NodeManager) Init(ctx context.Context) error {
	//

	return nil
}

func (m *NodeManager) AddRaft(r *raft.Raft) error {

	return nil
}

func (m *NodeManager) GetTransport() (raft.Transport, error) {
	t := m.tm.Transport()
	return t, nil
}
