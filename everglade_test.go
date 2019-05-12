package main

import (
	"bytes"
	"crypto/aes"
	"testing"
)

/*
// Test error handling
if err != nil {
		t.Errorf("[!] Caught error [%s]", err)
	}
*/

func TestPad(t *testing.T) {
	text := []byte("Hello, world!")
	if !bytes.Equal(trim(pad(text, aes.BlockSize), aes.BlockSize), text) {
		t.Errorf("[!] Padding error")
		t.Logf("Padded \t:[%s]", trim(pad(text, aes.BlockSize), aes.BlockSize))
		t.Logf("Original \t:[%s]", text)
	}
}

func TestEncryption(t *testing.T) {
	key, _ := getBytes(aes.BlockSize)
	plaintext, _ := getBytes(aes.BlockSize * 3)

	t.Logf("Key \t :%x", key)
	t.Logf("Plaintext :%x", plaintext)

	ciphertext := encrypt(key, plaintext)
	if bytes.Equal(plaintext, decrypt(key, ciphertext)) != true {
		t.Errorf("[!] Crypto error")
	}

}

type SubKeyTestData struct {
	key        []byte
	plaintext  []byte
	ciphertext []byte
}

type SubKeyTest struct {
	keys []SubKeyTestData
}

func NewSubKeyTest(n int) SubKeyTest {

	masterKey, _ := getPass()
	keySession := NewFileKey(masterKey)

	subKeyTest := SubKeyTest{}

	for i := 1; i <= n; i++ {
		key := keySession.getKey()[:aes.BlockSize]
		plaintext, _ := getBytes(aes.BlockSize)
		ciphertext := encrypt(key, plaintext)

		subKeyTest.keys = append(subKeyTest.keys, SubKeyTestData{key, plaintext, ciphertext})
	}
	return subKeyTest
}

/*
 * Here I
 * 1 - Create 2 subkeys
 * 2 - encrypt string 2 times
 * 3 - export master key
 * 4 - import master key
 * 5 - decrypt string 2 times
 */
func TestSubKeys(t *testing.T) {

	// Collect tests
	subKeys := NewSubKeyTest(2)
	importedKeys := SubKeyTest{}

	// Import tests
	for _, test := range subKeys.keys {
		importedKeys.keys = append(importedKeys.keys, test)
	}

	// Display tests
	for _, test := range importedKeys.keys {
		t.Logf("---")
		t.Logf("Key \t\t:%x", test.key)
		t.Logf("Plaintext \t:%x", test.plaintext)
		t.Logf("Ciphertext \t:%x", test.ciphertext)
	}
}
