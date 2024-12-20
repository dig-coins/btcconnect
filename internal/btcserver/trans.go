package btcserver

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/dig-coins/btcconnect/pkg/btctx"
	"github.com/dig-coins/btcconnect/pkg/helper"
	"github.com/okx/go-wallet-sdk/coins/bitcoin"
	"github.com/sgostarter/i/commerr"
	"github.com/sgostarter/libeasygo/cuserror"
	"golang.org/x/exp/slices"
)

type TransInput struct {
	TxID         string
	VOut         uint32
	PrivateKey   string
	Address      string
	Amount       int64
	RedeemScript string
}

type TransOutput struct {
	Address    string `json:"address"`
	Amount     int64  `json:"amount"`
	ChangeFlag bool   `json:"change_flag"`
	Comment    string `json:"comment,omitempty"`
}

func (s *BTCServer) genTransToTxCalcChange(inputs []TransInput, outputs []TransOutput,
	changeAddress string, estimateSmartFee, totalFee int64) (changeAmount int64, err error) {
	var totalAmount, outputAmount int64

	uncompleted := btctx.NewUncompleted()

	var oInputs []btctx.Input

	msPrivateKeys := make(map[int][]string)
	privateKeys := make(map[int]string)

	for idx, input := range inputs {
		oInputs = append(oInputs, btctx.Input{
			TxID:         input.TxID,
			VOut:         input.VOut,
			Address:      input.Address,
			Amount:       input.Amount,
			RedeemScript: input.RedeemScript,
		})

		totalAmount += input.Amount

		i, ok := s.cfg.MultiSignAddressInfos[input.Address]
		if !ok {
			privateKeys[idx] = "cRGtUsda56iwyg7svGUJfcMP3bFLFvyhWpdzoYcKwfuZZeqQoij7"

			continue
		}

		uncompleted.MultiSignInputInfos[idx] = btctx.MultiSignInfo{}

		msPrivateKeys[idx] = make([]string, 0, 2)

		for iFakeKey := 0; iFakeKey < i.MinSignNum; iFakeKey++ {
			msPrivateKeys[idx] = append(msPrivateKeys[idx], "cRGtUsda56iwyg7svGUJfcMP3bFLFvyhWpdzoYcKwfuZZeqQoij7")
		}
	}

	var oOutputs []btctx.Output

	for _, output := range outputs {
		oOutputs = append(oOutputs, btctx.Output{
			Address: output.Address,
			Amount:  output.Amount,
			Comment: output.Comment,
		})

		outputAmount += output.Amount
	}

	if totalAmount <= outputAmount {
		err = errors.New("insufficient")

		return
	}

	if changeAddress == "" {
		changeAddress = "2NC7ACW9obBtsgAgmciLrqJJ36iVcG3Gkgq"
	}

	oOutputs = append(oOutputs, btctx.Output{
		Address: changeAddress,
		Amount:  totalAmount - outputAmount,
	})

	var tx *wire.MsgTx

	if len(uncompleted.MultiSignInputInfos) == 0 {
		tx, err = btctx.GenSignedTx(&btctx.UnsignedTx{
			Inputs:  oInputs,
			Outputs: oOutputs,
		}, s.cfg.GetBTCNetParams())
	} else {
		tx, err = btctx.GenMultiSignedTx(&btctx.UnsignedTx{
			Inputs:  oInputs,
			Outputs: oOutputs,
		}, uncompleted, s.cfg.GetBTCNetParams())
	}

	if err != nil {
		return
	}

	btInputs := make([]bitcoin.Input, 0, len(inputs))
	for _, input := range inputs {
		btInputs = append(btInputs, bitcoin.GenInput(input.TxID, input.VOut,
			"cRGtUsda56iwyg7svGUJfcMP3bFLFvyhWpdzoYcKwfuZZeqQoij7",
			input.RedeemScript, input.Address, input.Amount))
	}

	if len(uncompleted.MultiSignInputInfos) == 0 {
		err = bitcoin.SignBuildTx(tx, btInputs, privateKeys, s.cfg.GetBTCNetParams())
	} else {
		err = bitcoin.MultiSignBuildTx(tx, btInputs, msPrivateKeys, privateKeys, s.cfg.GetBTCNetParams())
	}

	if err != nil {
		return
	}

	size := bitcoin.GetTxVirtualSize2(tx)

	fee := size*estimateSmartFee/1000 + 3
	if totalFee > 0 {
		fee = totalFee
	}

	changeAmount = totalAmount - outputAmount - fee

	if changeAmount < 0 {
		err = errors.New("insufficient")

		return
	}

	return
}

func (s *BTCServer) genTransToTx(inputs []TransInput, outputs []TransOutput) (txHex string, err error) {
	txBuild := bitcoin.NewTxBuild(1, s.cfg.GetBTCNetParams())

	for _, input := range inputs {
		txBuild.AddInput2(input.TxID, input.VOut, input.PrivateKey, input.Address, input.Amount)
	}

	for _, output := range outputs {
		if output.Comment != "" {
			var pkScript []byte

			pkScript, err = txscript.NullDataScript([]byte(output.Comment))
			if err != nil {
				return
			}

			txBuild.AppendTransparentOutput(wire.NewTxOut(int64(0), pkScript))
		} else {
			txBuild.AddOutput(output.Address, output.Amount)
		}
	}

	tx, err := txBuild.Build()
	if err != nil {
		return
	}

	txHex, err = bitcoin.GetTxHex(tx)

	return
}

func (s *BTCServer) WalletPayTo(wallet string, outputs []TransOutput, confirmationTarget int,
	changeAddress string) (txID string, err error) {
	if changeAddress == "" {
		err = commerr.ErrInvalidArgument

		return
	}

	unspentList, err := s.listWalletUnspent(wallet)
	if err != nil {
		return
	}

	estimateSmartFee, err := s.estimateSmartFee(confirmationTarget)
	if err != nil {
		return
	}

	inputs, outputs, err := s.selectUnspentInputs(unspentList, nil, estimateSmartFee, 0,
		outputs, changeAddress, true)
	if err != nil {
		return
	}

	txHex, err := s.genTransToTx(inputs, outputs)
	if err != nil {
		return
	}

	txID, err = s.SendRawTransaction(txHex)
	if err != nil {
		return
	}

	return
}

func (s *BTCServer) GetUnsignedTxEx(inputs []TransInput, outputs []TransOutput) (unsignedTx *btctx.UnsignedTx, unsignedTxHex string, err error) {
	unsignedTx = &btctx.UnsignedTx{
		Inputs:  make([]btctx.Input, 0, len(inputs)),
		Outputs: make([]btctx.Output, 0, len(outputs)),
	}

	for _, input := range inputs {
		unsignedTx.Inputs = append(unsignedTx.Inputs, btctx.Input{
			TxID:         input.TxID,
			VOut:         input.VOut,
			Address:      input.Address,
			Amount:       input.Amount,
			RedeemScript: input.RedeemScript,
		})
	}

	for _, output := range outputs {
		unsignedTx.Outputs = append(unsignedTx.Outputs, btctx.Output{
			Address:    output.Address,
			Amount:     output.Amount,
			ChangeFlag: output.ChangeFlag,
			Comment:    output.Comment,
		})
	}

	unsignedTxHex, err = btctx.MarshalUnsignedTx(unsignedTx)

	return
}

//
//
//

func (s *BTCServer) WalletPayToUnsignedTx(wallet string, payAddresses []string, rawOutputs []TransOutput,
	confirmationTarget int, changeAddress string, minTransFlag bool) (unsignedTxHex string, err error) {
	unspentList, err := s.listWalletUnspent(wallet)
	if err != nil {
		return
	}

	estimateSmartFee, err := s.estimateSmartFee(confirmationTarget)
	if err != nil {
		return
	}

	inputs, outputs, err := s.selectUnspentInputs(unspentList, payAddresses, estimateSmartFee, 0, rawOutputs,
		changeAddress, minTransFlag)
	if err != nil {
		return
	}

	_, unsignedTxHex, err = s.GetUnsignedTxEx(inputs, outputs)

	return
}

func (s *BTCServer) GenUnsignedTx4Gather(wallet string, fromAddresses []string, confirmationTarget int,
	changeAddress string, minTransFlag bool) (wpi string, err error) {
	return s.genUnsignedTx(wallet, fromAddresses, nil, confirmationTarget, changeAddress, minTransFlag)
}

func (s *BTCServer) GenUnsignedTx4TransTo(wallet string, fromAddresses []string, outputs []TransOutput,
	confirmationTarget int, changeAddress string, minTransFlag bool) (wpi string, err error) {
	return s.genUnsignedTx(wallet, fromAddresses, outputs, confirmationTarget, changeAddress, minTransFlag)
}

func (s *BTCServer) GenUnsignedTx4TransToMulti(wallet string, fromAddress string, outputs []TransOutput,
	confirmationTarget int, changeAddress string, minTransFlag bool) (wpi string, err error) {
	return s.genUnsignedTx(wallet, []string{fromAddress}, outputs, confirmationTarget, changeAddress, minTransFlag)
}

func (s *BTCServer) GenUnsignedTx4TransToOne(wallet string, fromAddress, toAddress string, amount int64,
	confirmationTarget int, changeAddress string, minTransFlag bool) (wpi string, err error) {
	return s.genUnsignedTx(wallet, []string{fromAddress}, []TransOutput{{
		Address: toAddress,
		Amount:  amount,
	}}, confirmationTarget, changeAddress, minTransFlag)
}

func (s *BTCServer) genUnsignedTx(wallet string, payAddresses []string, outputs []TransOutput,
	confirmationTarget int, changeAddress string, minTransFlag bool) (wpi string, err error) {
	unspentList, err := s.listWalletUnspent(wallet)
	if err != nil {
		return
	}

	estimateSmartFee, err := s.estimateSmartFee(confirmationTarget)
	if err != nil {
		return
	}

	inputs, outputs, err := s.selectUnspentInputs(unspentList, payAddresses, estimateSmartFee, 0, outputs,
		changeAddress, minTransFlag)
	if err != nil {
		return
	}

	_, wpi, err = s.GetUnsignedTxEx(inputs, outputs)

	return
}

//
//
//

func (s *BTCServer) calcChange4Trans(inputs []TransInput, outputs []TransOutput, changeAddress string,
	estimateSmartFee, totalFee int64) (newOutput []TransOutput, err error) {
	newOutput = append(newOutput, outputs...)

	newOutput = append(newOutput, TransOutput{
		Address: changeAddress,
		Amount:  0,
	})

	changeAmount, err := s.genTransToTxCalcChange(inputs, outputs, changeAddress, estimateSmartFee, totalFee)
	if err != nil {
		return
	}

	newOutput = newOutput[:len(newOutput)-1]

	newOutput = append(newOutput, TransOutput{
		Address:    changeAddress,
		Amount:     changeAmount,
		ChangeFlag: true,
	})

	return
}

func (s *BTCServer) selectUnspentInputs(unspentList []Unspent, payAddresses []string,
	estimateSmartFee, totalFee int64, rawOutputs []TransOutput, changeAddress string, minTransFlag bool) (
	inputs []TransInput, outputs []TransOutput, err error) {
	if changeAddress == "" {
		if len(payAddresses) == 1 && len(rawOutputs) > 0 {
			changeAddress = payAddresses[0]
		}
	}

	var outputAmount int64

	for idx, output := range rawOutputs {
		if output.Comment != "" {
			continue
		}

		if output.Amount <= 0 {
			err = cuserror.NewWithErrorMsg(fmt.Sprintf("invalid amount on input %d", idx))

			return
		}

		outputAmount += output.Amount
	}

	var totalAmount int64

	inputs = make([]TransInput, 0, 2)

	success := false

	for _, unspent := range unspentList {
		if !unspent.Safe {
			continue
		}

		if len(payAddresses) > 0 {
			if !slices.Contains(payAddresses, unspent.Address) {
				continue
			}
		}

		amount := helper.UnitBTC2SatoshiBTC(unspent.Amount)
		inputs = append(inputs, TransInput{
			TxID:    unspent.TxID,
			VOut:    unspent.VOut,
			Address: unspent.Address,
			Amount:  amount,
		})

		i, ok := s.cfg.MultiSignAddressInfos[unspent.Address]
		if ok {
			var redeemScript []byte

			redeemScript, err = bitcoin.GetRedeemScript(i.PublicKeys, i.MinSignNum, s.cfg.GetBTCNetParams())
			if err != nil {
				return
			}

			inputs[len(inputs)-1].RedeemScript = hex.EncodeToString(redeemScript)
		}

		totalAmount += amount

		if changeAddress == "" {
			changeAddress = unspent.Address
		}

		if !minTransFlag || totalAmount <= outputAmount {
			continue
		}

		outputs, err = s.calcChange4Trans(inputs, rawOutputs, changeAddress, estimateSmartFee, totalFee)
		if err == nil {
			success = true

			break
		}
	}

	if totalAmount <= outputAmount {
		err = commerr.ErrResourceExhausted

		return
	}

	if minTransFlag {
		if err == nil && !success {
			err = commerr.ErrResourceExhausted
		}

		return
	}

	if changeAddress == "" {
		err = commerr.ErrInvalidArgument

		return
	}

	outputs, err = s.calcChange4Trans(inputs, rawOutputs, changeAddress, estimateSmartFee, totalFee)

	return
}
