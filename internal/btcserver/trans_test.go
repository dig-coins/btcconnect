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

// nolint
func TestBroadcastTx(t *testing.T) {
	s := utNewBTCServer(t)

	txID, err := s.SendRawTransaction("01000000000103dc2dd538755d0f5c4326cab47e02ea184f3628ffde112235f34be7e91b3e382700000000fd670100483045022100fddc32bb043c9753a8d72e572e681fa189055da9617743049da31154a5537087022076f8b2e76256a68f451242a3f680acc40dfbf50bf99cbfdb0f12148c8479685b0147304402200d8297611c58b6a322272224d504a17eb50c3e24cc52fb53264b199d681259270220566b3ab6677a62a09fabf3c0abf5378a203a794f513ab74273c51ee69f33eb5f014730440220162a745e8345cfa8ac095ae006f8c0bb4ef9631780f80bf561a1504474c645ad022077881eec386e9f909217d29cc59bb4f3d0678272d24856e8981a289910f27b93014c8b532102f410e07213396b8d6289ca6f1c217380a2787db5a7487b0978bd792cbd32343e21028e7bbb364b64687db98ad39674e70c194aa01ef3da3148791536f25ab8861c0221037cb184baf184ec6d3f24c927c6243dbb122ae7f9353425f84f3dfcd95fdcd0f921037f9e1716a16efb9cf977a628942a5ee193aa29d31d99a6ddbb80d1d7dd69b5da54aeffffffff625bd1211b518d6db681280f0a87be58705a6c514858802a9df8c97708a28a1201000000fd670100473044022072ff935343767dd1203ac91ec6d6c6a2d375c540bf0e3e39ef3ee6f95f2c97d802202b2611d2729160ceb62cac94e65dc9c81c26670dd593eeb441de54c1b97e4e7d01473044022043468bd2613f7059c7419329b5741ce88dd4d1a6475e61ef62f522b49627576c022005414d1b8f9e343f7ec35271722800c09e4d2a04513cef22c79ce491fc2c9e3101483045022100d9fdfd37bb95447aa09e0a42f46cf095b8e526ac63d5d6f5ed6709aedf3529fe022027c5d6c8a04a827b6197c0ab8fbb882df8bf86ce0fd35e133133f80cfc6b3360014c8b532102f410e07213396b8d6289ca6f1c217380a2787db5a7487b0978bd792cbd32343e21028e7bbb364b64687db98ad39674e70c194aa01ef3da3148791536f25ab8861c0221037cb184baf184ec6d3f24c927c6243dbb122ae7f9353425f84f3dfcd95fdcd0f921037f9e1716a16efb9cf977a628942a5ee193aa29d31d99a6ddbb80d1d7dd69b5da54aefffffffff8295e64ab521e1a242861796a710d1c0886105441dec1ad2e79605b1c7c8aa60000000000ffffffff0173340000000000001976a914b34f7caffd40224fa0e717020870aed0c04b3a9b88ac000002483045022100842d119544616925c44d83788b523137b8dff41c86c8e12ec742e3cba5bf14d80220124b0499a55080636f992771a181a8eca0a74273598247f5dfb9248a7ba59e51012102f410e07213396b8d6289ca6f1c217380a2787db5a7487b0978bd792cbd32343e00000000")
	assert.Nil(t, err)
	t.Log(txID)
}

func TestTrans2(t *testing.T) {
	s := utNewBTCServer(t)

	wpi, err := s.GenUnsignedTx4TransToOne("hoho2", "mws4UFRP8XE8JhweXhgyMGkVPZfCMSFgmx",
		"2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i", 10000, 2, "")
	assert.Nil(t, err)
	t.Log(wpi)
}

func TestTrans3(t *testing.T) {
	s := utNewBTCServer(t)

	wpi, err := s.GenUnsignedTx4Gather("hoho2", []string{
		"tb1qkd8hetlagq3ylg88zupqsu9w6rqykw5mvthhr4", "2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i"},
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

	wpi, err := s.GenUnsignedTx4TransTo("hoho2", []string{
		"tb1qkd8hetlagq3ylg88zupqsu9w6rqykw5mvthhr4",
		"2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i",
	}, []TransOutput{
		{
			Address: "mws4UFRP8XE8JhweXhgyMGkVPZfCMSFgmx",
			Amount:  70100,
		},
	}, 2, "2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i")
	assert.Nil(t, err)
	t.Log(wpi)
}
