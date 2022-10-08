package everglade

import "testing"

func TestNew(t *testing.T) {
	_, err := New(1)
	if err != nil {
		t.Fatal(err)
	}
}
