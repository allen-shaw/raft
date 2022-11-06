package raft

import (
	"context"
	"github.com/AllenShaw19/raft/raft"
	"github.com/AllenShaw19/raft/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"sync"
)

type NodeConfig struct {
	GroupID string
}

type NodeManager struct {
	lock   sync.RWMutex
	nodes  map[string]*Node // nodeID, Node
	tm     *transport.Manager
	server *Server
	logger Logger
}

func NewNodeManager(server *Server, logger Logger) *NodeManager {
	m := &NodeManager{}
	m.server = server
	m.logger = logger
	m.nodes = make(map[string]*Node, 0)
	return m
}

func (m *NodeManager) Init(ctx context.Context) error {
	// init transport
	addr := raft.ServerAddress(m.server.Address())
	m.tm = transport.New(addr, []grpc.DialOption{grpc.WithInsecure()})
	go func() {
		err := m.startGRPCService()
		if err != nil {
			panic(err)
		}
	}()

	return nil
}

func (m *NodeManager) NewRaftNode(ctx context.Context, conf *NodeConfig, fsm raft.FSM) error {
	n := NewNode(conf.GroupID, m, m.logger, m.server, fsm)
	err := n.Init(ctx)
	if err != nil {
		return err
	}
	m.nodes[n.ID()] = n
	return nil
}

func (m *NodeManager) GetTransport() (raft.Transport, error) {
	t := m.tm.Transport()
	return t, nil
}

func (m *NodeManager) startGRPCService() error {
	sock, err := net.Listen("tcp", m.server.Address())
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	m.tm.Register(s)
	reflection.Register(s)
	if err := s.Serve(sock); err != nil {
		return err
	}
	return nil
}

func (m *NodeManager) GetNodeInfos() map[string]*Node {
	return m.nodes
}
