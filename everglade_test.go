package main

import (
	"bytes"
	"crypto/aes"
	"testing"
)

func TestPad(t *testing.T) {
	text := []byte("Hello, world!")

	t.Logf("Padded \t:[%x]", pad(text, aes.BlockSize))
	t.Logf("Original \t:[%x]", text)

	if !bytes.Equal(trim(pad(text, aes.BlockSize), aes.BlockSize), text) {
		t.Errorf("[!] Padding error")
		t.Logf("%x\n", trim(pad(text, aes.BlockSize), aes.BlockSize))
	}
}

func TestCBCEncryption(t *testing.T) {

	type testPair struct {
		key        CBCKey
		ciphertext []byte
		plaintext  [2][]byte
	}

	test := testPair{}
	test.key = NewCBCKey()
	test.plaintext[0] = []byte("hello, world")

	test.ciphertext = encryptCBC(test.key.key, test.plaintext[0], test.key.iv)
	test.plaintext[1] = trim(decryptCBC(test.key.key, test.ciphertext), aes.BlockSize)

	t.Logf("plaintext[0] \t:%x", test.plaintext[0])
	t.Logf("plaintext[1] \t:%x", test.plaintext[1])
	t.Logf("ciphertext \t:%x", trim(test.ciphertext, aes.BlockSize))

	if !bytes.Equal(test.plaintext[0], test.plaintext[1]) {
		t.Error("CBC Error\n")

	}

}

func TestGCMEncryption(t *testing.T) {

	type testPair struct {
		key        GCMKey
		ciphertext []byte
		plaintext  [2][]byte
	}

	test := testPair{}

	test.key = GenerateGMCKey()
	test.plaintext[0] = []byte("test")
	test.ciphertext = test.key.encrypt(test.plaintext[0])
	test.plaintext[1] = decryptGCM(test.key, test.ciphertext)

	t.Logf("plaintext[0] \t:%x", test.plaintext[0])
	t.Logf("plaintext[1] \t:%x", test.plaintext[1])
	t.Logf("ciphertext \t:%x", test.ciphertext)

	if !bytes.Equal(test.plaintext[0], test.plaintext[1]) {
		t.Errorf("GCM error\n")
	}
}
