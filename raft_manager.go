package raft

import (
	"context"
	"github.com/AllenShaw19/raft/raft"
	"github.com/AllenShaw19/raft/transport"
	"sync"
)

type NodeManager struct {
	lock  sync.RWMutex
	nodes []*Node
	tm    *transport.Manager
}

func NewNodeManager() *NodeManager {
	m := &NodeManager{}
	m.nodes = make([]*Node, 0)
	return m
}

func (m *NodeManager) Init(ctx context.Context) error {
	//
}

func (m *NodeManager) AddRaft(r *raft.Raft) error {

}

func (m *NodeManager) GetTransport() (raft.Transport, error) {
}
