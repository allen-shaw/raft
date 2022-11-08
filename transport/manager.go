package transport

import (
	"github.com/AllenShaw19/raft/raft"
	pb "github.com/AllenShaw19/raft/transport/proto"
	"google.golang.org/grpc"
	"sync"
)

type Manager struct {
	localAddress raft.ServerAddress
	dialOptions  []grpc.DialOption

	connectionMtx sync.Mutex
	connections   map[raft.ServerID]*conn

	apis map[string]*raftApi // map<groupID,api>
}

func New(localAddress raft.ServerAddress, dialOptions []grpc.DialOption) *Manager {
	m := &Manager{
		localAddress: localAddress,
		dialOptions:  dialOptions,
		connections:  make(map[raft.ServerID]*conn),
		apis:         make(map[string]*raftApi),
	}
	return m
}

func (m *Manager) Register(s *grpc.Server) {
	pb.RegisterRaftServer(s, &rpcService{manager: m})
}

func (m *Manager) Transport(groupID string) raft.Transport {
	if t, ok := m.apis[groupID]; ok {
		return t
	}

	m.apis[groupID] = newRaftApi(groupID, m)
	return m.apis[groupID]
}

func (m *Manager) GetRPCChan(groupID string) chan raft.RPC {
	return m.apis[groupID].rpcChan
}

func (m *Manager) GetHeartbeatFunc(groupID string) func(rpc raft.RPC) {
	return m.apis[groupID].GetHeartbeatFunc()
}
