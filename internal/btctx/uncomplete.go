package btctx

import (
	"encoding/hex"
	"encoding/json"
)

type MultiSignInfo struct {
	PublicKeys []string
	MinSignNum int
	SignedNum  int
}

type Uncompleted struct {
	MultiSignInfos map[int]MultiSignInfo
}

func NewUncompleted() *Uncompleted {
	return &Uncompleted{MultiSignInfos: make(map[int]MultiSignInfo)}
}

func (c *Uncompleted) Completed() bool {
	for _, msi := range c.MultiSignInfos {
		if msi.SignedNum < msi.MinSignNum {
			return false
		}
	}

	return true
}

func MarshalUncompleted(c *Uncompleted) (hexD string, err error) {
	d, err := json.Marshal(c)
	if err != nil {
		return
	}

	hexD = hex.EncodeToString(d)

	return
}

func UnmarshalUncompleted(hexD string) (c *Uncompleted, err error) {
	d, err := hex.DecodeString(hexD)
	if err != nil {
		return
	}

	c = NewUncompleted()

	err = json.Unmarshal(d, c)

	return
}
