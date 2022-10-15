package everglade

import (
	"os"
	"path/filepath"
	"strings"
)

// DiscoverFilesInDirectory automates discovery of files in directory and returns []string
func DiscoverFilesInDirectory(dir string) ([]string, error) {
	var fl []string
	err := filepath.Walk(dir,
		// For each file/folder in dir: append to fl if it isn't a directory
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				fl = append(fl, p)
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return fl, nil
}

// FindFileInDirectory takes a directory and a filename and returns the relative path of that file if exists
func FindFileInDirectory(dir, fn string) (string, error) {
	fs, err := DiscoverFilesInDirectory(dir)
	if err != nil {
		return "", err
	}

	for _, f := range fs {
		name := strings.Split(f, string(os.PathSeparator))
		if name[len(name)-1] == fn {
			return f, nil
		}
	}
	return "", nil
}

// FindFilesByTypeInDirectory returns the path of all files in the directory of a specific extension
func FindFilesByTypeInDirectory(dir, ex string) ([]string, error) {
	// Discover all files
	fs, err := DiscoverFilesInDirectory(dir)
	if err != nil {
		return nil, err
	}

	// For each file, compare its extension to the target extension
	var r []string
	for _, f := range fs {
		fileExt := filepath.Ext(f)
		if fileExt == ex {
			r = append(r, f)
		}
	}
	return r, nil
}
