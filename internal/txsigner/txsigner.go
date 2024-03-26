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

func (signer *TxSigner) GetSignedTxHex(unsignedTxHex string) (r string, rIsMultiSignMidTx bool, err error) {
	unsignedTx, err := btctx.UnmarshalUnsignedTx(unsignedTxHex)
	if err != nil {
		return
	}

	var uncompleted *btctx.Uncompleted

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
			continue
		}

		if uncompleted == nil {
			uncompleted = btctx.NewUncompleted()
		}

		msPrivateKeys[idx] = make([]string, 0, 2)

		mgInfo := uncompleted.MultiSignInfos[idx]
		mgInfo.MinSignNum = i.MinSignNum

		for _, key := range i.PublicKeys {
			if privateKey, _ := signer.keyPool.GetPrivateKeyByPublicKey(key); privateKey != "" {
				if mgInfo.SignedNum >= mgInfo.MinSignNum {
					continue
				}
				msPrivateKeys[idx] = append(msPrivateKeys[idx], privateKey)
				mgInfo.SignedNum++
			} else {
				mgInfo.PublicKeys = append(uncompleted.MultiSignInfos[idx].PublicKeys, key)
			}
		}

		uncompleted.MultiSignInfos[idx] = mgInfo
	}

	var tx *wire.MsgTx

	if uncompleted == nil {
		tx, err = btctx.GenSignedTx(unsignedTx, signer.netParams, signer)
	} else {
		tx, err = btctx.GenMultiSignedTx(unsignedTx, uncompleted, signer.netParams)
		if err != nil {
			return
		}

		inputs := make([]bitcoin.Input, 0, len(unsignedTx.Inputs))
		for idx, input := range unsignedTx.Inputs {
			inputs = append(inputs, bitcoin.GenInput(input.TxID, input.VOut, privateKeys[idx],
				input.RedeemScript, input.Address, input.Amount))
		}

		err = btctx.UpdateMultiSignedTx(tx, inputs, msPrivateKeys, privateKeys, signer.netParams)
	}

	if err != nil {
		return
	}

	if uncompleted == nil || uncompleted.Completed() {
		r, err = bitcoin.GetTxHex(tx)

		return
	}

	r, err = MarshalMultiSignMidTx(&MultiSignMidTx{
		unsignedTx:  unsignedTx,
		uncompleted: uncompleted,
		tx:          tx,
	})
	if err != nil {
		return
	}

	rIsMultiSignMidTx = true

	return
}

func (signer *TxSigner) UpdateMultiSignedTxHex(msTxHex string, uncompleted *btctx.Uncompleted) (
	r string, rIsMultiSignMidTx bool, err error) {
	msTx, err := UnmarshalMultiSignMidTx(msTxHex)
	if err != nil {
		return
	}

	var updateCount int

	msPrivateKeys := make(map[int][]string)
	privateKeys := make(map[int]string)

	for idx, input := range msTx.unsignedTx.Inputs {
		var iPrivateKey string

		iPrivateKey, err = signer.keyPool.GetPrivateKeyByAddress(input.Address)
		if err == nil {
			privateKeys[idx] = iPrivateKey
		}

		i, ok := signer.multiSignAddressInfo[input.Address]
		if !ok {
			continue
		}

		msPrivateKeys[idx] = make([]string, 0, 2)

		mgInfo := uncompleted.MultiSignInfos[idx]

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

		uncompleted.MultiSignInfos[idx] = mgInfo
	}

	if updateCount <= 0 {
		err = commerr.ErrResourceExhausted

		return
	}

	inputs := make([]bitcoin.Input, 0, len(msTx.unsignedTx.Inputs))
	for idx, input := range msTx.unsignedTx.Inputs {
		inputs = append(inputs, bitcoin.GenInput(input.TxID, input.VOut, privateKeys[idx],
			input.RedeemScript, input.Address, input.Amount))
	}

	err = btctx.UpdateMultiSignedTx(msTx.tx, inputs, msPrivateKeys, privateKeys, signer.netParams)
	if err != nil {
		return
	}

	if msTx.uncompleted.Completed() {
		r, err = bitcoin.GetTxHex(msTx.tx)
		if err != nil {
			return
		}

		return
	}

	r, err = MarshalMultiSignMidTx(msTx)
	if err != nil {
		return
	}

	rIsMultiSignMidTx = true

	return
}

//
//
//

func (signer *TxSigner) GetPrivateKey(address string) (string, error) {
	return signer.keyPool.GetPrivateKeyByAddress(address)
}
