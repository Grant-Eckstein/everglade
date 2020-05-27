package everglade

import (
	"crypto/sha256"
)

func (o Object) Hash(salt, pt []byte) []byte {
	saltedHash := sha256.Sum256(append(salt, pt...))
	return saltedHash[:]
}
