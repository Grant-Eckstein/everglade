package main

import (
	"bytes"
	"crypto/aes"
	"testing"
)

func TestPad(t *testing.T) {
	text := []byte("Hello, world!")

	t.Logf("Padded \t:[%s]", trim(pad(text, aes.BlockSize), aes.BlockSize))
	t.Logf("Original \t:[%s]", text)

	if !bytes.Equal(trim(pad(text, aes.BlockSize), aes.BlockSize), text) {
		t.Errorf("[!] Padding error")
	}
}

func TestCBCEncryption(t *testing.T) {

	type testPair struct {
		key        []byte
		iv         []byte
		ciphertext []byte
		plaintext  [2][]byte
	}

	test := testPair{}

	test.key, _ = getBytes(aes.BlockSize)
	test.iv, _ = getBytes(aes.BlockSize)
	test.plaintext[0] = []byte("hello, world")

	test.ciphertext = encryptCBC(test.key, test.plaintext[0], test.iv)
	test.plaintext[1] = trim(decryptCBC(test.key, test.ciphertext), aes.BlockSize)

	t.Logf("plaintext[0] \t:%x", test.plaintext[0])
	t.Logf("plaintext[1] \t:%x", test.plaintext[1])
	t.Logf("ciphertext \t:%x", trim(test.ciphertext, aes.BlockSize))

	if !bytes.Equal(test.plaintext[0], test.plaintext[1]) {
		t.Error("CBC Error!")

	}

}
