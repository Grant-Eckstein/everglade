package everglade

import (
	"bytes"
	"testing"
)

func TestObject_Export_Import(t *testing.T) {
	eo := New()
	mesg := []byte("Hello, world!")
	err, salt := Bytes(32)
	if err != nil {
		t.Fatal(getError(bytesGeneration, err))
	}

	eh := eo.Hash(salt, mesg)

	j := eo.Export()

	io := Import(j)
	ih := io.Hash(salt, mesg)

	if !bytes.Equal(ih, eh) {
		t.Fatal("Error verifying export process")
	}
}
