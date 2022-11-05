package raft

import (
	"context"
	"github.com/AllenShaw19/raft/raft"
	"github.com/AllenShaw19/raft/store/boltdb"
	"github.com/AllenShaw19/raft/utils"
	"os"
	"path/filepath"
)

type Node struct {
	id      string
	groupID string
	server  *Server
	raft    *raft.Raft
	logger  Logger
	mgr     *NodeManager
	fsm     raft.FSM
}

func NewNode(groupID string, mgr *NodeManager, logger Logger, svr *Server) *Node {
	n := &Node{}
	n.groupID = groupID
	n.id = genNodeID(svr.ID, groupID)
	n.logger = logger
	n.mgr = mgr
	n.server = svr
	return n
}

func (n *Node) Init(ctx context.Context) error {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(n.id)

	baseDir := filepath.Join(n.server.RaftDir, n.id)
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

func genNodeID(serverID, groupID string) string {
	nodeID := serverID + "/" + groupID
	return nodeID
}
