package everglade

import "testing"

func TestEverglade_Add(t *testing.T) {
	e, err := New()
	if err != nil {
		t.Fatal(err)
	}

	err = e.Add("everglade.go")
	if err != nil {
		t.Fatal(err)
	}

}
