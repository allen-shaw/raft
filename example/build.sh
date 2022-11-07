#!/usr/bin/env zsh
#go mod tidy
go mod tidy -compat=1.17
go build -o app *.go
