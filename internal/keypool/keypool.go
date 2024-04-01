package keypool

import (
	hdwallet "github.com/dig-coins/hd-wallet"
	"github.com/sgostarter/i/commerr"
)

type KeyInfo struct {
	Idx                 uint32
	WifPrivateKey       string
	PublicKey           string
	Address             string
	AddressSegWit       string
	AddressSegWitNative string
}

type KeyPool struct {
	kM   map[string]string
	kP   map[string]string
	keys []KeyInfo
}

func NewKeyPool(seeds []string, coinType uint32) (kp *KeyPool, err error) {
	kp = &KeyPool{
		kM: make(map[string]string),
		kP: make(map[string]string),
	}

	err = kp.init(seeds, coinType)
	if err != nil {
		return
	}

	return
}

func (kp *KeyPool) init(seeds []string, coinType uint32) (err error) {
	for _, seed := range seeds {
		err = kp.initOneSeed(seed, coinType)
		if err != nil {
			break
		}
	}

	return
}

func (kp *KeyPool) initOneSeed(seed string, coinType uint32) (err error) {
	seedD, err := hdwallet.NewSeed(seed, "", "")
	if err != nil {
		return
	}

	master, err := hdwallet.NewKey(
		hdwallet.Seed(seedD),
	)
	if err != nil {
		return
	}

	for idx := uint32(0); idx < 10; idx++ {
		wallet, _ := master.GetWallet(hdwallet.Purpose(hdwallet.ZeroQuote+49),
			hdwallet.CoinType(coinType), hdwallet.AddressIndex(idx))

		iKey := KeyInfo{
			Idx: idx,
		}

		iKey.WifPrivateKey, err = wallet.GetKey().PrivateWIF(true)
		if err != nil {
			return
		}

		iKey.PublicKey = wallet.GetKey().PublicHex(true)

		iKey.Address, err = wallet.GetAddress()
		if err != nil {
			return
		}

		kp.kM[iKey.Address] = iKey.WifPrivateKey

		iKey.AddressSegWit, err = wallet.GetKey().AddressP2WPKHInP2SH()
		if err != nil {
			return
		}

		kp.kM[iKey.AddressSegWit] = iKey.WifPrivateKey

		iKey.AddressSegWitNative, err = wallet.GetKey().AddressP2WPKH()
		if err != nil {
			return
		}

		kp.kM[iKey.AddressSegWitNative] = iKey.WifPrivateKey

		kp.kP[wallet.GetKey().PublicHex(true)] = iKey.WifPrivateKey

		kp.keys = append(kp.keys, iKey)
	}

	return
}

func (kp *KeyPool) GetKeys() []KeyInfo {
	keys := append([]KeyInfo{}, kp.keys...)

	for idx := 0; idx < len(keys); idx++ {
		keys[idx].WifPrivateKey = ""
	}

	return keys
}

func (kp *KeyPool) GetPrivateKeyByAddress(address string) (string, error) {
	key, ok := kp.kM[address]
	if !ok {
		return "", commerr.ErrNotFound
	}

	return key, nil
}

func (kp *KeyPool) GetPrivateKeyByPublicKey(publicKey string) (string, error) {
	key, ok := kp.kP[publicKey]
	if !ok {
		return "", commerr.ErrNotFound
	}

	return key, nil
}
