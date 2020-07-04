#!/usr/bin/env bash

#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
#CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
env GOOS=windows GOARCH=386 go build -ldflags "-s -w" -o xuexi.exe main.go