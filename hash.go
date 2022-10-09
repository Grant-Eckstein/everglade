package everglade

import "crypto/sha256"

func Hash(pt []byte) []byte {
	s := sha256.Sum256(pt)
	return s[:]
}
