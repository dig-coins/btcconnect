package config

import (
	"sync"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/dig-coins/btcconnect/internal/btcserver"
	"github.com/dig-coins/btcconnect/internal/share"
	"github.com/sgostarter/libconfig"
)

type CoinType int

const (
	CoinTypeStart      CoinType = 1
	CoinTypeBTCTestnet          = iota
	CoinTypeBTC
)

type BTCManConfig struct {
	CoinType              CoinType                               `json:"coin_type" yaml:"coin_type"`
	RpcConfig             btcserver.Config                       `json:"rpc_config" yaml:"rpc_config"`
	MultiSignAddressInfos map[string]*share.MultiSignAddressInfo `json:"multi_sign_address_infos" yaml:"multi_sign_address_infos"`
}

func (cfg *BTCManConfig) GetBTCNetParams() *chaincfg.Params {
	switch cfg.CoinType {
	case CoinTypeBTCTestnet:
		return &chaincfg.TestNet3Params
	case CoinTypeBTC:
		fallthrough
	default:
		return &chaincfg.MainNetParams
	}
}

var (
	_config BTCManConfig
	_once   sync.Once
)

func GetBTCManConfig() *BTCManConfig {
	_once.Do(func() {
		_, err := libconfig.Load("btc-man.yaml", &_config)
		if err != nil {
			panic(err)
		}
	})

	return &_config
}
