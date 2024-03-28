package btctx

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/okx/go-wallet-sdk/coins/bitcoin"
	"github.com/sgostarter/i/commerr"
)

func GenSignedTx(unsignedTx *UnsignedTx, netParams *chaincfg.Params) (tx *wire.MsgTx, err error) {
	if unsignedTx == nil || len(unsignedTx.Inputs) == 0 || len(unsignedTx.Outputs) == 0 {
		err = commerr.ErrInvalidArgument

		return
	}

	txBuild := bitcoin.NewTxBuild(1, netParams)

	for _, input := range unsignedTx.Inputs {
		txBuild.AddInput2(input.TxID, input.VOut, "", input.Address, input.Amount)
	}

	for _, output := range unsignedTx.Outputs {
		txBuild.AddOutput(output.Address, output.Amount)
	}

	tx, err = txBuild.Build2()

	return
}

func GenMultiSignedTx(unsignedTx *UnsignedTx, uncompleted *Uncompleted,
	netParams *chaincfg.Params) (tx *wire.MsgTx, err error) {
	if unsignedTx == nil || len(unsignedTx.Inputs) == 0 || len(unsignedTx.Outputs) == 0 {
		err = commerr.ErrInvalidArgument

		return
	}

	txBuild := bitcoin.NewTxBuild(1, netParams)

	for idx, input := range unsignedTx.Inputs {
		if _, ok := uncompleted.MultiSignInputInfos[idx]; ok {
			txBuild.AddInput(input.TxID, input.VOut, "", input.RedeemScript, "", input.Amount)
		} else {
			txBuild.AddInput2(input.TxID, input.VOut, "", input.Address, input.Amount)
		}
	}

	for _, output := range unsignedTx.Outputs {
		txBuild.AddOutput(output.Address, output.Amount)
	}

	hexTx, err := txBuild.SingleBuild2()
	if err != nil {
		return
	}

	tx, err = bitcoin.NewTxFromHex(hexTx)

	return
}
