package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/blake2b"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// CBCKey organizes CBC keys with nonce
type CBCKey struct {
	key []byte
	iv  []byte
}

// GCMKey organizes GCM keys with nonce
type GCMKey struct {
	key         []byte
	keyID       int
	salt        []byte
	nonce       []byte
	nonceLength int
	cbcKey      CBCKey
}

/* --- Type Constructors --- */

// NewCBCKey initializes CBCKey with random key and IV
func NewCBCKey() CBCKey {
	cbcKey := CBCKey{}

	cbcKey.key, _ = getBytes(aes.BlockSize)
	cbcKey.iv, _ = getBytes(aes.BlockSize)
	return cbcKey
}

// NewGMCKey initializes GMCKey with key, salt, and nonce
func NewGMCKey() GCMKey {
	gcmKey := GCMKey{}
	gcmKey.nonceLength = 12

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	gcmKey.nonce, _ = getBytes(gcmKey.nonceLength)
	gcmKey.key, gcmKey.salt = getPass()

	gcmKey.cbcKey = NewCBCKey()

	return gcmKey
}

/* --- Encryption Methods --- */

func encryptCBC(key []byte, plaintext []byte, iv []byte) []byte {
	block, err := aes.NewCipher(pad(key, aes.BlockSize))
	check(err)
	if len(plaintext)%aes.BlockSize != 0 {
		plaintext = pad(plaintext, block.BlockSize())
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	bm := cipher.NewCBCEncrypter(block, iv)
	ciphertext = append(iv, make([]byte, len(plaintext))...)
	bm.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext
}

func decryptCBC(key []byte, ciphertext []byte) []byte {
	block, err := aes.NewCipher(pad(key, aes.BlockSize))
	check(err)
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	bm := cipher.NewCBCDecrypter(block, iv)
	bm.CryptBlocks(ciphertext, ciphertext)
	return ciphertext
}

// @DEBUG Prints each element of GCMKey #TODO automate this
func (k *GCMKey) print() {
	fmt.Printf("---Start Key---\n")
	fmt.Printf("KeyID :%v\n", k.keyID)
	fmt.Printf("Salt :%x\n", k.salt)
	fmt.Printf("Nonce :%x\n", k.nonce)
	fmt.Printf("Key :%x\n", k.key)
	fmt.Printf("\t---Start SubKey---\n")
	fmt.Printf("\tIV :%x\n", k.cbcKey.iv)
	fmt.Printf("\tKey :%x\n", k.cbcKey.key)
	fmt.Printf("\t---End SubKey---\n")
	fmt.Printf("---End Key---\n")

}

// Encrypts authenticated data
func encryptGCM(gcmKey GCMKey, pt []byte) []byte {
	block, err := aes.NewCipher(gcmKey.key)
	check(err)

	aesgcm, err := cipher.NewGCM(block)
	check(err)

	// Using the encrypted nonce for the Authenticated Data
	ad := encryptCBC(gcmKey.cbcKey.key, gcmKey.nonce, gcmKey.cbcKey.iv)

	var ciphertext []byte
	ciphertext = aesgcm.Seal(ciphertext, gcmKey.nonce, pt, ad)
	return ciphertext
}

// Decrypts authenticated data
func decryptGCM(gcmKey GCMKey, ct []byte) []byte {

	block, err := aes.NewCipher(gcmKey.key)
	check(err)

	aesgcm, err := cipher.NewGCM(block)
	check(err)

	// Using the encrypted nonce for the Authenticated Data
	ad := encryptCBC(gcmKey.cbcKey.key, gcmKey.nonce, gcmKey.cbcKey.iv)

	var pt []byte
	pt, err = aesgcm.Open(pt, gcmKey.nonce, ct, ad)
	check(err)

	return pt
}

// Pads data
func pad(b []byte, blocksize int) []byte {
	if len(b)%blocksize == 0 {
		return b
	}
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return []byte(pb)
}

// Trims padding
func trim(b []byte, blocksize int) []byte {
	c := b[len(b)-1]
	n := int(c)
	return b[:len(b)-n]
}

// Returns random bytes of n length
func getBytes(n int) ([]byte, error) {
	k := make([]byte, n)
	_, err := rand.Read(k)
	return k, err
}

// Fixes the typing error [n]byte -> []byte
func getHash(dat []byte) []byte {
	hash := blake2b.Sum256(dat)
	var r = hash[:]
	return r
}

// Return salted key and salt
func getPass() ([]byte, []byte) {
	k, _ := getBytes(aes.BlockSize)
	s, _ := getBytes(aes.BlockSize / 2)
	sk := append(k, s...)

	return getHash(sk), s
}
