package core

import (
	"github.com/AllenShaw19/raft/raft/entity"
	"github.com/AllenShaw19/raft/raft/option"
	"github.com/AllenShaw19/raft/raft/rpc"
)

type ReplicatorGroup struct {
}

func (r ReplicatorGroup) Init(nodeId *entity.NodeId, opts *option.ReplicatorGroupOptions) bool {
	//TODO implement me
	panic("implement me")
}

func (r ReplicatorGroup) AddReplicatorForFollower(peer *entity.PeerId) bool {
	//TODO implement me
	panic("implement me")
}

func (r ReplicatorGroup) AddReplicatorSync(peer *entity.PeerId, replicatorType ReplicatorType) bool {
	//TODO implement me
	panic("implement me")
}

func (r ReplicatorGroup) AddReplicator(peer *entity.PeerId, replicatorType ReplicatorType, sync bool) bool {
	//TODO implement me
	panic("implement me")
}

func (r ReplicatorGroup) SendHeartbeat(peer *entity.PeerId, closure rpc.RpcResponseClosure[interface{}]) {
	//TODO implement me
	panic("implement me")
}
