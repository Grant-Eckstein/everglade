package everglade

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEverglade_Sign(t *testing.T) {
	// Create test area
	fn := "test.txt"
	dir, err := createTestArea(fn)
	if err != nil {
		t.Fatal(err)
	}
	fn = filepath.Join(dir, fn)

	// Create new everglade instance with new file
	e, err := New(1)
	if err != nil {
		t.Fatal(err)
	}

	// Add test file
	err = e.Add(fn)
	if err != nil {
		t.Fatal(err)
	}

	// Get signatures
	sigs, err := e.Sign()
	if err != nil {
		t.Fatal(err)
	}

	// Verify signature
	c, err := os.ReadFile(fn)
	if err != nil {
		t.Fatal(err)
	}

	g := e.Verify(sigs[0], c)
	if !g {
		t.Fatal("Signature is not valid")
	}

	// Clean
	err = os.RemoveAll(dir)
	if err != nil {
		t.Fatal("error removing testing directory")
	}
}
