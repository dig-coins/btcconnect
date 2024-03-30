package redistorage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test1(t *testing.T) {
	stg, err := NewRedisStorage()
	assert.Nil(t, err)

	block, txIndex, err := stg.LoadProcessedBlockPosition()
	assert.Nil(t, err)
	t.Log(block, txIndex)

	err = stg.StoreProcessedBlockPosition(1, 0)
	assert.Nil(t, err)

	block, txIndex, err = stg.LoadProcessedBlockPosition()
	assert.Nil(t, err)
	assert.EqualValues(t, 1, block)
	assert.EqualValues(t, 0, txIndex)

	scriptPubKeyType, scriptPubKeyHex, address, err := stg.GetUnspentTXO("t1", 0)
	assert.Nil(t, err)
	t.Log(scriptPubKeyType, scriptPubKeyHex, address)

	err = stg.NewTXO("t1", 0, "scriptPubKeyType", "scryptPubKeyHex", "address")
	assert.Nil(t, err)

	scriptPubKeyType, scriptPubKeyHex, address, err = stg.GetUnspentTXO("t1", 0)
	assert.Nil(t, err)
	assert.EqualValues(t, "scriptPubKeyType", scriptPubKeyType)
	assert.EqualValues(t, "scryptPubKeyHex", scriptPubKeyHex)
	assert.EqualValues(t, "address", address)

	err = stg.SpendTXO("t1", 0)
	assert.Nil(t, err)

	err = stg.SpendTXO("t1", 0)
	assert.Nil(t, err)
}
