package main

import (
	"context"
	"fmt"
	"github.com/AllenShaw19/raft"
	"github.com/AllenShaw19/raft/builtin"
	"log"
	"os"
	"os/signal"
)

func main() {
	server := &raft.Server{
		ID:      "node01",
		IP:      "127.0.0.1",
		Port:    8081,
		RaftDir: "./data",
	}
	ctx := context.Background()

	logger := &Logger{}
	nodeMgr := raft.NewNodeManager(server, logger)
	err := nodeMgr.Init(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("node manager start success")

	nodeConfigs := make([]*raft.NodeConfig, 2)
	for i := range nodeConfigs {
		nodeConfigs[i] = &raft.NodeConfig{
			GroupID: fmt.Sprintf("group%02d", i),
		}
	}

	for i, conf := range nodeConfigs {
		err := nodeMgr.NewRaftNode(ctx, conf)
		if err != nil {
			fmt.Printf("new raft node fail,i=%v, conf=%+v, err:%v \n", i, conf, err)
			panic(err)
		}
	}
	fmt.Println("raft node add success")

	apiService := builtin.NewApiService(nodeMgr, &builtin.Config{Address: "127.0.0.1:19948"}, logger)
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
