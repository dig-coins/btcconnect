package txsigner

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/dig-coins/btcconnect/internal/btctx"
	"github.com/dig-coins/btcconnect/internal/edfile"
	"github.com/dig-coins/btcconnect/internal/keypool"
	"github.com/dig-coins/btcconnect/internal/share"
	"github.com/okx/go-wallet-sdk/coins/bitcoin"
	"github.com/sgostarter/i/commerr"
	"github.com/sgostarter/libeasygo/cuserror"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

type TxSigner struct {
	netParams            *chaincfg.Params
	keyPool              *keypool.KeyPool
	multiSignAddressInfo map[string]*share.MultiSignAddressInfo
}

func NewTxSigner(coinType uint32, netParams *chaincfg.Params, seedFileName, seedSecKey string,
	multiSignAddressInfo map[string]*share.MultiSignAddressInfo) (signer *TxSigner, err error) {
	if coinType == 0 || netParams == nil {
		err = commerr.ErrInvalidArgument

		return
	}

	signer = &TxSigner{
		netParams:            netParams,
		multiSignAddressInfo: multiSignAddressInfo,
	}

	err = signer.init(seedFileName, seedSecKey, coinType)
	if err != nil {
		return
	}

	return
}

func (signer *TxSigner) init(seedFileName, seedSecKey string, coinType uint32) (err error) {
	d, err := edfile.ReadSecFile(seedFileName, seedSecKey)
	if err != nil {
		return
	}

	var seeds []string

	err = yaml.Unmarshal(d, &seeds)
	if err != nil {
		return
	}

	if len(seeds) == 0 {
		err = cuserror.NewWithErrorMsg("no seeds")

		return
	}

	signer.keyPool, err = keypool.NewKeyPool(seeds, coinType)
	if err != nil {
		return
	}

	return
}

func (signer *TxSigner) GetKeys() []keypool.KeyInfo {
	return signer.keyPool.GetKeys()
}

func (signer *TxSigner) SignTx(unsignedTxHex string) (r string, rAllSigned bool, err error) {
	unsignedTx, err := btctx.UnmarshalUnsignedTx(unsignedTxHex)
	if err != nil {
		return
	}

	uncompleted := btctx.NewUncompleted()
	msPrivateKeys := make(map[int][]string)
	privateKeys := make(map[int]string)

	for idx, input := range unsignedTx.Inputs {
		var iPrivateKey string

		iPrivateKey, err = signer.keyPool.GetPrivateKeyByAddress(input.Address)
		if err == nil {
			privateKeys[idx] = iPrivateKey
		}

		i, ok := signer.multiSignAddressInfo[input.Address]
		if !ok {
			uncompleted.SignInputFlag[idx] = privateKeys[idx] != ""

			continue
		}

		msPrivateKeys[idx] = make([]string, 0, 2)

		mgInfo := uncompleted.MultiSignInputInfos[idx]
		mgInfo.MinSignNum = i.MinSignNum

		for _, key := range i.PublicKeys {
			if privateKey, _ := signer.keyPool.GetPrivateKeyByPublicKey(key); privateKey != "" {
				if mgInfo.SignedNum >= mgInfo.MinSignNum {
					continue
				}

				msPrivateKeys[idx] = append(msPrivateKeys[idx], privateKey)
				mgInfo.SignedNum++
			} else {
				mgInfo.PublicKeys = append(uncompleted.MultiSignInputInfos[idx].PublicKeys, key)
			}
		}

		uncompleted.MultiSignInputInfos[idx] = mgInfo
	}

	var tx *wire.MsgTx

	if len(uncompleted.MultiSignInputInfos) == 0 {
		tx, err = btctx.GenSignedTx(unsignedTx, signer.netParams)
	} else {
		tx, err = btctx.GenMultiSignedTx(unsignedTx, uncompleted, signer.netParams)
	}

	if err != nil {
		return
	}

	inputs := make([]bitcoin.Input, 0, len(unsignedTx.Inputs))
	for _, input := range unsignedTx.Inputs {
		inputs = append(inputs, bitcoin.GenInput(input.TxID, input.VOut, "",
			input.RedeemScript, input.Address, input.Amount))
	}

	if len(uncompleted.MultiSignInputInfos) == 0 {
		err = bitcoin.SignBuildTx(tx, inputs, privateKeys, signer.netParams)
	} else {
		err = bitcoin.MultiSignBuildTx(tx, inputs, msPrivateKeys, privateKeys, signer.netParams)
	}

	if err != nil {
		return
	}

	if uncompleted == nil || uncompleted.Completed() {
		r, err = bitcoin.GetTxHex(tx)
		if err != nil {
			return
		}

		rAllSigned = true

		return
	}

	r, err = MarshalMiddleSignMidTx(&MiddleSignMidTx{
		unsignedTx:  unsignedTx,
		uncompleted: uncompleted,
		tx:          tx,
	})
	if err != nil {
		return
	}

	return
}

func (signer *TxSigner) UpdateMiddleSignedTxHex(msTxHex string) (
	r string, rAllSigned bool, err error) {
	msTx, err := UnmarshalMiddleSignMidTx(msTxHex)
	if err != nil {
		return
	}

	if !msTx.Check(signer.netParams) {
		err = commerr.ErrBadFormat

		return
	}

	var updateCount int

	msPrivateKeys := make(map[int][]string)
	privateKeys := make(map[int]string)

	for idx, input := range msTx.unsignedTx.Inputs {
		if msTx.uncompleted.SignInputFlag[idx] {
			continue
		}

		var iPrivateKey string

		iPrivateKey, err = signer.keyPool.GetPrivateKeyByAddress(input.Address)
		if err == nil {
			privateKeys[idx] = iPrivateKey

			updateCount++
		}

		i, ok := signer.multiSignAddressInfo[input.Address]
		if !ok {
			msTx.uncompleted.SignInputFlag[idx] = true

			continue
		}

		msPrivateKeys[idx] = make([]string, 0, 2)

		mgInfo := msTx.uncompleted.MultiSignInputInfos[idx]

		for _, key := range i.PublicKeys {
			if !slices.Contains(mgInfo.PublicKeys, key) {
				continue
			}

			if mgInfo.SignedNum >= mgInfo.MinSignNum {
				break
			}

			if privateKey, _ := signer.keyPool.GetPrivateKeyByPublicKey(key); privateKey != "" {
				msPrivateKeys[idx] = append(msPrivateKeys[idx], privateKey)

				mgInfo.SignedNum++

				updateCount++
			}
		}

		msTx.uncompleted.MultiSignInputInfos[idx] = mgInfo
	}

	if updateCount <= 0 {
		err = commerr.ErrResourceExhausted

		return
	}

	if len(msTx.uncompleted.MultiSignInputInfos) == 0 {
		err = bitcoin.SignBuildTx(msTx.tx, msTx.bitcoinInputs, privateKeys, signer.netParams)
	} else {
		err = bitcoin.MultiSignBuildTx(msTx.tx, msTx.bitcoinInputs, msPrivateKeys, privateKeys, signer.netParams)
	}

	if err != nil {
		return
	}

	if msTx.uncompleted.Completed() {
		r, err = bitcoin.GetTxHex(msTx.tx)
		if err != nil {
			return
		}

		rAllSigned = true

		return
	}

	r, err = MarshalMiddleSignMidTx(msTx)

	return
}

//
//
//

func (signer *TxSigner) GetPrivateKey(address string) (string, error) {
	return signer.keyPool.GetPrivateKeyByAddress(address)
}
