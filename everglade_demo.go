package main

import (
	"fmt"
	"encoding/json"
)

type test struct {
	Name string
	Age int
}


func main() {
	pt := "Hello, world!"

	// Encrypt
	key := GenerateGMCKey()
	ct := key.encrypt([]byte(pt))

	// Decrypt
	dKey := GMCKey(key.NonceLength, key.Nonce, key.Key, key.Salt, key.CbcKey.Key, key.CbcKey.Iv)
	result := dKey.decrypt(ct)

	fmt.Println(string(result))

	sJson, _ := json.Marshal(&dKey)
	fmt.Println(string(sJson))

	eKey := GCMKey{}
	json.Unmarshal(sJson, &eKey)

	pt = "Test2"
	ct = eKey.encrypt([]byte(pt))
	result = eKey.decrypt(ct)
	fmt.Println(string(result))


}
