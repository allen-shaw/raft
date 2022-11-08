package raft

import (
	"context"
	"fmt"
	"github.com/AllenShaw19/raft/raft"
	"github.com/AllenShaw19/raft/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"net"
	"sync"
)

type NodeConfig struct {
	GroupID     string
	PeerAddress []string
}

type nodeMap struct {
	m    map[string]map[string]*Node
	lock sync.RWMutex
}

func newNodeMap() *nodeMap {
	nm := &nodeMap{}
	nm.m = make(map[string]map[string]*Node)
	return nm
}

func (m *nodeMap) Set(n *Node) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.m[n.GroupID()]; !ok {
		m.m[n.GroupID()] = make(map[string]*Node)
	}

	m.m[n.GroupID()][n.ID()] = n
}

func (m *nodeMap) Get(groupID string) map[string]*Node {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.m[groupID]
}

func (m *nodeMap) Dump() map[string]*Node {
	m.lock.RLock()
	defer m.lock.RUnlock()

	r := make(map[string]*Node)

	for _, subM := range m.m {
		for id, n := range subM {
			r[id] = n
		}
	}
	return r
}

type NodeManager struct {
	lock   sync.RWMutex
	nodes  *nodeMap // map<GroupID, map<NodeID,Node>>
	tm     *transport.Manager
	server *Server
	logger Logger
}

func NewNodeManager(server *Server, logger Logger) *NodeManager {
	m := &NodeManager{}
	m.server = server
	m.logger = logger
	m.nodes = newNodeMap()
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
	err := n.Init(ctx, conf.PeerAddress)
	if err != nil {
		return err
	}
	m.nodes.Set(n)
	return nil
}

func (m *NodeManager) GetTransport(groupID string) (raft.Transport, error) {
	t := m.tm.Transport(groupID)
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
	return m.nodes.Dump()
}

func (m *NodeManager) Join(ctx context.Context, nodeID, groupID, address string) error {
	// FIXME: 加锁
	nodes := m.nodes.Get(groupID)

	m.logger.Debugf(ctx, "[NodeManager] Join, nodeID %v, groupID %v, address %v, myNodesLen %v", nodeID, groupID, address, len(nodes))

	for _, node := range nodes {
		f := node.raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(address), 0, 0)
		if f.Error() != nil {
			return f.Error()
		}
	}

	fmt.Printf("node %v, group %v, address %v\n", nodeID, groupID, address)

	return nil
}

func (m *NodeManager) GetNode(groupID string) *Node {
	nodes := m.nodes.Get(groupID)
	for _, n := range nodes {
		return n
	}
	return nil
}
