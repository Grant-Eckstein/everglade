package everglade

import (
	"crypto/rand"
)

func Bytes(n int) (error, []byte) {
	b := make([]byte, n)
	err := randBytes(b)
	if err != nil {
		return getError(hashGeneration, err), nil
	} else {
		return nil, b
	}
}

func randBytes(b []byte) error {
	_, err := rand.Read(b)
	return err
}
