package everglade

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestFile_EncryptCBC_DecryptCBC(t *testing.T) {
	obj := New()
	mesg := []byte("Hello, world!")
	f := NewFile("./test.file")

	// Write test message to file
	w, err := os.Create(f.Name)
	if err != nil {
		t.Fatal(getError(createFile, err))
	}

	defer w.Close()
	_, err = w.Write(mesg)
	if err != nil {
		t.Fatal(getError(writeFile, err))
	}

	// Encrypt file
	err = f.EncryptCBC(obj)
	if err != nil {
		t.Fatal(getError(encryptFile, err))
	}

	// Decrypt file
	err = f.DecryptCBC(obj)
	if err != nil {
		t.Fatal(getError(decryptFile, err))
	}

	// Read test message from file
	d, err := ioutil.ReadFile(f.Name)
	if err != nil {
		t.Fatal(getError(readFile, err))
	}

	// Compare read data after decryption to origional message
	if !bytes.Equal(d, mesg) {
		t.Fatal("Error verifying file CBC encryption process")
	}

	// Cleanup test file
	err = os.Remove(f.Name)
	if err != nil {
		log.Fatal(getError(deleteFile, err))
	}
}

func TestFile_EncryptGCM_DecryptGCM(t *testing.T) {
	obj := New()
	mesg := []byte("Hello, world!")
	ad := []byte("Auth")
	f := NewFile("./test.file")

	// Write test message to file
	w, err := os.Create(f.Name)
	if err != nil {
		t.Fatal(getError(createFile, err))
	}

	defer w.Close()
	_, err = w.Write(mesg)
	if err != nil {
		t.Fatal(getError(writeFile, err))
	}

	// Encrypt file
	err = f.EncryptGCM(obj, ad)
	if err != nil {
		t.Fatal(getError(encryptFile, err))
	}

	// Decrypt file
	err = f.DecryptGCM(obj, ad)
	if err != nil {
		t.Fatal(getError(decryptFile, err))
	}

	// Read test message from file
	d, err := ioutil.ReadFile(f.Name)
	if err != nil {
		t.Fatal(getError(readFile, err))
	}

	// Compare read data after decryption to origional message
	if !bytes.Equal(d, mesg) {
		t.Fatal("Error verifying file GCM encryption process")
	}

	// Cleanup test file
	err = os.Remove(f.Name)
	if err != nil {
		log.Fatal(getError(deleteFile, err))
	}
}

func TestFile_EncryptOAEP_DecryptOAEP(t *testing.T) {
	obj := New()
	mesg := []byte("Hello, world!")
	l := []byte("label")
	f := NewFile("./test.file")

	// Write test message to file
	w, err := os.Create(f.Name)
	if err != nil {
		t.Fatal(getError(createFile, err))
	}

	defer w.Close()
	_, err = w.Write(mesg)
	if err != nil {
		t.Fatal(getError(writeFile, err))
	}

	// Encrypt file
	err = f.EncryptOAEP(obj, l)
	if err != nil {
		t.Fatal(getError(encryptFile, err))
	}

	// Decrypt file
	err = f.DecryptOAEP(obj, l)
	if err != nil {
		t.Fatal(getError(decryptFile, err))
	}

	// Read test message from file
	d, err := ioutil.ReadFile(f.Name)
	if err != nil {
		t.Fatal(getError(readFile, err))
	}

	// Compare read data after decryption to origional message
	if !bytes.Equal(d, mesg) {
		t.Fatal("Error verifying file OAEP encryption process")
	}

	// Cleanup test file
	err = os.Remove(f.Name)
	if err != nil {
		log.Fatal(getError(deleteFile, err))
	}
}