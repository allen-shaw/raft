package rpc

import "github.com/AllenShaw19/raft/raft"

type RpcResponseClosure[T any] interface {
	raft.Closure
	// SetResponse Called by request handler to set response.
	// @param resp rpc response
	SetResponse(resp T)
}
