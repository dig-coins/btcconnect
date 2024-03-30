package txsigner

import (
	"bytes"
	"encoding/hex"
	"encoding/json"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/dig-coins/btcconnect/internal/btctx"
	"github.com/okx/go-wallet-sdk/coins/bitcoin"
	"github.com/sgostarter/i/commerr"
)

type MiddleSignMidTx struct {
	UnsignedTxHex  string
	UncompletedHex string
	TxHex          string

	unsignedTx    *btctx.UnsignedTx
	uncompleted   *btctx.Uncompleted
	tx            *wire.MsgTx
	bitcoinInputs []bitcoin.Input
}

func (msTx *MiddleSignMidTx) Check(netParams *chaincfg.Params) bool {
	if msTx.unsignedTx == nil || msTx.tx == nil || msTx.uncompleted == nil {
		return false
	}

	if len(msTx.unsignedTx.Inputs) != len(msTx.tx.TxIn) || len(msTx.unsignedTx.Outputs) != len(msTx.tx.TxOut) {
		return false
	}

	inputs := make([]bitcoin.Input, 0, len(msTx.unsignedTx.Inputs))

	for idx, input := range msTx.unsignedTx.Inputs {
		if input.TxID != msTx.tx.TxIn[idx].PreviousOutPoint.Hash.String() ||
			input.VOut != msTx.tx.TxIn[idx].PreviousOutPoint.Index {
			return false
		}

		inputs = append(inputs, bitcoin.GenInput(input.TxID, input.VOut, "",
			input.RedeemScript, input.Address, input.Amount))
	}

	for idx, output := range msTx.unsignedTx.Outputs {
		if msTx.tx.TxOut[idx].Value != output.Amount {
			return false
		}

		pkScript, err := bitcoin.AddrToPkScript(output.Address, netParams)
		if err != nil {
			return false
		}

		if !bytes.Equal(msTx.tx.TxOut[idx].PkScript, pkScript) {
			return false
		}
	}

	msTx.bitcoinInputs = inputs

	return true
}

func (msTx *MiddleSignMidTx) GetUnsignedTx() *btctx.UnsignedTx {
	return msTx.unsignedTx
}

func MarshalMiddleSignMidTx(msTx *MiddleSignMidTx) (hexD string, err error) {
	if msTx == nil || msTx.unsignedTx == nil || msTx.uncompleted == nil || msTx.tx == nil {
		err = commerr.ErrInvalidArgument

		return
	}

	msTx.UnsignedTxHex, err = btctx.MarshalUnsignedTx(msTx.unsignedTx)
	if err != nil {
		return
	}

	msTx.UncompletedHex, err = btctx.MarshalUncompleted(msTx.uncompleted)
	if err != nil {
		return
	}

	msTx.TxHex, err = bitcoin.GetTxHex(msTx.tx)
	if err != nil {
		return
	}

	d, err := json.Marshal(msTx)
	if err != nil {
		return
	}

	hexD = hex.EncodeToString(d)

	return
}

func UnmarshalMiddleSignMidTx(hexD string) (msTx *MiddleSignMidTx, err error) {
	d, err := hex.DecodeString(hexD)
	if err != nil {
		return
	}

	msTx = &MiddleSignMidTx{}

	err = json.Unmarshal(d, msTx)
	if err != nil {
		return
	}

	if msTx.UnsignedTxHex == "" || msTx.UncompletedHex == "" || msTx.TxHex == "" {
		err = commerr.ErrInvalidArgument

		return
	}

	msTx.unsignedTx, err = btctx.UnmarshalUnsignedTx(msTx.UnsignedTxHex)
	if err != nil {
		return
	}

	msTx.uncompleted, err = btctx.UnmarshalUncompleted(msTx.UncompletedHex)
	if err != nil {
		return
	}

	msTx.tx, err = bitcoin.NewTxFromHex(msTx.TxHex)
	if err != nil {
		return
	}

	return
}
