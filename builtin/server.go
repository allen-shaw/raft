package builtin

import (
	"context"
	"github.com/AllenShaw19/raft"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Config struct {
	address string
}

type ApiService struct {
	conf   *Config
	e      *gin.Engine
	logger raft.Logger
	mgr    *raft.RaftManager
}

func NewApiService(mgr *raft.RaftManager, conf *Config, logger raft.Logger) *ApiService {
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
	err := s.e.Run(s.conf.address)
	if err != nil {
		s.logger.Error(ctx, "run api service fail", zap.Any("addr", s.conf.address), zap.Error(err))
		return err
	}
	return nil
}

func (s *ApiService) init() {
	s.e.Use(gin.Logger())
	s.e.Use(gin.Recovery())

	api := s.e.Group("/api")
	api.GET("/leader", s.getLeader)
}

func (s *ApiService) getLeader(ctx *gin.Context) {

}
