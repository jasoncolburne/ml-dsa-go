package mldsa

import (
	"fmt"
)

const SEEDLENGTH = 32

type MLDSA struct {
	parameters ParameterSet
}

func Init(parameters ParameterSet) *MLDSA {
	return &MLDSA{parameters: parameters}
}

func (dsa *MLDSA) KeyGen() (public []byte, private []byte, err error) {
	rnd, err := rbg(SEEDLENGTH)
	if err != nil {
		return nil, nil, err
	}

	return keyGen(dsa.parameters, rnd)
}

func (dsa *MLDSA) Sign(sk, message, ctx []byte) ([]byte, error) {
	if len(ctx) > 255 {
		return nil, fmt.Errorf("ctx length > 255")
	}

	rnd, err := rbg(SEEDLENGTH)
	if err != nil {
		return nil, err
	}

	// for KAT testing
	// rnd = make([]byte, SEEDLENGTH)

	mPrime := integerToBytes(0, 1)
	mPrime = append(mPrime, integerToBytes(len(ctx), 1)...)
	mPrime = append(mPrime, ctx...)
	mPrime = append(mPrime, message...)

	sigma := sign(dsa.parameters, sk, mPrime, rnd)
	return sigma, nil
}

func (dsa *MLDSA) Verify(pk, message, signature, ctx []byte) (bool, error) {
	if len(ctx) > 255 {
		return false, fmt.Errorf("ctx length > 255")
	}

	mPrime := integerToBytes(0, 1)
	mPrime = append(mPrime, integerToBytes(len(ctx), 1)...)
	mPrime = append(mPrime, ctx...)
	mPrime = append(mPrime, message...)

	return verify(dsa.parameters, pk, mPrime, signature), nil
}
