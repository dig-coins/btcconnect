package share

import (
	"encoding/hex"
	"encoding/json"
)

type CommandResult struct {
	ErrMessage    string `json:"errMessage,omitempty" yaml:"errMessage,omitempty"`
	Tx            string `json:"tx,omitempty" yaml:"tx,omitempty"`
	AllSignedFlag bool   `json:"all_signed_flag,omitempty" yaml:"all_signed_flag,omitempty"`
}

func UnmarshalCommandResult(s string) (cr CommandResult, err error) {
	d, err := hex.DecodeString(s)
	if err != nil {
		return
	}

	err = json.Unmarshal(d, &cr)

	return
}

func MarshalCommandResult(cr CommandResult) (s string, err error) {
	d, err := json.Marshal(&cr)
	if err != nil {
		return
	}

	s = hex.EncodeToString(d)

	return
}
