package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/AllenShaw19/raft"
	"github.com/AllenShaw19/raft/builtin"
	"log"
	"os"
	"os/signal"
)

var (
	serverID string
	IP       string
	Port     int
	HttpPort int
	RaftDir  string
	Peer     string
)

func init() {
	flag.StringVar(&serverID, "serverid", "node01", "serverID,node01")
	flag.StringVar(&IP, "ip", "127.0.0.1", "server ip, 127.0.0.1")
	flag.IntVar(&Port, "port", 8081, "port, 8081")
	flag.IntVar(&HttpPort, "httpport", 18081, "http port, 18081")
	flag.StringVar(&RaftDir, "dir", "./data", "raft dir, ./data")
	flag.StringVar(&Peer, "peer", "", "peer, 127.0.0.1:8081")
}

func main() {
	flag.Parse()

	serverID = "node02"
	Port = 8082
	HttpPort = 18082
	Peer = "127.0.0.1:18081"

	server := &raft.Server{
		ID:       serverID,
		IP:       IP,
		Port:     Port,
		HttpPort: HttpPort,
		RaftDir:  RaftDir,
	}
	peerAddrs := make([]string, 0)
	if Peer != "" {
		peerAddrs = append(peerAddrs, Peer)
	}
	ctx := context.Background()

	logger := &Logger{}
	nodeMgr := raft.NewNodeManager(server, logger)
	err := nodeMgr.Init(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("node manager start success")

	nodeConfigs := make([]*raft.NodeConfig, 1)
	for i := range nodeConfigs {
		nodeConfigs[i] = &raft.NodeConfig{
			GroupID:     fmt.Sprintf("group%02d", i),
			PeerAddress: peerAddrs,
		}
	}

	for i, conf := range nodeConfigs {
		err := nodeMgr.NewRaftNode(ctx, conf, NewTable(conf.GroupID))
		if err != nil {
			fmt.Printf("new raft node fail,i=%v, conf=%+v, err:%v \n", i, conf, err)
			panic(err)
		}
	}
	fmt.Println("raft node add success")

	apiService := builtin.NewApiService(nodeMgr, &builtin.Config{Address: server.HttpAddr()}, logger)
	err = apiService.Run(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("api service start success")

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	log.Println("raft exiting")
}
