package everglade

import (
	"errors"
	"fmt"
)

const (
	invalidSignature        = "Invalid signature error"
	noKeySet                = "No key set in object"
	rsaKeyGeneration        = "Error generating RSA Key"
	rsaSignatureProcessing  = "Error generating signature"
	rsaEncryptionProcessing = "Error encrypting data with RSA"
	rsaDecryptionProcessing = "Error decrypting data with RSA"
	hashGeneration          = "Error generating hash"
	bytesGeneration         = "Error generating random bytes"
	cipherBlockGeneration   = "Error generating cipher.Block"
	gcmCipherGeneration     = "Error generating GCM cipher"
	nonceWrongLength        = "Nonce is incorrect size"
	nonceNil      = "Nonce is nil"
	jsonMarshal   = "Error marshaling data"
	jsonUnmarshal = "Error unmarshaling data"
	fileWalkRead  = "Error reading in file walk"
	fileWalk      = "Error walking directory"
	readFile      = "Error reading file"
	createFile    = "Error creating file"
	writeFile     = "Error writing file"
	deleteFile    = "Error deleting file"
	encryptFile   = "Error encrypting file"
	decryptFile   = "Error decrypting file"
)

var nullError = errors.New("")

func getError(s string, e error) error {
	if s == "" {
		return errors.New(s)
	}
	return errors.New(fmt.Sprintf("%v: %v", s, e))
}
