// nolint
package helper

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseScript2(t *testing.T) {
	d, err := hex.DecodeString("532102f410e07213396b8d6289ca6f1c217380a2787db5a7487b0978bd792cbd32343e21028e7bbb364b64687db98ad39674e70c194aa01ef3da3148791536f25ab8861c0221037cb184baf184ec6d3f24c927c6243dbb122ae7f9353425f84f3dfcd95fdcd0f921037f9e1716a16efb9cf977a628942a5ee193aa29d31d99a6ddbb80d1d7dd69b5da54ae")
	assert.Nil(t, err)

	t.Log(ParseScript(d))
}
