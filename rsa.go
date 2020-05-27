package everglade

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

func (o Object) EncryptOAEP(pt, l []byte) (error, []byte) {
	ct, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &o.RSAKey.PublicKey, pt, l)
	if err != nil {
		return getError(rsaEncryptionProcessing, err), nil
	}

	return nil, ct
}

func (o Object) DecryptOAEP(ct, l []byte) (error, []byte) {
	pt, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, o.RSAKey, ct, l)
	if err != nil {
		return getError(rsaDecryptionProcessing, err), nil
	}

	return nil, pt
}

// Sign data with RSASSA-PKCS1-V1_5-SIGN from RSA PKCS#1 v1.5 using o.RSAKey
func (o Object) Sign(pt []byte) (error, []byte) {

	hashed := sha256.Sum256(pt)

	signature, err := rsa.SignPKCS1v15(rand.Reader, o.RSAKey, crypto.SHA256, hashed[:])
	if err != nil {
		return getError(rsaSignatureProcessing, err), signature
	} else {
		return nil, signature
	}
}

// Verify a signature with RSASSA-PKCS1-V1_5-SIGN from RSA PKCS#1 v1.5 using o.RSAKey
func (o Object) Verify(pt, s []byte) (error, bool) {
	if len(o.BlockKey) != 0 {
		err, newSignature := o.Sign(pt)
		if err != nil {
			return err, false
		}
		if bytes.Equal(newSignature, s) {
			return nil, true
		}
		return getError(invalidSignature, err), false
	} else {
		return getError(noKeySet, nil), false
	}
}
