#!/usr/bin/env bash

export GO111MODULE=on
export GOPROXY=https://mirrors.aliyun.com/goproxy/,https://goproxy.cn,direct
export GOOS=linux

echo 1. install dependencies
go mod download

echo 2. build executable
go build -o beacon

echo 3. stop server
ssh -p 22222 naeemo@comoc.ink "
~/.local/bin/supervisorctl stop beacon
"

echo 4. upload executable
scp -P 22222 ./beacon naeemo@comoc.ink:/opt/apps/beacon

echo 5. start new server
ssh -p 22222 naeemo@comoc.ink "
~/.local/bin/supervisorctl start beacon
"
