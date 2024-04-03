package main

import (
	"github.com/dig-coins/btcconnect/internal/btcserver"
	"github.com/dig-coins/btcconnect/internal/config"
	"github.com/sgostarter/i/l"
)

func main() {
	logger := l.NewConsoleLoggerWrapper()
	//logger.GetLogger().SetLevel(l.LevelDebug)

	cfg := config.GetBTCManConfig()

	/*
		stg, err := redistorage.NewRedisStorage()
		if err != nil {
			logger.Fatal(err)
		}
	*/

	s := btcserver.NewBTCServer(cfg, nil, logger)

	s.Wait()
}
