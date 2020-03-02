package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type test struct {
	Name string
	Age  int
}

/*
	Challenge 1 [COMPLETED]
		1 - Encrypt text
	 	2 - Write gcmkey to file
		3 - Read gcmkey from file
		4 - Decrypt text with imported key

	Challenge 2
		- Execute challenge 1 with keyList

	Challenge 3
		- Execute challenge 2 with FileList

	Then implment with file encryption
*/
func main() {
	pt := "Hello, world!"

	/* Challenge 1 */
	fmt.Printf("Encoding text [%+v]\n", pt)
	// Encrypt
	key := GenerateGMCKey()
	ct := key.encrypt([]byte(pt))

	// Write to file
	kJ, _ := json.Marshal(&key)
	err := ioutil.WriteFile("testfile", kJ, 0644)
	check(err)

	// Import key
	iKeyJson, err := ioutil.ReadFile("testfile")
	check(err)

	var iKey GCMKey
	err = json.Unmarshal(iKeyJson, &iKey)
	check(err)

	// Decrypt
	result := iKey.decrypt(ct)
	fmt.Printf("Recovered text [%+s] using imported key\n", result)


	/* Challenge 3ish? */
	files := DiscoverFilesInDirectory(".")
	// fmt.Printf("%+v\n", files)

	// Jsonify
	fJ, err := json.Marshal(&files)
	check(err)
	fmt.Printf("fj - [%+s]\n", fJ)

	// Write to file
	err = ioutil.WriteFile("fileList", fJ, 0644)
	check(err)

	// Import key
	iFileListJson, err := ioutil.ReadFile("fileList")
	check(err)

	var iFileList FileList
	err = json.Unmarshal(iFileListJson, &iFileList)
	check(err)

	// err = ioutil.WriteFile("testdata", kJ, 0644)
	// check(err)
	//
	// fmt.Printf("files [%+v]\n", fJ)


	// // Decrypt
	// dKey := GMCKey(key.NonceLength, key.Nonce, key.Key, key.Salt, key.CbcKey.Key, key.CbcKey.Iv)
	// result := dKey.decrypt(ct)
	//
	// fmt.Println(string(result))
	//
	// sJson, _ := json.Marshal(&dKey)
	// fmt.Println(string(sJson))
	//
	// eKey := GCMKey{}
	// json.Unmarshal(sJson, &eKey)
	//
	// pt = "Test2"
	// ct = eKey.encrypt([]byte(pt))
	// result = eKey.decrypt(ct)
	// fmt.Println(string(result))

}
