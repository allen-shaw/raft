package codec

import (
	"github.com/AllenShaw19/raft/raft"
	pb "github.com/AllenShaw19/raft/transport/proto"
	"github.com/AllenShaw19/raft/utils"
)

func encodeRPCHeader(s raft.RPCHeader) *pb.RPCHeader {
	//fmt.Printf("[encodeRPCHeader] id:%v addr:%v \n", string(s.ID), string(s.Addr))
	_, groupID := utils.ParseNodeID(string(s.ID))
	return &pb.RPCHeader{
		GroupId:         groupID,
		ProtocolVersion: int64(s.ProtocolVersion),
	}
}

func encodeLogs(s []*raft.Log) []*pb.Log {
	ret := make([]*pb.Log, len(s))
	for i, l := range s {
		ret[i] = encodeLog(l)
	}
	return ret
}

func encodeLog(s *raft.Log) *pb.Log {
	return &pb.Log{
		Index:      s.Index,
		Term:       s.Term,
		Type:       encodeLogType(s.Type),
		Data:       s.Data,
		Extensions: s.Extensions,
	}
}

func encodeLogType(s raft.LogType) pb.Log_LogType {
	switch s {
	case raft.LogCommand:
		return pb.Log_LOG_COMMAND
	case raft.LogNoop:
		return pb.Log_LOG_NOOP
	case raft.LogAddPeerDeprecated:
		return pb.Log_LOG_ADD_PEER_DEPRECATED
	case raft.LogRemovePeerDeprecated:
		return pb.Log_LOG_REMOVE_PEER_DEPRECATED
	case raft.LogBarrier:
		return pb.Log_LOG_BARRIER
	case raft.LogConfiguration:
		return pb.Log_LOG_CONFIGURATION
	default:
		panic("invalid LogType")
	}
}

func EncodeAppendEntriesRequest(s *raft.AppendEntriesRequest) *pb.AppendEntriesRequest {
	return &pb.AppendEntriesRequest{
		RpcHeader:         encodeRPCHeader(s.RPCHeader),
		Term:              s.Term,
		Leader:            s.Leader,
		PrevLogEntry:      s.PrevLogEntry,
		PrevLogTerm:       s.PrevLogTerm,
		Entries:           encodeLogs(s.Entries),
		LeaderCommitIndex: s.LeaderCommitIndex,
	}
}

func EncodeAppendEntriesResponse(s *raft.AppendEntriesResponse) *pb.AppendEntriesResponse {
	return &pb.AppendEntriesResponse{
		RpcHeader:      encodeRPCHeader(s.RPCHeader),
		Term:           s.Term,
		LastLog:        s.LastLog,
		Success:        s.Success,
		NoRetryBackoff: s.NoRetryBackoff,
	}
}

func EncodeRequestVoteRequest(s *raft.RequestVoteRequest) *pb.RequestVoteRequest {
	return &pb.RequestVoteRequest{
		RpcHeader:          encodeRPCHeader(s.RPCHeader),
		Term:               s.Term,
		Candidate:          s.Candidate,
		LastLogIndex:       s.LastLogIndex,
		LastLogTerm:        s.LastLogTerm,
		LeadershipTransfer: s.LeadershipTransfer,
	}
}

func EncodeRequestVoteResponse(s *raft.RequestVoteResponse) *pb.RequestVoteResponse {
	return &pb.RequestVoteResponse{
		RpcHeader: encodeRPCHeader(s.RPCHeader),
		Term:      s.Term,
		Peers:     s.Peers,
		Granted:   s.Granted,
	}
}

func EncodeInstallSnapshotRequest(s *raft.InstallSnapshotRequest) *pb.InstallSnapshotRequest {
	return &pb.InstallSnapshotRequest{
		RpcHeader:          encodeRPCHeader(s.RPCHeader),
		Term:               s.Term,
		Leader:             s.Leader,
		LastLogIndex:       s.LastLogIndex,
		LastLogTerm:        s.LastLogTerm,
		Peers:              s.Peers,
		Configuration:      s.Configuration,
		ConfigurationIndex: s.ConfigurationIndex,
		Size:               s.Size,
	}
}

func EncodeInstallSnapshotResponse(s *raft.InstallSnapshotResponse) *pb.InstallSnapshotResponse {
	return &pb.InstallSnapshotResponse{
		RpcHeader: encodeRPCHeader(s.RPCHeader),
		Term:      s.Term,
		Success:   s.Success,
	}
}

func EncodeTimeoutNowRequest(s *raft.TimeoutNowRequest) *pb.TimeoutNowRequest {
	return &pb.TimeoutNowRequest{
		RpcHeader: encodeRPCHeader(s.RPCHeader),
	}
}

func EncodeTimeoutNowResponse(s *raft.TimeoutNowResponse) *pb.TimeoutNowResponse {
	return &pb.TimeoutNowResponse{
		RpcHeader: encodeRPCHeader(s.RPCHeader),
	}
}
