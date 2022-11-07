package builtin

import (
	"context"
	"github.com/AllenShaw19/raft"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Config struct {
	Address string
}

type ApiService struct {
	conf   *Config
	e      *gin.Engine
	logger raft.Logger
	mgr    *raft.NodeManager
}

func NewApiService(mgr *raft.NodeManager, conf *Config, logger raft.Logger) *ApiService {
	gin.SetMode(gin.ReleaseMode)

	s := &ApiService{}
	s.mgr = mgr
	s.conf = conf
	s.logger = logger
	s.e = gin.Default()
	// set route
	s.init()

	return s
}

func (s *ApiService) Run(ctx context.Context) error {
	var err error
	go func() {
		err = s.e.Run(s.conf.Address)
		if err != nil {
			s.logger.Error(ctx, "run api service fail", zap.Any("addr", s.conf.Address), zap.Error(err))
		}
	}()

	return err
}

func (s *ApiService) init() {
	s.e.Use(gin.Logger())
	s.e.Use(gin.Recovery())

	api := s.e.Group("/api")
	api.GET("/leader", s.getLeader)
	api.GET("/status", s.getStatus)
	api.POST("/join", s.handleJoin)
}

func (s *ApiService) getLeader(ctx *gin.Context) {

}

func (s *ApiService) getStatus(ctx *gin.Context) {
	status := &Status{Nodes: make([]*NodeInfo, 0)}

	nodes := s.mgr.GetNodeInfos()
	for _, node := range nodes {
		info := &NodeInfo{
			ID:       node.ID(),
			IsLeader: node.IsLeader(),
		}
		if !info.IsLeader {
			leaderAddr, leaderID := node.GetLeader()
			info.Leader = &Leader{ID: leaderID, Address: leaderAddr}
		}
		status.Nodes = append(status.Nodes, info)
	}

	resp := &Response{Code: 0, Msg: "ok", Data: status}
	ctx.JSON(http.StatusOK, resp)
}

func (s *ApiService) handleJoin(ctx *gin.Context) {
	req := &JoinRequest{}
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	err = s.mgr.Join(ctx, req.NodeID, req.GroupID, req.Address)
	if err != nil {
		resp := &Response{Code: 10001, Msg: err.Error(), Data: nil}
		ctx.JSON(http.StatusOK, resp)
	}

	resp := &Response{Code: 0, Msg: "ok", Data: nil}
	ctx.JSON(http.StatusOK, resp)
}
