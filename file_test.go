package everglade

import (
	"bytes"
	"os"
	"testing"
)

func TestEverglade_Add(t *testing.T) {
	e, err := New(1)
	if err != nil {
		t.Fatal(err)
	}

	err = e.Add("everglade.go")
	if err != nil {
		t.Fatal(err)
	}
}

func TestEverglade_EncryptDecrypt(t *testing.T) {
	pt, err := os.ReadFile("everglade.go")
	if err != nil {
		t.Fatal(err)
	}

	e, err := New(1)
	if err != nil {
		t.Fatal(err)
	}

	err = e.Add("everglade.go")
	if err != nil {
		t.Fatal(err)
	}

	err = e.Encrypt()
	if err != nil {
		t.Fatal(err)
	}

	err = e.Decrypt()
	if err != nil {
		t.Fatal(err)
	}

	ot, err := os.ReadFile("everglade.go")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(pt, ot) {
		t.Fatal("Output did not match input")
	}

}

func TestEverglade_EncryptDecryptLayers(t *testing.T) {
	layers := 5
	testFilename := "everglade.go"

	// Read in plaintext
	pt, err := os.ReadFile(testFilename)
	if err != nil {
		t.Fatal(err)
	}

	e, err := New(layers)
	if err != nil {
		t.Fatal(err)
	}

	err = e.Add(testFilename)
	if err != nil {
		t.Fatal(err)
	}

	err = e.Encrypt()
	if err != nil {
		t.Fatal(err)
	}

	err = e.Decrypt()
	if err != nil {
		t.Fatal(err)
	}

	// Read in output text
	ot, err := os.ReadFile(testFilename)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(pt, ot) {
		t.Fatal("Output did not match input")
	}
}
