package helper

import (
	"encoding/json"

	"github.com/0xb10c/rawtx"
	"github.com/okx/go-wallet-sdk/coins/bitcoin"
)

func ParseTx(txHex string) (txInfo string, err error) {
	wireTx, err := bitcoin.NewTxFromHex(txHex)
	if err != nil {
		return
	}

	tx := rawtx.Tx{}
	tx.FromWireMsgTx(wireTx)

	txStats := tx.Stats()

	txStatJSON, err := json.MarshalIndent(txStats, "", "  ")
	if err != nil {
		return
	}

	txInfo = string(txStatJSON)

	return
}

func ParseScript(scripts []byte) string {
	return rawtx.BitcoinScript(scripts).Parse().String()
}
