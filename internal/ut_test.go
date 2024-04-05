// nolint
package internal

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/dig-coins/btcconnect/internal/btcserver"
	"github.com/dig-coins/btcconnect/internal/btctx"
	"github.com/dig-coins/btcconnect/internal/share"
	"github.com/dig-coins/btcconnect/internal/utl"
	"github.com/stretchr/testify/assert"
)

var _uTConfig *utl.UTConfig

func TestMain(m *testing.M) {
	_uTConfig = utl.Setup()

	m.Run()
}
func TestBroadcastTx(t *testing.T) {
	s := utNewBTCServer(t)

	txID, err := s.SendRawTransaction("01000000000101d431feaf8d888bd78e22d97c8fd990d1e98b38051473daebe06957068febc6170000000000ffffffff010e0803000000000017a914384e5c254516849a9d98f709d88eca37f5b8cc658702483045022100cd49dfa65757e64ff79cbff78271e58743400c7bc6799f091d2bcfaeaa3de9a102200fd969699b8ed1bc0f55a463711ca0e125c13c8b69d81b1ae344a7f296e275050121030c8bf410a7b008b914c44d4aa88792764fb12b53a6052144db73f4b5e504878b00000000")
	assert.Nil(t, err)
	t.Log(txID)
}

func utNewBTCServer(t *testing.T) *btcserver.BTCServer {
	return btcserver.NewBTCServer(&btcserver.Config{
		CoinType:              share.CoinTypeBTC,
		Listens:               ":9001",
		RPCHost:               _uTConfig.GetBTCRpcHost(t),
		RPCPort:               _uTConfig.GetBTCRpcPort(t),
		RPCUser:               _uTConfig.GetBTCRpcUser(t),
		RPCPassword:           _uTConfig.GetBTCRpcPassword(t),
		RPCUseSSL:             _uTConfig.GetBTCRpcUseTLS(t),
		MultiSignAddressInfos: nil,
	}, nil, nil)
}

func TestTrans3(t *testing.T) {
	s := utNewBTCServer(t)

	unsignedTx, err := s.GenUnsignedTx4Gather(_uTConfig.GetWallet(t), []string{
		"bc1qxvqmxrnnfqtwz42pl3vr7ahfu3axlslf4809f2"},
		40, "36pjdq8Nb7XL2AuerdSeBxr6yazKomLgdV", true) // 4,21
	assert.Nil(t, err)

	tx, err := btctx.UnmarshalUnsignedTx(unsignedTx)
	assert.Nil(t, err)

	var totalAmount, outputAmount int64

	for _, input := range tx.Inputs {
		totalAmount += input.Amount
	}

	for _, output := range tx.Outputs {
		outputAmount += output.Amount
	}

	t.Logf("In: %d, Out: %d, Fee: %d\n", totalAmount, outputAmount, totalAmount-outputAmount)
	utSaveUnsignedTxToFile(t, unsignedTx)
}

func utSaveUnsignedTxToFile(t *testing.T, unsignedTx string) {
	d, err := json.Marshal(share.Command{
		CommandType: share.CommandTypeGenTx,
		Input:       unsignedTx,
	})
	assert.Nil(t, err)

	err = os.WriteFile("../tmp/unsigned-tx.json", d, 0600)
	assert.Nil(t, err)
}
