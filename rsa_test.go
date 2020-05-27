package everglade

import (
	"bytes"
	"testing"
)

func TestObject_EncryptOAEP_DecryptOAEP(t *testing.T) {
	obj := New()
	var mesg = []byte("Hello, world")
	var label = []byte("label")

	err, ct := obj.EncryptOAEP(mesg, label)
	if err != nil {
		t.Fatalf("Error encrypting data: %v", err)
	}

	err, pt := obj.DecryptOAEP(ct, label)
	if err != nil {
		t.Fatalf("Error decrypting data: %v", err)
	}

	if !bytes.Equal(pt, mesg) {
		t.Fatal("Error verifying OAEP encryption process")
	}
}

func TestObject_Sign_Verify(t *testing.T) {
	obj := New()
	var mesg = []byte("Hello, world!")

	err, signature := obj.Sign(mesg)
	if err != nil {
		t.Fatalf("Error creating signature: %v", err)
	}
	err, val := obj.Verify(mesg, signature)
	if err != nil {
		t.Fatalf("Error verifying signature: %v", err)
	} else if val == false {
		t.Fatal("Signature invalid")
	}
}
