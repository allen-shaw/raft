package raft

import (
	"context"
	"go.uber.org/zap"
)

type Field = zap.Field

type Logger interface {
	Debugf(ctx context.Context, template string, args ...interface{})
	Infof(ctx context.Context, template string, args ...interface{})
	Warnf(ctx context.Context, template string, args ...interface{})
	Errorf(ctx context.Context, template string, args ...interface{})
	Panicf(ctx context.Context, template string, args ...interface{})
	Fatalf(ctx context.Context, template string, args ...interface{})
	Debug(ctx context.Context, msg string, fields ...Field)
	Info(ctx context.Context, msg string, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, fields ...Field)
	Panic(ctx context.Context, msg string, fields ...Field)
	Fatal(ctx context.Context, msg string, fields ...Field)
}
