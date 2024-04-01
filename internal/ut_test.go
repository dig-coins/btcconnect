// nolint
package internal

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/dig-coins/btcconnect/internal/btcserver"
	"github.com/dig-coins/btcconnect/internal/btctx"
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
func TestBroadcastTx(t *testing.T) {
	s := utNewBTCServer(t)

	txID, err := s.SendRawTransaction("01000000000101ff3bdd0e728614edcb78bf9e386d1e104fc2aa86a3d8dade792a686c172f668400000000171600143301b30e734816e15541fc583f76e9e47a6fc3e9ffffffff01dd9a07000000000017a914384e5c254516849a9d98f709d88eca37f5b8cc65870247304402202f13d4f835c8323ca5a14ceba66f311c077f3e504fe08601a9699f4dbf443a800220734ed885178d57548455880a47c2c9928324f832f441ae7698fca249d1d85cc70121030c8bf410a7b008b914c44d4aa88792764fb12b53a6052144db73f4b5e504878b00000000")
	assert.Nil(t, err)
	t.Log(txID)
}

func utNewBTCServer(t *testing.T) *btcserver.BTCServer {
	stg, err := redistorage.NewRedisStorage()
	assert.Nil(t, err)

	return btcserver.NewBTCServer(&btcserver.Config{
		RPCHost:     _uTConfig.GetBTCRpcHost(t),
		RPCPort:     _uTConfig.GetBTCRpcPort(t),
		RPCUser:     _uTConfig.GetBTCRpcUser(t),
		RPCPassword: _uTConfig.GetBTCRpcPassword(t),
		RPCUseSSL:   _uTConfig.GetBTCRpcUseTLS(t),
	}, stg, &chaincfg.MainNetParams, nil, nil)
}

func TestTrans3(t *testing.T) {
	s := utNewBTCServer(t)

	unsignedTx, err := s.GenUnsignedTx4Gather(_uTConfig.GetWallet(t), []string{
		"34KAAvz5Nad7G5wk56PKjPpVLfMjjkaeQb"},
		40, "36pjdq8Nb7XL2AuerdSeBxr6yazKomLgdV") // 4,21
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
