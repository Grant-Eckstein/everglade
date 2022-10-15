package everglade

import (
	"bytes"
	"os"
)

type Signature []byte

func (e *Everglade) Sign() ([]Signature, error) {
	var sigs []Signature
	for _, f := range e.Paths {
		// Get file contents
		fc, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}

		// Sign each file
		sig, err := e.Blind.RSA.Sign(fc)
		if err != nil {
			return nil, err
		}

		// Add signature to return
		sigs = append(sigs, sig)
	}
	return sigs, nil
}

func (e *Everglade) Verify(sig, pt []byte) bool {
	// Create verification signature
	n, err := e.Blind.RSA.Sign(pt)
	if err != nil {
		return false
	}

	// If they do not match, the signature is invalid
	if !bytes.Equal(sig, n) {
		return false
	}
	
	return true
}
