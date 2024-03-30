package edfile

import (
	"crypto/md5" // nolint: gosec
	erand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"os"

	"github.com/sgostarter/i/commerr"
	"github.com/sgostarter/libeasygo/crypt/aes"
	"github.com/sgostarter/libeasygo/pathutils"
)

func EncodePlainFile(d []byte) []byte {
	r := make([]byte, len(d)+100)

	_, _ = erand.Read(r)

	startPos := rand.Intn(60) // nolint: gosec

	binary.LittleEndian.PutUint32(r[10:], uint32(len(d)))
	binary.LittleEndian.PutUint32(r[14:], uint32(startPos))

	copy(r[startPos+20:], d)

	return r
}

func DecodePlainFile(d []byte) (dd []byte, ok bool) {
	if len(d) < 100 {
		return
	}

	dLen := binary.LittleEndian.Uint32(d[10:])
	dPos := binary.LittleEndian.Uint32(d[14:]) + 20

	if dPos+dLen >= uint32(len(d)) {
		return
	}

	dd = make([]byte, dLen)
	copy(dd[:], d[dPos:])

	ok = true

	return
}

func deriveSecKeyFromKeyS(key string) []byte {
	sum := md5.Sum([]byte(key)) // nolint: gosec

	return sum[:]
}

func WriteSecFile(name string, key string, data []byte) (err error) {
	ed, err := aes.CBCEncrypt(EncodePlainFile(data), deriveSecKeyFromKeyS(key))
	if err != nil {
		return
	}

	_ = pathutils.MustDirOfFileExists(name)

	err = os.WriteFile(name, ed, 0600)

	return
}

func ReadSecFile(name string, key string) (data []byte, err error) {
	d, err := os.ReadFile(name)
	if err != nil {
		return
	}

	dd, err := aes.CBCDecrypt(d, deriveSecKeyFromKeyS(key))
	if err != nil {
		return
	}

	data, ok := DecodePlainFile(dd)
	if !ok {
		err = commerr.ErrBadFormat

		return
	}

	return
}
