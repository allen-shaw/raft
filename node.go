package raft

import (
	"context"
	"fmt"
	"github.com/AllenShaw19/raft/raft"
	"github.com/AllenShaw19/raft/store/boltdb"
	"github.com/AllenShaw19/raft/utils"
	"os"
	"path/filepath"
	"time"
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

func NewNode(groupID string, mgr *NodeManager, logger Logger, svr *Server, fsm raft.FSM) *Node {
	n := &Node{}
	n.groupID = groupID
	n.id = genNodeID(svr.ID, groupID)
	n.logger = logger
	n.mgr = mgr
	n.server = svr
	n.fsm = fsm
	return n
}

func (n *Node) Init(ctx context.Context, peers []string) error {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(n.id)
	config.HeartbeatTimeout = 3 * time.Second
	config.ElectionTimeout = 3 * time.Second

	baseDir := filepath.Join(n.server.RaftDir, n.id)
	isNewNode := !utils.IsPathExists(baseDir)
	if isNewNode {
		err := os.MkdirAll(baseDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// get transport
	transport, err := n.mgr.GetTransport(n.groupID)
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
	if isNewNode && len(peers) == 0 {
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

	for _, peer := range peers {
		err := n.join(peer)
		if err != nil {
			fmt.Printf("%v join to %v fail, err: %v\n", n.id, peer, err)
		}
	}

	return nil
}

func (n *Node) GroupID() string {
	return n.groupID
}

func (n *Node) ID() string {
	return n.id
}

func (n *Node) SetLogger(logger Logger) {
	n.logger = logger
}

func (n *Node) SetMgr(mgr *NodeManager) {
	n.mgr = mgr
}

func (n *Node) IsLeader() bool {
	address, _ := n.raft.LeaderWithID()
	return address == raft.ServerAddress(n.server.Address())
}

func (n *Node) GetLeader() (string, string) {
	addr, id := n.raft.LeaderWithID()
	return string(addr), string(id)
}

func (n *Node) join(peer string) error {
	req := map[string]string{"group_id": n.GroupID(), "node_id": n.ID(), "address": n.server.Address()}
	url := fmt.Sprintf("http://%v/api/join", peer)
	_, err := utils.PostJson(url, req)
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) Apply(cmd []byte) error {
	f := n.raft.Apply(cmd, 3*time.Second)
	return f.Error()
}

func genNodeID(serverID, groupID string) string {
	nodeID := serverID + "/" + groupID
	return nodeID
}
