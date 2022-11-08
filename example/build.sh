#!/usr/bin/env zsh
rm -rf data
go mod tidy -compat=1.17
go build -o app *.go
