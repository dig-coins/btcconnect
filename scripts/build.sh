#!/bin/bash

#
# https://github.com/golang/go/issues/61229
#


dest=./dest
rm -rf $dest
mkdir -p $dest

go build -ldflags=-linkmode=internal -o $dest/seeds-data cmd/seeds-data/main.go
go build -ldflags=-linkmode=internal -o $dest/btc-man cmd/btc-man/main.go


GOARCH=amd64 GOOS=windows go build -ldflags "-s -w" -o $dest/btc-man.exe cmd/btc-man/main.go
GOARCH=arm64 GOOS=windows go build -ldflags "-s -w" -o $dest/btc-man.arm64.exe cmd/btc-man/main.go

GOARCH=amd64 GOOS=linux go build -ldflags "-s -w" -o $dest/btc-man.linux cmd/btc-man/main.go

GOARCH=arm64 GOOS=darwin go build -ldflags "-s -w" -o $dest/btc-man.macos.arm64 cmd/btc-man/main.go
GOARCH=amd64 GOOS=darwin go build -ldflags "-s -w" -o $dest/btc-man.macos.amd64 cmd/btc-man/main.go
