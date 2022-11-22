package transport

import (
	"context"
	"fmt"
	"github.com/AllenShaw19/raft/raft"
	"github.com/AllenShaw19/raft/transport/codec"
	pb "github.com/AllenShaw19/raft/transport/proto"
	"io"
	"time"
)

type snapshotStream struct {
	s   pb.Raft_InstallSnapshotServer
	buf []byte
}

func (s *snapshotStream) Read(b []byte) (int, error) {
	if len(s.buf) > 0 {
		n := copy(b, s.buf)
		s.buf = s.buf[n:]
		return n, nil
	}
	m, err := s.s.Recv()
	if err != nil {
		return 0, err
	}
	n := copy(b, m.GetData())
	if n < len(m.GetData()) {
		s.buf = m.GetData()[n:]
	}
	return n, nil
}

type rpcService struct {
	pb.UnimplementedRaftServer
	manager *Manager
}

func (s rpcService) handleRPC(groupID string, command interface{}, data io.Reader) (interface{}, error) {
	ch := make(chan raft.RPCResponse, 1)
	rpc := raft.RPC{
		Command:  command,
		RespChan: ch,
		Reader:   data,
	}
	if isHeartbeat(command) {
		// We can take the fast path and use the heartbeat callback and skip the queue in s.manager.rpcChan.
		fn := s.manager.GetHeartbeatFunc(groupID)
		if fn != nil {
			fn(rpc)
			goto wait
		}
	}
	s.manager.GetRPCChan(groupID) <- rpc
wait:
	resp := <-ch
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Response, nil
}

func (s rpcService) AppendEntriesPipeline(server pb.Raft_AppendEntriesPipelineServer) error {
	for {
		msg, err := server.Recv()
		if err != nil {
			return err
		}
		fmt.Printf("[%v][AppendEntriesPipeline] %v\n", time.Now().Format("2006-01-02 15:04:05"), msg.RpcHeader.GroupId)
		resp, err := s.handleRPC(msg.RpcHeader.GroupId, codec.DecodeAppendEntriesRequest(msg), nil)
		if err != nil {
			// TODO(quis): One failure doesn't have to break the entire stream?
			// Or does it all go wrong when it's out of order anyway?
			return err
		}
		if err := server.Send(codec.EncodeAppendEntriesResponse(resp.(*raft.AppendEntriesResponse))); err != nil {
			return err
		}
	}
}

func (s rpcService) AppendEntries(ctx context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	fmt.Printf("[%v][AppendEntries] %v\n", time.Now().Format("2006-01-02 15:04:05"), req.RpcHeader.GroupId)
	resp, err := s.handleRPC(req.RpcHeader.GroupId, codec.DecodeAppendEntriesRequest(req), nil)
	if err != nil {
		return nil, err
	}
	return codec.EncodeAppendEntriesResponse(resp.(*raft.AppendEntriesResponse)), nil
}

func (s rpcService) RequestVote(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	resp, err := s.handleRPC(req.RpcHeader.GroupId, codec.DecodeRequestVoteRequest(req), nil)
	if err != nil {
		return nil, err
	}
	return codec.EncodeRequestVoteResponse(resp.(*raft.RequestVoteResponse)), nil
}

func (s rpcService) TimeoutNow(ctx context.Context, req *pb.TimeoutNowRequest) (*pb.TimeoutNowResponse, error) {
	resp, err := s.handleRPC(req.RpcHeader.GroupId, codec.DecodeTimeoutNowRequest(req), nil)
	if err != nil {
		return nil, err
	}
	return codec.EncodeTimeoutNowResponse(resp.(*raft.TimeoutNowResponse)), nil
}

func (s rpcService) InstallSnapshot(server pb.Raft_InstallSnapshotServer) error {
	isr, err := server.Recv()
	if err != nil {
		return err
	}
	resp, err := s.handleRPC(isr.RpcHeader.GroupId, codec.DecodeInstallSnapshotRequest(isr), &snapshotStream{server, isr.GetData()})
	if err != nil {
		return err
	}
	return server.SendAndClose(codec.EncodeInstallSnapshotResponse(resp.(*raft.InstallSnapshotResponse)))
}

func isHeartbeat(command interface{}) bool {
	req, ok := command.(*raft.AppendEntriesRequest)
	if !ok {
		return false
	}
	return req.Term != 0 && len(req.Leader) != 0 && req.PrevLogEntry == 0 && req.PrevLogTerm == 0 && len(req.Entries) == 0 && req.LeaderCommitIndex == 0
}
