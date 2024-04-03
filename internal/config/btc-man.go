package config

import (
	"sync"

	"github.com/dig-coins/btcconnect/internal/btcserver"
	"github.com/sgostarter/libconfig"
)

var (
	_config btcserver.Config
	_once   sync.Once
)

func GetBTCManConfig() *btcserver.Config {
	_once.Do(func() {
		_, err := libconfig.Load("btc-man.yaml", &_config)
		if err != nil {
			panic(err)
		}
	})

	return &_config
}
