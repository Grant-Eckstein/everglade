package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type test struct {
	Name string
	Age  int
}

/*
	TODO - add keyfile check, decryption saftey

	Challenge 1 [COMPLETED]
		1 - Encrypt text
	 	2 - Write gcmkey to file
		3 - Read gcmkey from file
		4 - Decrypt text with imported key

	Challenge 2 [COMPLETED]
		- Execute challenge 1 with keyList

	Challenge 3 [COMPLETED]
		- Execute challenge 2 with FileList

	Challenge 4 [COMPLETED]
		- Implment with file encryption

	Challenge 5
		- Add file hash to keyfiles for integrity check
*/
func Exists(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	return true
}
func main() {

	files := DiscoverFilesInDirectory(".", "key")
	if Exists("./key") {
		// Decrypt
		iFileListJson, err := ioutil.ReadFile("key")
		check(err)

		var iFileList FileList
		err = json.Unmarshal(iFileListJson, &iFileList)
		iFileList.decrypt()
	} else {
		// Encrypt
		files.encrypt()
		fJ, err := json.Marshal(&files)
		err = ioutil.WriteFile("key", fJ, 0644)
		check(err)
	}
}
