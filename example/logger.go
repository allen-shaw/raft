package main

import (
	"context"
	"fmt"
	"github.com/AllenShaw19/raft"
	"log"
	"os"
)

type Logger struct {
}

func (l *Logger) Debugf(ctx context.Context, template string, args ...interface{}) {
	log.Printf(template, args)
}

func (l Logger) Infof(ctx context.Context, template string, args ...interface{}) {
	log.Printf(template, args)
}

func (l Logger) Warnf(ctx context.Context, template string, args ...interface{}) {
	log.Printf(template, args)
}

func (l Logger) Errorf(ctx context.Context, template string, args ...interface{}) {
	log.Printf(template, args)
}

func (l Logger) Panicf(ctx context.Context, template string, args ...interface{}) {
	log.Printf(template, args)
	panic(fmt.Sprintf(template, args))
}

func (l Logger) Fatalf(ctx context.Context, template string, args ...interface{}) {
	log.Printf(template, args)
	os.Exit(255)
}

func (l Logger) Debug(ctx context.Context, msg string, fields ...raft.Field) {
	log.Println(msg, fields)
}

func (l Logger) Info(ctx context.Context, msg string, fields ...raft.Field) {
	log.Println(msg, fields)
}

func (l Logger) Warn(ctx context.Context, msg string, fields ...raft.Field) {
	log.Println(msg, fields)
}

func (l Logger) Error(ctx context.Context, msg string, fields ...raft.Field) {
	log.Println(msg, fields)
}

func (l Logger) Panic(ctx context.Context, msg string, fields ...raft.Field) {
	log.Println(msg, fields)
	panic(msg)
}

func (l Logger) Fatal(ctx context.Context, msg string, fields ...raft.Field) {
	log.Println(msg, fields)
	os.Exit(255)
}
