package btctx

import (
	"encoding/hex"
	"encoding/json"
)

type Input struct {
	TxID         string
	VOut         uint32
	Address      string
	Amount       int64
	RedeemScript string
}

type Output struct {
	Address string
	Amount  int64
}

type UnsignedTx struct {
	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`
}

func MarshalUnsignedTx(tx *UnsignedTx) (hexD string, err error) {
	d, err := json.Marshal(tx)
	if err != nil {
		return
	}

	hexD = hex.EncodeToString(d)

	return
}

func UnmarshalUnsignedTx(hexD string) (tx *UnsignedTx, err error) {
	d, err := hex.DecodeString(hexD)
	if err != nil {
		return
	}

	tx = &UnsignedTx{}

	err = json.Unmarshal(d, tx)
	if err != nil {
		return
	}

	return
}
