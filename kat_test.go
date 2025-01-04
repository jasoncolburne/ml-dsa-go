package mldsa_test

import (
	"crypto/subtle"
	"encoding/hex"
	"testing"

	mldsa "github.com/jasoncolburne/ml-dsa-go"
)

type TestVector struct {
	Count      int32
	Seed       string
	PrivateKey string
	PublicKey  string
	Message    string
	Signature  string
	Context    string
}

func testKatVectors(vectors []TestVector, parameters mldsa.ParameterSet, t *testing.T) {
	dsa := mldsa.Init(parameters)
	for _, vector := range vectors {
		seed, _ := hex.DecodeString(vector.Seed)
		pk, sk, err := dsa.KeyGenWithSeed(seed)
		if err != nil {
			t.Fatalf("err :%v", err)
		}

		expectedPk, _ := hex.DecodeString(vector.PublicKey)
		if subtle.ConstantTimeCompare(pk, expectedPk) != 1 {
			t.Fatalf("bad pk")
		}

		expectedSk, _ := hex.DecodeString(vector.PrivateKey)
		if subtle.ConstantTimeCompare(sk, expectedSk) != 1 {
			t.Fatalf("bad sk")
		}

		message, _ := hex.DecodeString(vector.Message)
		ctx, _ := hex.DecodeString(vector.Context)
		sig, err := dsa.SignDeterministically(sk, message, ctx)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		sm := make([]byte, len(sig))
		copy(sm, sig)
		sm = append(sm, message...)
		expectedSm, _ := hex.DecodeString(vector.Signature)

		if subtle.ConstantTimeCompare(sm, expectedSm) != 1 {
			t.Fatalf("bad sm")
		}
	}
}

func TestMLDSA44Kat(t *testing.T) {
	testKatVectors(ML_DSA_44_TestVectors, mldsa.ML_DSA_44_Parameters, t)
}

func TestMLDSA65Kat(t *testing.T) {
	testKatVectors(ML_DSA_65_TestVectors, mldsa.ML_DSA_65_Parameters, t)
}

func TestMLDSA87Kat(t *testing.T) {
	testKatVectors(ML_DSA_87_TestVectors, mldsa.ML_DSA_87_Parameters, t)
}
