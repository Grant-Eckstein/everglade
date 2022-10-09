package everglade

import (
	"os"
	"path/filepath"
	"strings"
)

// DiscoverFilesInDirectory automates discovery of files in directory and returns []string
func DiscoverFilesInDirectory(dir, ex string) ([]string, error) {
	var fl []string
	err := filepath.Walk(dir,
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && p != ex {
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
	fs, err := DiscoverFilesInDirectory(dir, "")
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

// FindFilesByTypeInDirectory returns the relative path of all files in the directory of a specific extension
func FindFilesByTypeInDirectory(dir, ex string) ([]string, error) {
	var r []string
	fs, err := DiscoverFilesInDirectory(dir, "")
	if err != nil {
		return nil, err
	}

	for _, f := range fs {
		path := strings.Split(f, string(os.PathSeparator))
		name := path[len(path)-1]
		ext := strings.Split(name, string(os.PathSeparator))

		if ext[len(ext)-1] == ex {
			r = append(r, f)
		}
	}
	return r, nil
}
