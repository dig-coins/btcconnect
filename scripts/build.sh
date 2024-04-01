#!/bin/bash

#
# https://github.com/golang/go/issues/61229
#


dest=./dest
rm -rf $dest
mkdir -p $dest

go build -ldflags=-linkmode=internal -o $dest/seeds-data cmd/seeds-data/main.go
go build -ldflags=-linkmode=internal -o $dest/btc-tx-signer cmd/btc-tx-signer/main.go



