package core

import (
	"github.com/AllenShaw19/raft/raft"
	"github.com/AllenShaw19/raft/raft/conf"
	"github.com/AllenShaw19/raft/raft/entity"
)

// configurationCtx Configuration commit context.
type Stage int

const (
	StageNone       Stage = iota // none stage
	StageCatchingUp              // the node is catching-up
	StageJoint                   // joint stage
	StageStable                  // stable stage
)

type configurationCtx struct {
	node  *Node
	stage Stage

	// Peers change times
	nChanges int
	version  int64
	// peers
	newPeers    []*entity.PeerId
	oldPeers    []*entity.PeerId
	addingPeers []*entity.PeerId

	// learners
	newLearners []*entity.PeerId
	oldLearners []*entity.PeerId
	done        raft.Closure
}

func NewConfigurationCtx(node *Node) *configurationCtx {
	ctx := &configurationCtx{}
	ctx.node = node
	ctx.stage = StageNone
	ctx.version = 0
	ctx.done = nil
	ctx.newPeers = make([]*entity.PeerId, 0)
	ctx.oldPeers = make([]*entity.PeerId, 0)
	ctx.addingPeers = make([]*entity.PeerId, 0)
	ctx.newLearners = make([]*entity.PeerId, 0)
	ctx.oldLearners = make([]*entity.PeerId, 0)
	return ctx
}

// Start change configuration.
func (ctx *configurationCtx) Start(oldConf, newConf *conf.Configuration, done raft.Closure) {

}

func (ctx *configurationCtx) IsBusy() bool {
	return ctx.stage != StageNone
}
