package btctx

import (
	"encoding/hex"
	"encoding/json"
)

type Input struct {
	TxID         string `json:"tx_id"`
	VOut         uint32 `json:"v_out"`
	Address      string `json:"address"`
	Amount       int64  `json:"amount"`
	RedeemScript string `json:"redeem_script"`
}

type Output struct {
	Address    string `json:"address"`
	Amount     int64  `json:"amount"`
	ChangeFlag bool   `json:"change_flag"`
	Comment    string `json:"comment,omitempty"`
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
