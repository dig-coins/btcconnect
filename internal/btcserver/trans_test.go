package btcserver

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/dig-coins/btcconnect/internal/redistorage"
	"github.com/dig-coins/btcconnect/internal/share"
	"github.com/dig-coins/btcconnect/internal/utl"
	"github.com/stretchr/testify/assert"
)

var _uTConfig *utl.UTConfig

func TestMain(m *testing.M) {
	_uTConfig = utl.Setup()

	m.Run()
}

func utNewBTCServer(t *testing.T) *BTCServer {
	stg, err := redistorage.NewRedisStorage()
	assert.Nil(t, err)

	return NewBTCServer(&Config{
		RPCHost:     _uTConfig.GetBTCRpcHost(t),
		RPCPort:     _uTConfig.GetBTCRpcPort(t),
		RPCUser:     _uTConfig.GetBTCRpcUser(t),
		RPCPassword: _uTConfig.GetBTCRpcPassword(t),
		RPCUseSSL:   _uTConfig.GetBTCRpcUseTLS(t),
	}, stg, &chaincfg.TestNet3Params, map[string]*share.MultiSignAddressInfo{
		"2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i": {
			PublicKeys: []string{
				"02f410e07213396b8d6289ca6f1c217380a2787db5a7487b0978bd792cbd32343e",
				"028e7bbb364b64687db98ad39674e70c194aa01ef3da3148791536f25ab8861c02",
				"037cb184baf184ec6d3f24c927c6243dbb122ae7f9353425f84f3dfcd95fdcd0f9",
				"037f9e1716a16efb9cf977a628942a5ee193aa29d31d99a6ddbb80d1d7dd69b5da",
			},
			MinSignNum: 3,
		},
	}, nil)
}

func TestBroadcastTx(t *testing.T) {
	s := utNewBTCServer(t)

	txID, err := s.SendRawTransaction("010000000001028aebc0cb584b62531438f663e4267901c26d6766078549369c12cbfcb7b0da7500000000fd680100483045022100a7181a0bfacc26d828555c5374ef3f561ee028263f8a58543cdaa2984116711002207f9a8bfb4dc12b48f4454f5f414d8430e0a0aa6e97dfd9cfc538ed388ff7197901483045022100a834a1a57a5864e49170eb00a009ac645f3de524d7d610df13412c0280c7ddc4022059255277dbbc75437b1a4e92e6fe42f18cc8f899e4e0bbc61c30266608c1fc2601473044022044280a73520c12a783c69c25997393cf6709001cac8ad4fd98f5da7e435ef64002205dd8b0196e1c4ba61c0cfa53e9dbc4c55c5ad9dfe3904041738bfe574e376f9f014c8b532102f410e07213396b8d6289ca6f1c217380a2787db5a7487b0978bd792cbd32343e21028e7bbb364b64687db98ad39674e70c194aa01ef3da3148791536f25ab8861c0221037cb184baf184ec6d3f24c927c6243dbb122ae7f9353425f84f3dfcd95fdcd0f921037f9e1716a16efb9cf977a628942a5ee193aa29d31d99a6ddbb80d1d7dd69b5da54aeffffffff4c6fbea4631176d6205c936bee2345564f2bd19c1eaa7f76db890da68eff97a90000000000ffffffff02d4110100000000001976a914b34f7caffd40224fa0e717020870aed0c04b3a9b88ac5c0100000000000017a914828554491599c7989112bdc632b14a64c1b5ed4f87000247304402202b56044f3c497995d9dfb6b2b5bedb4d32229f4d3426c4c2dadea978400632db022016345b7ab15caa277e795c0860b46e3c7ac6f8fcee5365a7d78d7fac96e4cfe80121037cb184baf184ec6d3f24c927c6243dbb122ae7f9353425f84f3dfcd95fdcd0f900000000")
	assert.Nil(t, err)
	t.Log(txID)
}

func TestTrans2(t *testing.T) {
	s := utNewBTCServer(t)

	wpi, err := s.GenUnsignedTx4TransToOne("hoho2", "mws4UFRP8XE8JhweXhgyMGkVPZfCMSFgmx",
		"tb1q88h06gxq5cgpy566832d6qzcyaey7xhddl2het", 1000, 2, "")
	assert.Nil(t, err)
	t.Log(wpi)
}

func TestTrans3(t *testing.T) {
	s := utNewBTCServer(t)

	wpi, err := s.GenUnsignedTx4Gather("hoho2", []string{
		"tb1qkd8hetlagq3ylg88zupqsu9w6rqykw5mvthhr4"},
		2, "mws4UFRP8XE8JhweXhgyMGkVPZfCMSFgmx")
	assert.Nil(t, err)
	t.Log(wpi)
}

func TestTrans4(t *testing.T) {
	s := utNewBTCServer(t)

	wpi, err := s.GenUnsignedTx4TransTo("hoho2", []string{
		"tb1qkd8hetlagq3ylg88zupqsu9w6rqykw5mvthhr4", "2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i",
	}, []TransOutput{
		{
			Address: "2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i",
			Amount:  70000,
		},
	}, 2, "mws4UFRP8XE8JhweXhgyMGkVPZfCMSFgmx")
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
		"mojMsSaRFDtF14NnhdNYpRnL1e5CbAzLUU")
	assert.Nil(t, err)
	assert.EqualValues(t, 3, len(inputs))
	assert.EqualValues(t, 1, len(outputs))
	assert.EqualValues(t, "mojMsSaRFDtF14NnhdNYpRnL1e5CbAzLUU", outputs[0].Address)
	assert.EqualValues(t, 9304, outputs[0].Amount)

	inputs, outputs, err = s.selectUnspentInputs(unspentList, nil, 100, []TransOutput{{
		Address: "2Mu5idgPdH9cSLb2kBUn9ZEB5z6rUdmkLDX",
		Amount:  100,
	}}, "bc1q37kc3s6fnwwvn973lff0sy22kznptyc2skx0rw")
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
	}}, "")
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
	}}, "bc1q37kc3s6fnwwvn973lff0sy22kznptyc2skx0rw")
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

	wpi, err := s.GenUnsignedTx4TransToOne("hoho2", "tb1q88h06gxq5cgpy566832d6qzcyaey7xhddl2het",
		"2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i", 20823, 2, "")
	assert.Nil(t, err)
	t.Log(wpi)
}

func TestTransFromMultiSigAndWitnessAddress(t *testing.T) {
	/*
			2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i 20823
			tb1q88h06gxq5cgpy566832d6qzcyaey7xhddl2het 19856

		->
			mws4UFRP8XE8JhweXhgyMGkVPZfCMSFgmx 30000
			tb1q88h06gxq5cgpy566832d6qzcyaey7xhddl2het *
	*/

	s := utNewBTCServer(t)

	wpi, err := s.GenUnsignedTx4TransTo("hoho2", []string{
		"2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i",
		"tb1q88h06gxq5cgpy566832d6qzcyaey7xhddl2het",
	}, []TransOutput{
		{
			Address: "mws4UFRP8XE8JhweXhgyMGkVPZfCMSFgmx",
			Amount:  70100,
		},
	}, 2, "2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i")
	assert.Nil(t, err)
	t.Log(wpi)
}
