package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"


	"golang.org/x/crypto/blake2b"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// CBCKey organizes CBC keys with nonce
type CBCKey struct {
	Key []byte
	Iv  []byte
}

// GCMKey organizes GCM keys with nonce
type GCMKey struct {
	Key         []byte
	KeyID       int
	Salt        []byte
	Nonce       []byte
	NonceLength int
	CbcKey      CBCKey
}

/* --- Type Constructors --- */

// NewCBCKey initializes CBCKey with random key and IV
func NewCBCKey() CBCKey {
	cbcKey := CBCKey{}

	cbcKey.Key, _ = getBytes(aes.BlockSize)
	cbcKey.Iv, _ = getBytes(aes.BlockSize)
	return cbcKey
}

// NewGMCKey initializes GMCKey with key, salt, and nonce
func GenerateGMCKey() GCMKey {
	gcmKey := GCMKey{}
	gcmKey.NonceLength = 12

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	gcmKey.Nonce, _ = getBytes(gcmKey.NonceLength)
	gcmKey.Key, gcmKey.Salt = getPass()

	gcmKey.CbcKey = NewCBCKey()

	return gcmKey
}

func GMCKey(nL int, n []byte, k []byte, s []byte, cK []byte, cI []byte) GCMKey {
	gcmKey := GCMKey{}
	gcmKey.NonceLength = nL

	gcmKey.Nonce = n
	gcmKey.Key = k
	gcmKey.Salt = s

	gcmKey.CbcKey.Key = cK
	gcmKey.CbcKey.Iv = cI

	return gcmKey
}

/* --- Encryption Methods --- */

func (k *CBCKey) encrypt(plaintext []byte) []byte {
	block, err := aes.NewCipher(pad(k.Key, aes.BlockSize))
	check(err)
	if len(plaintext)%aes.BlockSize != 0 {
		plaintext = pad(plaintext, block.BlockSize())
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	bm := cipher.NewCBCEncrypter(block, k.Iv)
	ciphertext = append(k.Iv, make([]byte, len(plaintext))...)
	bm.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext
}

func (k *CBCKey) decrypt(ciphertext []byte) []byte {
	block, err := aes.NewCipher(pad(k.Key, aes.BlockSize))
	check(err)
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	bm := cipher.NewCBCDecrypter(block, iv)
	bm.CryptBlocks(ciphertext, ciphertext)
	return ciphertext
}


// Encrypts authenticated data
func (gcmKey *GCMKey) encrypt(pt []byte) []byte {
	block, err := aes.NewCipher(gcmKey.Key)
	check(err)

	aesgcm, err := cipher.NewGCM(block)
	check(err)

	// Using the encrypted nonce for the Authenticated Data
	ad := gcmKey.CbcKey.encrypt(gcmKey.Nonce)

	var ciphertext []byte
	ciphertext = aesgcm.Seal(ciphertext, gcmKey.Nonce, pt, ad)
	return ciphertext
}

// Decrypts authenticated data
func (gcmKey *GCMKey) decrypt(ct []byte) []byte {

	block, err := aes.NewCipher(gcmKey.Key)
	check(err)

	aesgcm, err := cipher.NewGCM(block)
	check(err)

	// Using the encrypted nonce for the Authenticated Data
	ad := gcmKey.CbcKey.encrypt(gcmKey.Nonce)

	var pt []byte
	pt, err = aesgcm.Open(pt, gcmKey.Nonce, ct, ad)
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
