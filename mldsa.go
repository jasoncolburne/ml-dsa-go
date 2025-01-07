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

func (dsa *MLDSA) KeyGenWithSeed(rnd []byte) (public []byte, private []byte, err error) {
	return keyGen(dsa.parameters, rnd)
}

// hedged signing
func (dsa *MLDSA) Sign(sk, message, ctx []byte) ([]byte, error) {
	if len(ctx) > 255 {
		return nil, fmt.Errorf("ctx length > 255")
	}

	rnd, err := rbg(SEEDLENGTH)
	if err != nil {
		return nil, err
	}

	mPrime := concatenateBytes(
		integerToBytes(0, 1),
		integerToBytes(int32(len(ctx)), 1),
		ctx,
		message,
	)

	sigma := sign(dsa.parameters, sk, mPrime, rnd)
	return sigma, nil
}

// deterministic signing
func (dsa *MLDSA) SignDeterministically(sk, message, ctx []byte) ([]byte, error) {
	if len(ctx) > 255 {
		return nil, fmt.Errorf("ctx length > 255")
	}

	rnd := make([]byte, SEEDLENGTH)

	mPrime := concatenateBytes(
		integerToBytes(0, 1),
		integerToBytes(int32(len(ctx)), 1),
		ctx,
		message,
	)

	sigma := sign(dsa.parameters, sk, mPrime, rnd)
	return sigma, nil
}

func (dsa *MLDSA) Verify(pk, message, signature, ctx []byte) (bool, error) {
	if len(ctx) > 255 {
		return false, fmt.Errorf("ctx length > 255")
	}

	mPrime := concatenateBytes(
		integerToBytes(0, 1),
		integerToBytes(int32(len(ctx)), 1),
		ctx,
		message,
	)

	return verify(dsa.parameters, pk, mPrime, signature), nil
}
