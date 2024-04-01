package btcserver

import (
	"github.com/dig-coins/btcconnect/pkg/helper"
	"github.com/sgostarter/i/commerr"
	"github.com/spf13/cast"
)

func (s *BTCServer) SendRawTransaction(hexTx string) (txID string, err error) {
	err = s.rpcClient.CallWrapper("sendrawtransaction", []any{hexTx}, &txID)

	return
}

type Unspent struct {
	TxID          string  `json:"txid"`
	VOut          uint32  `json:"vout"`
	Label         string  `json:"label"`
	Address       string  `json:"address"`
	ScriptPubKey  string  `json:"scriptPubKey"`
	Safe          bool    `json:"safe"`
	Confirmations int     `json:"confirmations"`
	Spendable     bool    `json:"spendable"`
	Amount        float64 `json:"amount"`
	Desc          string  `json:"desc"`
}

// nolint: unused
func (s *BTCServer) listUnspent() (unspent []Unspent, err error) {
	err = s.rpcClient.CallWrapper("listunspent", []any{nil}, &unspent)

	return
}

func (s *BTCServer) listWalletUnspent(wallet string) (unspent []Unspent, err error) {
	err = s.rpcClient.CallWalletWrapper(wallet, "listunspent", []any{nil}, &unspent)

	return
}

func (s *BTCServer) estimateSmartFee(confirmationTarget int) (fee int64, err error) {
	var r map[string]any

	err = s.rpcClient.CallWrapper("estimatesmartfee", []any{confirmationTarget, "ECONOMICAL"}, &r)
	if err != nil {
		return
	}

	fee = helper.UnitBTC2SatoshiBTC(cast.ToFloat64(r["feerate"]))
	if fee == 0 {
		err = commerr.ErrInternal

		return
	}

	return
}
