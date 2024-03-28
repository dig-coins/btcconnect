package txsigner

import (
	"encoding/hex"
	"encoding/json"

	"github.com/btcsuite/btcd/wire"
	"github.com/dig-coins/btcconnect/internal/btctx"
	"github.com/okx/go-wallet-sdk/coins/bitcoin"
	"github.com/sgostarter/i/commerr"
)

type MiddleSignMidTx struct {
	UnsignedTxHex  string
	UncompletedHex string
	TxHex          string

	unsignedTx  *btctx.UnsignedTx
	uncompleted *btctx.Uncompleted
	tx          *wire.MsgTx
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
