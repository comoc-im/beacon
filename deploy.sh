#!/usr/bin/env bash

export GOPROXY=https://mirrors.aliyun.com/goproxy/,https://goproxy.cn,direct
export GOOS=linux
export GOARCH=386

echo 1. install dependencies
go mod download

echo 2. build executable
go build -o beacon ./cmd/beacon/main.go

echo 3. stop server
ssh root@chummy.fun "
supervisorctl stop beacon
"

echo 4. upload executable
scp ./beacon root@chummy.fun:/opt/apps/beacon

echo 5. start new server
ssh root@chummy.fun "
supervisorctl start beacon
"
