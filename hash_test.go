package everglade

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestObject_Hash(t *testing.T) {
	obj := New()

	message := []byte("Hello, world!")
	salt, err := hex.DecodeString("4451da5cd759ffa88e7802ce7722aa10")
	if err != nil {
		t.Fatalf("Error decoding precomputed salt: %v", err)
	}

	ver_hash, err := hex.DecodeString("d6a50a85127c087e6f7ee732bda4c9345b82a81c5a5b5ec05409cb51f48902c7")
	if err != nil {
		t.Fatalf("Error decoding precomputed hash: %v", err)
	}

	hash := obj.Hash(salt, message)
	if !bytes.Equal(hash, ver_hash) {
		t.Fatalf("Error comparing test to precomputed hash: %v", getError("Hash function doesn't match", nullError))
	}
}
