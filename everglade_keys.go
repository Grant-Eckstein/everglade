package main

import (
	"crypto/aes"
	"encoding/json"
	"fmt"
)

/* --- New Types --- */

func (k *CBCKey) jsonify() []byte {
	b, err := json.Marshal(k)
	check(err)
	return b
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

	key := GenerateGMCKey()

	key.KeyID = len(kl.keys)

	kl.keys = append(kl.keys, key)

	return key
}

// getKey returns key given ID in KeyList
func (kl *KeyList) getKey(keyID int) GCMKey {
	for _, key := range kl.keys {
		if key.KeyID == keyID {
			return key
		}
	}
	fmt.Printf("KeyID lookup failed\n KeyID :%v, key list length :%v\n", keyID, len(kl.keys))
	panic("No key found with that ID")
}

// printAllKeys returns key given ID in KeyList
func (kl *KeyList) printAllKeys() {
	for _, key := range kl.keys {
		fmt.Printf("KeyID:%v Key:%x\n", key.KeyID, key.Key)
	}
}
