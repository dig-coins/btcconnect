// nolint
package btcserver

import (
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/dig-coins/btcconnect/internal/share"
	"github.com/dig-coins/btcconnect/internal/utl"
	"github.com/okx/go-wallet-sdk/coins/bitcoin"
	"github.com/stretchr/testify/assert"
)

var _uTConfig *utl.UTConfig

func TestMain(m *testing.M) {
	_uTConfig = utl.Setup()

	m.Run()
}

func utNewBTCServer(t *testing.T) *BTCServer {
	return NewBTCServer(&Config{
		CoinType:    share.CoinTypeBTCTestnet,
		Listens:     ":9000",
		RPCHost:     _uTConfig.GetBTCRpcHost(t),
		RPCPort:     _uTConfig.GetBTCRpcPort(t),
		RPCUser:     _uTConfig.GetBTCRpcUser(t),
		RPCPassword: _uTConfig.GetBTCRpcPassword(t),
		RPCUseSSL:   _uTConfig.GetBTCRpcUseTLS(t),
		MultiSignAddressInfos: map[string]*share.MultiSignAddressInfo{
			"2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i": {
				PublicKeys: []string{
					"02f410e07213396b8d6289ca6f1c217380a2787db5a7487b0978bd792cbd32343e",
					"028e7bbb364b64687db98ad39674e70c194aa01ef3da3148791536f25ab8861c02",
					"037cb184baf184ec6d3f24c927c6243dbb122ae7f9353425f84f3dfcd95fdcd0f9",
					"037f9e1716a16efb9cf977a628942a5ee193aa29d31d99a6ddbb80d1d7dd69b5da",
				},
				MinSignNum: 3,
			},
		},
	}, nil, nil)
}

// nolint
func TestBroadcastTx(t *testing.T) {
	s := utNewBTCServer(t)

	txID, err := s.SendRawTransaction("0100000001699266e5d7e176129f8b333385769d82cb1a05568b3dd79b297f179606b964f000000000fd680100483045022100b82732cd52474520dfa352057866d6ffffe330be4a52a22bb61995a7d10645b502200edf59b715fb1768d3a5bb7fd560a225165dab7eae2fcb9277e142229ead46d801483045022100d1337fd16c222fc3a06db9a94f51e2a9ef142e2b12ebaf3e0826c46e428af92c0220470137999b0c0ada536bab6115b76831bf0d9efa888baefe33a5d3235c5194be0147304402206fad35910b1b572481e4ea59a9f8cfbe263df7a384b8ff3d44985526704da76a02205828f5547b251205626bf1ae27688625d2137952153f29c27056e576cf68e444014c8b532102f410e07213396b8d6289ca6f1c217380a2787db5a7487b0978bd792cbd32343e21028e7bbb364b64687db98ad39674e70c194aa01ef3da3148791536f25ab8861c0221037cb184baf184ec6d3f24c927c6243dbb122ae7f9353425f84f3dfcd95fdcd0f921037f9e1716a16efb9cf977a628942a5ee193aa29d31d99a6ddbb80d1d7dd69b5da54aeffffffff02e8030000000000001976a914b34f7caffd40224fa0e717020870aed0c04b3a9b88acf10500000000000017a914828554491599c7989112bdc632b14a64c1b5ed4f8700000000")
	assert.Nil(t, err)
	t.Log(txID)
}

func TestTrans2(t *testing.T) {
	s := utNewBTCServer(t)

	wpi, err := s.GenUnsignedTx4TransToOne(_uTConfig.GetWallet(t), "mws4UFRP8XE8JhweXhgyMGkVPZfCMSFgmx",
		"tb1q9esscxqfepv26juy6eay4ejf4vsc3dywk9xqqv", 3000, 2, "", true)
	assert.Nil(t, err)
	t.Log(wpi)
}

func TestTrans3(t *testing.T) {
	s := utNewBTCServer(t)

	wpi, err := s.GenUnsignedTx4Gather(_uTConfig.GetWallet(t), []string{
		"34KAAvz5Nad7G5wk56PKjPpVLfMjjkaeQb"},
		40, "36pjdq8Nb7XL2AuerdSeBxr6yazKomLgdV", true) // 4,21
	assert.Nil(t, err)
	t.Log(wpi)
}

func TestTrans4(t *testing.T) {
	s := utNewBTCServer(t)

	wpi, err := s.GenUnsignedTx4TransTo(_uTConfig.GetWallet(t), []string{
		"tb1qkd8hetlagq3ylg88zupqsu9w6rqykw5mvthhr4", "2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i",
	}, []TransOutput{
		{
			Address: "2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i",
			Amount:  70000,
		},
	}, 2, "mws4UFRP8XE8JhweXhgyMGkVPZfCMSFgmx", true)
	assert.Nil(t, err)
	t.Log(wpi)
}

// nolint
func Test_selectUnspentInputs(t *testing.T) {
	s := utNewBTCServer(t)

	unspentList := []Unspent{{
		TxID:          "63a7a7344e530552cd3937cf51f4323fb96a206e88208a0cca60205036738307",
		VOut:          0,
		Safe:          true,
		Confirmations: 6,
		Spendable:     true,
		Amount:        0.00001111,
		Address:       "bc1q37kc3s6fnwwvn973lff0sy22kznptyc2skx0rw",
	}, {
		TxID:          "63a7a7344e530552cd3937cf51f4323fb96a206e88208a0cca60205036738307",
		VOut:          1,
		Safe:          true,
		Confirmations: 6,
		Spendable:     true,
		Amount:        0.00000111,
		Address:       "bc1q37kc3s6fnwwvn973lff0sy22kznptyc2skx0rw",
	}, {
		TxID:          "b9f72a4d3d6f25f1c97a34f9452fa68710341f5e8cbedd7cc99772955b5c9577",
		VOut:          0,
		Safe:          true,
		Confirmations: 6,
		Spendable:     true,
		Amount:        0.00008111,
		Address:       "2NC7ACW9obBtsgAgmciLrqJJ36iVcG3Gkgq",
	}}

	inputs, outputs, err := s.selectUnspentInputs(unspentList, nil, 100, nil,
		"mojMsSaRFDtF14NnhdNYpRnL1e5CbAzLUU", true)
	assert.Nil(t, err)
	assert.EqualValues(t, 3, len(inputs))
	assert.EqualValues(t, 1, len(outputs))
	assert.EqualValues(t, "mojMsSaRFDtF14NnhdNYpRnL1e5CbAzLUU", outputs[0].Address)
	assert.EqualValues(t, 9304, outputs[0].Amount)

	inputs, outputs, err = s.selectUnspentInputs(unspentList, nil, 100, []TransOutput{{
		Address: "2Mu5idgPdH9cSLb2kBUn9ZEB5z6rUdmkLDX",
		Amount:  100,
	}}, "bc1q37kc3s6fnwwvn973lff0sy22kznptyc2skx0rw", true)
	assert.EqualValues(t, 1, len(inputs))
	assert.EqualValues(t, "bc1q37kc3s6fnwwvn973lff0sy22kznptyc2skx0rw", inputs[0].Address)

	assert.EqualValues(t, 2, len(outputs))
	assert.EqualValues(t, "2Mu5idgPdH9cSLb2kBUn9ZEB5z6rUdmkLDX", outputs[0].Address)
	assert.EqualValues(t, 100, outputs[0].Amount)

	assert.EqualValues(t, "bc1q37kc3s6fnwwvn973lff0sy22kznptyc2skx0rw", outputs[1].Address)
	assert.EqualValues(t, 995, outputs[1].Amount)

	inputs, outputs, err = s.selectUnspentInputs(unspentList, nil, 100, []TransOutput{{
		Address: "2Mu5idgPdH9cSLb2kBUn9ZEB5z6rUdmkLDX",
		Amount:  100,
	}}, "", true)
	assert.EqualValues(t, 1, len(inputs))
	assert.EqualValues(t, "bc1q37kc3s6fnwwvn973lff0sy22kznptyc2skx0rw", inputs[0].Address)

	assert.EqualValues(t, 2, len(outputs))
	assert.EqualValues(t, "2Mu5idgPdH9cSLb2kBUn9ZEB5z6rUdmkLDX", outputs[0].Address)
	assert.EqualValues(t, 100, outputs[0].Amount)

	assert.EqualValues(t, "bc1q37kc3s6fnwwvn973lff0sy22kznptyc2skx0rw", outputs[1].Address)
	assert.EqualValues(t, 995, outputs[1].Amount)

	inputs, outputs, err = s.selectUnspentInputs(unspentList, []string{"2NC7ACW9obBtsgAgmciLrqJJ36iVcG3Gkgq"}, 100, []TransOutput{{
		Address: "2Mu5idgPdH9cSLb2kBUn9ZEB5z6rUdmkLDX",
		Amount:  100,
	}}, "bc1q37kc3s6fnwwvn973lff0sy22kznptyc2skx0rw", true)
	assert.EqualValues(t, 1, len(inputs))
	assert.EqualValues(t, "2NC7ACW9obBtsgAgmciLrqJJ36iVcG3Gkgq", inputs[0].Address)

	assert.EqualValues(t, 2, len(outputs))
	assert.EqualValues(t, "2Mu5idgPdH9cSLb2kBUn9ZEB5z6rUdmkLDX", outputs[0].Address)
	assert.EqualValues(t, 100, outputs[0].Amount)

	assert.EqualValues(t, "bc1q37kc3s6fnwwvn973lff0sy22kznptyc2skx0rw", outputs[1].Address)
	assert.EqualValues(t, 7993, outputs[1].Amount)
}

func TestBTCServer_GenMultiSignatureAddress(t *testing.T) {
	s := utNewBTCServer(t)

	address, err := s.GenMultiSignatureAddress([]string{
		"02f410e07213396b8d6289ca6f1c217380a2787db5a7487b0978bd792cbd32343e",
		"028e7bbb364b64687db98ad39674e70c194aa01ef3da3148791536f25ab8861c02",
		"037cb184baf184ec6d3f24c927c6243dbb122ae7f9353425f84f3dfcd95fdcd0f9",
		"037f9e1716a16efb9cf977a628942a5ee193aa29d31d99a6ddbb80d1d7dd69b5da",
	}, 3, &chaincfg.TestNet3Params)
	assert.Nil(t, err)
	assert.EqualValues(t, "2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i", address)
}

func TestTrans2MultiSignAddress(t *testing.T) {
	s := utNewBTCServer(t)

	wpi, err := s.GenUnsignedTx4TransToOne(_uTConfig.GetWallet(t), "tb1q88h06gxq5cgpy566832d6qzcyaey7xhddl2het",
		"2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i", 20823, 2, "", true)
	assert.Nil(t, err)
	t.Log(wpi)
}

// nolint
func TestTransFromMultiSigAndWitnessAddress(t *testing.T) {
	/*
			2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i 20823
			tb1q88h06gxq5cgpy566832d6qzcyaey7xhddl2het 19856

		->
			mws4UFRP8XE8JhweXhgyMGkVPZfCMSFgmx 30000
			tb1q88h06gxq5cgpy566832d6qzcyaey7xhddl2het *
	*/

	s := utNewBTCServer(t)

	wpi, err := s.GenUnsignedTx4TransTo(_uTConfig.GetWallet(t), []string{
		"tb1qkd8hetlagq3ylg88zupqsu9w6rqykw5mvthhr4",
		"2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i",
	}, []TransOutput{
		{
			Address: "mws4UFRP8XE8JhweXhgyMGkVPZfCMSFgmx",
			Amount:  70100,
		},
	}, 2, "2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i", true)
	assert.Nil(t, err)
	t.Log(wpi)
}

func TestTransFromMultiSigAndWitnessAddress2(t *testing.T) {
	s := utNewBTCServer(t)

	wpi, err := s.GenUnsignedTx4TransTo(_uTConfig.GetWallet(t), []string{
		"2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i",
	}, []TransOutput{
		{
			Address: "mws4UFRP8XE8JhweXhgyMGkVPZfCMSFgmx",
			Amount:  1000,
		},
	}, 2, "2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i", true)
	assert.Nil(t, err)
	t.Log(wpi)
}

func TestStatisticsFee(t *testing.T) {
	s := utNewBTCServer(t)

	fnPrintFee := func(label string, allAmount int64, outputs []TransOutput) {
		var outAmount int64

		for _, output := range outputs {
			outAmount += output.Amount
		}

		t.Logf("%s, fee is %d\n", label, allAmount-outAmount)
	}

	redeemScript, err := bitcoin.GetRedeemScript([]string{
		"02f410e07213396b8d6289ca6f1c217380a2787db5a7487b0978bd792cbd32343e",
		"028e7bbb364b64687db98ad39674e70c194aa01ef3da3148791536f25ab8861c02",
		"037cb184baf184ec6d3f24c927c6243dbb122ae7f9353425f84f3dfcd95fdcd0f9",
		"037f9e1716a16efb9cf977a628942a5ee193aa29d31d99a6ddbb80d1d7dd69b5da",
	}, 3, &chaincfg.TestNet3Params)
	assert.Nil(t, err)
	redeemScriptHex := hex.EncodeToString(redeemScript)
	t.Log(redeemScriptHex)

	//
	//
	//

	outputs, err := s.calcChange4Trans([]TransInput{
		{
			TxID:    "63a7a7344e530552cd3937cf51f4323fb96a206e88208a0cca60205036738307",
			VOut:    0,
			Address: "mojMsSaRFDtF14NnhdNYpRnL1e5CbAzLUU",
			Amount:  100000,
		},
	}, []TransOutput{}, "mojMsSaRFDtF14NnhdNYpRnL1e5CbAzLUU", 1000)
	assert.Nil(t, err)
	fnPrintFee("L", 100000, outputs)

	outputs, err = s.calcChange4Trans([]TransInput{
		{
			TxID:    "63a7a7344e530552cd3937cf51f4323fb96a206e88208a0cca60205036738307",
			VOut:    0,
			Address: "mojMsSaRFDtF14NnhdNYpRnL1e5CbAzLUU",
			Amount:  100000,
		},
	}, []TransOutput{}, "tb1qkd8hetlagq3ylg88zupqsu9w6rqykw5mvthhr4", 1000)
	assert.Nil(t, err)
	fnPrintFee("L", 100000, outputs)

	outputs, err = s.calcChange4Trans([]TransInput{
		{
			TxID:    "63a7a7344e530552cd3937cf51f4323fb96a206e88208a0cca60205036738307",
			VOut:    0,
			Address: "mojMsSaRFDtF14NnhdNYpRnL1e5CbAzLUU",
			Amount:  100000,
		},
	}, []TransOutput{}, "2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i", 1000)
	assert.Nil(t, err)
	fnPrintFee("L", 100000, outputs)

	//
	//
	//

	outputs, err = s.calcChange4Trans([]TransInput{
		{
			TxID:    "63a7a7344e530552cd3937cf51f4323fb96a206e88208a0cca60205036738307",
			VOut:    0,
			Address: "2NC7ACW9obBtsgAgmciLrqJJ36iVcG3Gkgq",
			Amount:  100000,
		},
	}, []TransOutput{}, "mojMsSaRFDtF14NnhdNYpRnL1e5CbAzLUU", 1000)
	assert.Nil(t, err)
	fnPrintFee("WN", 100000, outputs)

	outputs, err = s.calcChange4Trans([]TransInput{
		{
			TxID:    "63a7a7344e530552cd3937cf51f4323fb96a206e88208a0cca60205036738307",
			VOut:    0,
			Address: "2NC7ACW9obBtsgAgmciLrqJJ36iVcG3Gkgq",
			Amount:  100000,
		},
	}, []TransOutput{}, "tb1qkd8hetlagq3ylg88zupqsu9w6rqykw5mvthhr4", 1000)
	assert.Nil(t, err)
	fnPrintFee("WN", 100000, outputs)

	outputs, err = s.calcChange4Trans([]TransInput{
		{
			TxID:    "63a7a7344e530552cd3937cf51f4323fb96a206e88208a0cca60205036738307",
			VOut:    0,
			Address: "2NC7ACW9obBtsgAgmciLrqJJ36iVcG3Gkgq",
			Amount:  100000,
		},
	}, []TransOutput{}, "2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i", 1000)
	assert.Nil(t, err)
	fnPrintFee("WN", 100000, outputs)

	//
	//
	//

	outputs, err = s.calcChange4Trans([]TransInput{
		{
			TxID:    "63a7a7344e530552cd3937cf51f4323fb96a206e88208a0cca60205036738307",
			VOut:    0,
			Address: "tb1qkd8hetlagq3ylg88zupqsu9w6rqykw5mvthhr4",
			Amount:  100000,
		},
	}, []TransOutput{}, "mojMsSaRFDtF14NnhdNYpRnL1e5CbAzLUU", 1000)
	assert.Nil(t, err)
	fnPrintFee("W", 100000, outputs)

	outputs, err = s.calcChange4Trans([]TransInput{
		{
			TxID:    "63a7a7344e530552cd3937cf51f4323fb96a206e88208a0cca60205036738307",
			VOut:    0,
			Address: "tb1qkd8hetlagq3ylg88zupqsu9w6rqykw5mvthhr4",
			Amount:  100000,
		},
	}, []TransOutput{}, "tb1qkd8hetlagq3ylg88zupqsu9w6rqykw5mvthhr4", 1000)
	assert.Nil(t, err)
	fnPrintFee("W", 100000, outputs)

	outputs, err = s.calcChange4Trans([]TransInput{
		{
			TxID:    "63a7a7344e530552cd3937cf51f4323fb96a206e88208a0cca60205036738307",
			VOut:    0,
			Address: "tb1qkd8hetlagq3ylg88zupqsu9w6rqykw5mvthhr4",
			Amount:  100000,
		},
	}, []TransOutput{}, "2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i", 1000)
	assert.Nil(t, err)
	fnPrintFee("W", 100000, outputs)

	//
	//
	//
	outputs, err = s.calcChange4Trans([]TransInput{
		{
			TxID:         "63a7a7344e530552cd3937cf51f4323fb96a206e88208a0cca60205036738307",
			VOut:         0,
			Address:      "2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i",
			Amount:       100000,
			RedeemScript: redeemScriptHex,
		},
	}, []TransOutput{}, "mojMsSaRFDtF14NnhdNYpRnL1e5CbAzLUU", 1000)
	assert.Nil(t, err)
	fnPrintFee("MS", 100000, outputs)

	outputs, err = s.calcChange4Trans([]TransInput{
		{
			TxID:         "63a7a7344e530552cd3937cf51f4323fb96a206e88208a0cca60205036738307",
			VOut:         0,
			Address:      "2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i",
			Amount:       100000,
			RedeemScript: redeemScriptHex,
		},
	}, []TransOutput{}, "tb1qkd8hetlagq3ylg88zupqsu9w6rqykw5mvthhr4", 1000)
	assert.Nil(t, err)
	fnPrintFee("MS", 100000, outputs)

	outputs, err = s.calcChange4Trans([]TransInput{
		{
			TxID:         "63a7a7344e530552cd3937cf51f4323fb96a206e88208a0cca60205036738307",
			VOut:         0,
			Address:      "2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i",
			Amount:       100000,
			RedeemScript: redeemScriptHex,
		},
	}, []TransOutput{}, "2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i", 1000)
	assert.Nil(t, err)
	fnPrintFee("MS", 100000, outputs)

}
