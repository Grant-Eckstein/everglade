package everglade

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"log"
)

func PKCS5Pad(data []byte, bs int) []byte {
	p := bs - len(data)%bs
	r := bytes.Repeat([]byte{byte(p)}, p)
	return append(data, r...)
}

func PKCS5Trim(data []byte) []byte {
	p := data[len(data)-1]
	return data[:len(data)-int(p)]
}

func (o Object) EncryptCBC(pt []byte) ([]byte, []byte) {
	pt = PKCS5Pad(pt, aes.BlockSize)

	err, iv := Bytes(aes.BlockSize)
	if err != nil {
		log.Fatalf("Error generating CBC initialization vector: %v", err)
	}

	b, err := aes.NewCipher(o.BlockKey)
	if err != nil {
		log.Fatal(getError(cipherBlockGeneration, nullError))
	}

	ct := make([]byte, len(pt))
	m := cipher.NewCBCEncrypter(b, iv)
	m.CryptBlocks(ct, pt)

	return iv, ct
}

func (o Object) DecryptCBC(iv, ct []byte) []byte {
	b, err := aes.NewCipher(o.BlockKey)
	if err != nil {
		log.Fatal(getError(cipherBlockGeneration, nullError))
	}

	pt := make([]byte, len(ct))
	m := cipher.NewCBCDecrypter(b, iv)
	m.CryptBlocks(pt, ct)

	return PKCS5Trim(pt)
}

// EncryptGCM using AES-GCM-256 with RSA Signature of component of RSA private key as associated data
func (o Object) EncryptGCM(pt, ad []byte) ([]byte, []byte) {

	pt = PKCS5Pad(pt, 32)

	b, err := aes.NewCipher(o.BlockKey)
	if err != nil {
		log.Fatal(getError(cipherBlockGeneration, err))
	}

	err, n := Bytes(12)
	if err != nil {
		log.Fatalf("Error generating nonce: %v", getError(bytesGeneration, err))
	}

	aesgcm, err := cipher.NewGCM(b)
	if err != nil {
		log.Fatal(getError(gcmCipherGeneration, err))
	}

	ct := aesgcm.Seal(nil, n, pt, ad)
	return n, ct
}

// DecryptGCM using AES-GCM-256 with RSA Signature of component of RSA private key as associated data
func (o Object) DecryptGCM(n, ct, ad []byte) []byte {
	if len(n) != nonceLength {
		if len(n) == 0 {
			log.Fatal(getError(nonceNil, nullError))
		} else {
			log.Fatal(getError(nonceWrongLength, nullError))
		}
	}

	b, err := aes.NewCipher(o.BlockKey)
	if err != nil {
		log.Fatal(getError(cipherBlockGeneration, err))
	}

	aesgcm, err := cipher.NewGCM(b)
	if err != nil {
		log.Fatal(getError(gcmCipherGeneration, err))
	}

	pt, err := aesgcm.Open(nil, n, ct, ad)
	if err != nil {
		log.Fatalf("Error decrypting: %v", err)
	}

	return PKCS5Trim(pt)
}
