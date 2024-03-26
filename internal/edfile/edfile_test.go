package edfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	secFileName = "../../datas/seed.dat"
	secKey      = "*&%$%tRR"
)

func TestSave(t *testing.T) {
	err := WriteSecFile(secFileName, secKey,
		[]byte("fxx"))
	assert.Nil(t, err)
}

func TestLoad(t *testing.T) {
	d, err := ReadSecFile(secFileName, secKey)
	assert.Nil(t, err)
	t.Log(string(d))
}
