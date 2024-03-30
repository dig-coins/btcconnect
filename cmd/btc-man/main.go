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
		RPCHost:     cfg.RPCConfig.RPCHost,
		RPCPort:     cfg.RPCConfig.RPCPort,
		RPCUser:     cfg.RPCConfig.RPCUser,
		RPCPassword: cfg.RPCConfig.RPCPassword,
		RPCUseSSL:   cfg.RPCConfig.RPCUseSSL,
	}, stg, cfg.GetBTCNetParams(), cfg.MultiSignAddressInfos, logger)

	s.Wait()
}
