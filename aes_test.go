package everglade

import (
	"bytes"
	"testing"
)

func TestObject_EncryptCBC_DecryptCBC(t *testing.T) {
	obj := New()
	mesg := []byte("Hello, world!")

	iv, ct := obj.EncryptCBC(mesg)

	pt := obj.DecryptCBC(iv, ct)

	if !bytes.Equal(pt, mesg) {
		t.Fatal("Error verifying CBC encryption process")
	}
}

func TestObject_EncryptGCM_DecryptGCM(t *testing.T) {
	obj := New()
	mesg := []byte("Hello, world!")
	ad := []byte("From grant")

	n, ct := obj.EncryptGCM(mesg, ad)

	pt := obj.DecryptGCM(n, ct, ad)

	if !bytes.Equal(pt, mesg) {
		t.Fatal("Error verifying GCM encryption process")
	}
}
