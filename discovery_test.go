package everglade

import (
	"os"
	"path/filepath"
	"testing"
)

func createTestArea(filename string) (string, error) {
	content := []byte("Hello, world!")
	// Create test directory
	dir, err := os.MkdirTemp("", "testing")
	if err != nil {
		return "", err
	}

	// Write to temp file
	fn := filepath.Join(dir, filename)
	err = os.WriteFile(fn, content, 777)
	if err != nil {
		return "", err
	}

	return dir, nil
}

func TestDiscoverFilesInDirectory(t *testing.T) {
	dir, err := createTestArea("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	fs, err := DiscoverFilesInDirectory(dir)
	if err != nil {
		return
	}

	if len(fs) != 1 {
		t.Fatal("Discovery failed, should detect 1 file")
	}

	// Clean up
	err = os.RemoveAll(dir)
	if err != nil {
		t.Fatal("error removing testing directory")
	}
}

func TestFindFilesByTypeInDirectory(t *testing.T) {
	// Create test area
	dir, err := createTestArea("test.txt")
	if err != nil {
		t.Fatal(err)
	}

	// Run test
	fs, err := FindFilesByTypeInDirectory(dir, ".txt")
	if err != nil {
		t.Fatal(err)
	}

	// Evaluate outcome
	if len(fs) != 1 {
		t.Logf("Length of fs: %v\n", len(fs))
		t.Logf("Contents: %v\n", fs)
		t.Logf("Dir is '%v'\n", dir)
		t.Fatal("Discovery failed, should detect 1 file")
	}

	// Clean up
	err = os.RemoveAll(dir)
	if err != nil {
		t.Fatal("error removing testing directory")
	}
}

func TestFindFileInDirectory(t *testing.T) {
	// Create test area
	fn := "test.txt"
	dir, err := createTestArea(fn)
	if err != nil {
		t.Fatal(err)
	}

	// Run test
	p, err := FindFileInDirectory(dir, fn)
	if err != nil {
		t.Fatal(err)
	}

	// Evaluate
	tc := filepath.Join(dir, fn)
	if p != tc || p == "" {
		t.Fatal("File not found")
	}

	// Clean up
	err = os.RemoveAll(dir)
	if err != nil {
		t.Fatal("error removing testing directory")
	}
}
