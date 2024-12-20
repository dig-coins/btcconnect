package share

import (
	"github.com/btcsuite/btcd/chaincfg"
	hdwallet "github.com/dig-coins/hd-wallets"
)

type CoinType int

const (
	CoinTypeStart      CoinType = 1
	CoinTypeBTCTestnet          = iota
	CoinTypeBTC
	CoinTypeBTCRegtest
)

func GetBTCNetParams(coinType CoinType) *chaincfg.Params {
	switch coinType {
	case CoinTypeBTCRegtest:
		return &chaincfg.RegressionNetParams
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
	case CoinTypeBTCRegtest:
		return hdwallet.BTCRegTest
	case CoinTypeBTCTestnet:
		return hdwallet.BTCTestnet
	case CoinTypeBTC:
		return hdwallet.BTC
	default:
		return 0
	}
}
