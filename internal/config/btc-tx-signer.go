package config

import (
	"sync"

	"github.com/dig-coins/btcconnect/internal/share"
	"github.com/sgostarter/libconfig"
)

type BTCTxSignerConfig struct {
	Listens               string                                 `json:"listens" yaml:"listens"`
	CoinType              share.CoinType                         `json:"coin_type" yaml:"coin_type"`
	SeedFileName          string                                 `json:"seed_file_name" yaml:"seed_file_name"`
	SeedSecKey            string                                 `json:"seed_sec_key" yaml:"seed_sec_key"`
	MultiSignAddressInfos map[string]*share.MultiSignAddressInfo `json:"multi_sign_address_infos" yaml:"multi_sign_address_infos"`
	AutoUnsignedTxRoot    string                                 `json:"auto_unsigned_tx_root" yaml:"auto_unsigned_tx_root"`
}

var (
	_btcTxSignerConfig BTCTxSignerConfig
	_btcTxSignerOnce   sync.Once
)

func GetBTCTxSignerConfig(configFile ...string) *BTCTxSignerConfig {
	_btcTxSignerOnce.Do(func() {
		f := "btc-tx-signer.yaml"
		if len(configFile) == 1 && configFile[0] != "" {
			f = configFile[0]
		}
		_, err := libconfig.Load(f, &_btcTxSignerConfig)
		if err != nil {
			panic(err)
		}
	})

	return &_btcTxSignerConfig
}
