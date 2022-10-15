package everglade

import (
	"encoding/hex"
	"testing"
)

func TestHash(t *testing.T) {
	c := []byte("Hello, world!")
	ex := "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3"

	o := hex.EncodeToString(Hash(c))

	if ex != o {
		t.Logf("Hash: %v\n", o)
		t.Logf("Ex  : %v\n", ex)
		t.Fatal("Error, hash does not match test case")
	}
}
