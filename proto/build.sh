#!/usr/bin/env bash
#protoc -I=./ --go_out=$GOPATH/src *.proto
protoc -I=. -I=$GOPATH/src -I=$GOPATH/src/github.com/gogo/protobuf --gofast_out=$GOPATH/src raft.proto