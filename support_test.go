package everglade

import (
	"bytes"
	"testing"
)

func TestBytes(t *testing.T) {
	err, b := Bytes(10)
	if err != nil {
		t.Fatalf("No bytes have been produced: %v", err)
	}
	if bytes.Equal(b, make([]byte, 10)) {
		t.Fatalf("Error generating random bytes: %v", err)
	}
}
