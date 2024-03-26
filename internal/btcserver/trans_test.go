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

	txID, err := s.SendRawTransaction("01000000000104e16cfc542c68d9d139a44c77f291f05e78a57ecef3d38a3851235c39df8bc54801000000fd670100483045022100ab0d6869e4809bd6a8e97b40cde79e4c66318a92cf1cb40499ec2958c1aa115602207cc4de577216cd23a8e43483750ab4d1f4e849fc160e85610817c20d391638a2014730440220317f7360f4132ff3a686443da49bbdf3409669eccf4bd60a35a87b49a0cb849a022041d4359f8054c012500d2dcdf62c0e75b1e594545555bd23a3e10f429c72b66d01473044022032506885b7adfaaa90d2bb1691b46a628ebb989b0e654339351b27748f534d4802207e203dd295a59641f03c7569ecc9c9cf996eaabf02502cfd4ca297e13bd3ef14014c8b532102f410e07213396b8d6289ca6f1c217380a2787db5a7487b0978bd792cbd32343e21028e7bbb364b64687db98ad39674e70c194aa01ef3da3148791536f25ab8861c0221037cb184baf184ec6d3f24c927c6243dbb122ae7f9353425f84f3dfcd95fdcd0f921037f9e1716a16efb9cf977a628942a5ee193aa29d31d99a6ddbb80d1d7dd69b5da54aeffffffff1db15c8edea8beeec5c17e1a2b2522e0077ac3c4d3782bb3a14e6827ed53286d01000000fd6701004830450221009258d0c4c6046e77b8e286c464315a80b5e684e597c28add38d86ccc7092ea730220479dc7a07b88c3f308498c2e6c6bf9d24be1e071562a2afd4fb5ef9d6ed1471c0147304402202944b36268bc8a06cd9377c0003eec01469231e77c47aec7d6db3c7dcf837ad202201749462e6de9c42dc24fd858eb5b1e5e395f9214480851e387c3709d6c762a7901473044022043b9c806b7161a319862c33503b7dfd09775c78b29257188a7d78c5e5a7a50530220586ecbae9971a007b9a2ed1dc3537682f5c9cfd58619ed00dedd1bbcfa978fc1014c8b532102f410e07213396b8d6289ca6f1c217380a2787db5a7487b0978bd792cbd32343e21028e7bbb364b64687db98ad39674e70c194aa01ef3da3148791536f25ab8861c0221037cb184baf184ec6d3f24c927c6243dbb122ae7f9353425f84f3dfcd95fdcd0f921037f9e1716a16efb9cf977a628942a5ee193aa29d31d99a6ddbb80d1d7dd69b5da54aeffffffff47aa78ca957f678f8bc5d5aa20d85c2eb552f73466bda7df9bc04336d13032a20100000000ffffffff7acca0029497b801cc6fe7ca3a804ca1f9fb9a574a8bb9ae82eff7aec2ed823a0000000000ffffffff0260ea0000000000001976a914b34f7caffd40224fa0e717020870aed0c04b3a9b88ac811500000000000017a914828554491599c7989112bdc632b14a64c1b5ed4f8700000247304402206bee4187a68df81174d7d566ed0cd28d9e37fcdd98b90992f2718021e4a0fd030220192ad8ba7db9bad03b21e8cfc12f99fe2ff927fd0b2348e9942e9ff778a8c7bb0121037cb184baf184ec6d3f24c927c6243dbb122ae7f9353425f84f3dfcd95fdcd0f90247304402207d541f93bd55ada33a49c3ab2285fd3304e77bb45d108a6ade35bea8561f9c9d02201628a568810d9a00a9463d9bebfac9046feb2aa2bbcf1a1fb967b2648a33c9bd0121037cb184baf184ec6d3f24c927c6243dbb122ae7f9353425f84f3dfcd95fdcd0f900000000")
	assert.Nil(t, err)
	t.Log(txID)
}

func TestTrans2(t *testing.T) {
	s := utNewBTCServer(t)

	wpi, err := s.GenUnsignedTx4TransToOne("hoho2", "tb1q88h06gxq5cgpy566832d6qzcyaey7xhddl2het",
		"tb1qkd8hetlagq3ylg88zupqsu9w6rqykw5mvthhr4", 8000, 2, "")
	assert.Nil(t, err)
	t.Log(wpi)
}

func TestTrans3(t *testing.T) {
	s := utNewBTCServer(t)

	wpi, err := s.GenUnsignedTx4Gather("hoho2", []string{
		"tb1qpll2cts2v39djv6jqq5f0y29jwun454mwt7z57", "tb1qkd8hetlagq3ylg88zupqsu9w6rqykw5mvthhr4"},
		2, "mws4UFRP8XE8JhweXhgyMGkVPZfCMSFgmx")
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
			Amount:  60000,
		},
	}, 2, "2N59MZ6kPV1qWvahhUzDzZGXo4ZsjAmF14i")
	assert.Nil(t, err)
	t.Log(wpi)
}
