package main

import (
	"github.com/dig-coins/btcconnect/internal/btcserver"
	"github.com/dig-coins/btcconnect/internal/config"
	"github.com/dig-coins/btcconnect/internal/redistorage"
	"github.com/sgostarter/i/l"
)

func main() {
	logger := l.NewConsoleLoggerWrapper()
	//logger.GetLogger().SetLevel(l.LevelDebug)

	cfg := config.GetBTCManConfig()

	stg, err := redistorage.NewRedisStorage()
	if err != nil {
		logger.Fatal(err)
	}

	s := btcserver.NewBTCServer(&btcserver.Config{
		RPCHost:     cfg.RpcConfig.RPCHost,
		RPCPort:     cfg.RpcConfig.RPCPort,
		RPCUser:     cfg.RpcConfig.RPCUser,
		RPCPassword: cfg.RpcConfig.RPCPassword,
		RPCUseSSL:   cfg.RpcConfig.RPCUseSSL,
	}, stg, cfg.GetBTCNetParams(), cfg.MultiSignAddressInfos, logger)

	s.Wait()
}
