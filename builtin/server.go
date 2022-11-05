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
}

func NewApiService() *ApiService {
	s := &ApiService{}
	s.e = gin.Default()
	// set route

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
