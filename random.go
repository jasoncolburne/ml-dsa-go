package mldsa

import (
	"crypto"

	drbg "github.com/canonical/go-sp800.90a-drbg"
)

func rbg(len int) ([]byte, error) {
	hashRbg, err := drbg.NewHash(crypto.SHA256, []byte{}, nil)
	if err != nil {
		return nil, err
	}

	bytes := make([]byte, len)
	err = hashRbg.Generate([]byte{}, bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
