package main

import (
	"context"
	"github.com/AllenShaw19/raft"
	"go.uber.org/zap"
)

type Logger struct {
	l *zap.Logger
}

func NewLogger() *Logger {
	logger := &Logger{}
	var err error
	logger.l, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return logger
}

func (l *Logger) Debugf(ctx context.Context, template string, args ...interface{}) {
	l.l.Sugar().Debugf(template, args)
}

func (l *Logger) Infof(ctx context.Context, template string, args ...interface{}) {
	l.l.Sugar().Infof(template, args)
}

func (l *Logger) Warnf(ctx context.Context, template string, args ...interface{}) {
	l.l.Sugar().Warnf(template, args)
}

func (l *Logger) Errorf(ctx context.Context, template string, args ...interface{}) {
	l.l.Sugar().Errorf(template, args)
}

func (l *Logger) Panicf(ctx context.Context, template string, args ...interface{}) {
	l.l.Sugar().Panicf(template, args)
}

func (l *Logger) Fatalf(ctx context.Context, template string, args ...interface{}) {
	l.l.Sugar().Fatalf(template, args)
}

func (l *Logger) Debug(ctx context.Context, msg string, fields ...raft.Field) {
	l.l.Debug(msg, fields...)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...raft.Field) {
	l.l.Info(msg, fields...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...raft.Field) {
	l.l.Warn(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...raft.Field) {
	l.l.Error(msg, fields...)
}

func (l Logger) Panic(ctx context.Context, msg string, fields ...raft.Field) {
	l.l.Panic(msg, fields...)
}

func (l Logger) Fatal(ctx context.Context, msg string, fields ...raft.Field) {
	l.l.Fatal(msg, fields...)
}
