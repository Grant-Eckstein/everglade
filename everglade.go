package everglade

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"log"
)

type Object struct {
	BlockKey []byte
	RSAKey   *rsa.PrivateKey
}

const (
	nonceLength = 12
)

func New() Object {
	err, bk := Bytes(32)
	if err != nil {
		log.Fatalf("%v", getError(bytesGeneration, err))
	}

	rk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("%v", getError(rsaKeyGeneration, err))
	}

	return Object{
		BlockKey: bk,
		RSAKey:   rk,
	}
}

func (o Object) Export() []byte {
	b, err := json.Marshal(o)
	if err != nil {
		log.Fatalf("%v", getError(jsonMarshal, err))
	}

	return b
}

func Import(j []byte) Object {
	var o Object
	err := json.Unmarshal(j, &o)

	if err != nil {
		log.Fatalf("%v", getError(jsonUnmarshal, err))
	}

	return o
}
