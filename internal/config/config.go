package config

import (
	"github.com/btcsuite/btcd/chaincfg"
	hdwallet "github.com/dig-coins/hd-wallet"
)

func GetBTCNetParams(coinType CoinType) *chaincfg.Params {
	switch coinType {
	case CoinTypeBTCTestnet:
		return &chaincfg.TestNet3Params
	case CoinTypeBTC:
		fallthrough
	default:
		return &chaincfg.MainNetParams
	}
}

func GetHDWalletCoinType(coinType CoinType) uint32 {
	switch coinType {
	case CoinTypeBTCTestnet:
		return hdwallet.BTCTestnet
	case CoinTypeBTC:
		return hdwallet.BTC
	default:
		return 0
	}
}
