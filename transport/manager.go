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

	rpcChan          chan raft.RPC
	heartbeatFunc    func(raft.RPC)
	heartbeatFuncMtx sync.Mutex

	connectionMtx sync.Mutex
	connections   map[raft.ServerID]*conn
}

func New(localAddress raft.ServerAddress, dialOptions []grpc.DialOption) *Manager {
	m := &Manager{
		localAddress: localAddress,
		dialOptions:  dialOptions,
		rpcChan:      make(chan raft.RPC),
		connections:  make(map[raft.ServerID]*conn),
	}
	return m
}

func (m *Manager) Register(s *grpc.Server) {
	pb.RegisterRaftServer(s, &rpcService{manager: m})
}

func (m *Manager) Transport() raft.Transport {
	return raftApi{m}
}
