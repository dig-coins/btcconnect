package btcserver

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/okx/go-wallet-sdk/coins/bitcoin"
	"github.com/sgostarter/i/commerr"
)

func (s *BTCServer) GenMultiSignatureAddress(pubKeys []string, minSignNum int, net *chaincfg.Params) (multiAddress string, err error) {
	if minSignNum > len(pubKeys) {
		err = commerr.ErrInvalidArgument

		return
	}

	redeemScript, err := bitcoin.GetRedeemScript(pubKeys, minSignNum, net)
	if err != nil {
		return
	}

	multiAddress, err = bitcoin.GenerateMultiAddress(redeemScript, s.cfg.GetBTCNetParams())

	return
}
