package raft

import (
	"context"
	"github.com/AllenShaw19/raft/raft"
	"github.com/AllenShaw19/raft/store/boltdb"
	"github.com/AllenShaw19/raft/utils"
	"os"
	"path/filepath"
)

type Config struct {
	RaftDir string
}

type Node struct {
	id     string
	raft   *raft.Raft
	logger Logger
	mgr    *NodeManager
	conf   *Config
	fsm    raft.FSM
}

func NewNode(id string, mgr *NodeManager, logger Logger, conf *Config) *Node {
	n := &Node{}
	n.id = id
	n.logger = logger
	n.mgr = mgr
	n.conf = conf
	return n
}

func (n *Node) Init(ctx context.Context) error {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(n.id)

	baseDir := filepath.Join(n.conf.RaftDir, n.id)
	isNewNode := !utils.IsPathExists(baseDir)

	// get transport
	transport, err := n.mgr.GetTransport()
	if err != nil {
		return err
	}

	// new log store
	ldb, err := boltdb.NewBoltStore(filepath.Join(baseDir, "logs.dat"))
	if err != nil {
		return err
	}

	// new stable store
	sdb, err := boltdb.NewBoltStore(filepath.Join(baseDir, "stable.dat"))
	if err != nil {
		return err
	}

	// new snapshot store
	fss, err := raft.NewFileSnapshotStore(baseDir, 3, os.Stderr)
	if err != nil {
		return err
	}

	// new raft node
	r, err := raft.NewRaft(config, n.fsm, ldb, sdb, fss, transport)
	if err != nil {
		return err
	}
	n.raft = r

	// bootstrap for new node
	if isNewNode {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		f := n.raft.BootstrapCluster(configuration)
		if err := f.Error(); err != nil {
			return err
		}
	}

	return nil
}

func (n *Node) SetLogger(logger Logger) {
	n.logger = logger
}

func (n *Node) SetMgr(mgr *NodeManager) {
	n.mgr = mgr
}
