package mldsa_test

import (
	"encoding/hex"
	"math/rand"
	"testing"

	mldsa "github.com/jasoncolburne/ml-dsa-go"
)

func testParamsRoundtrip(params mldsa.ParameterSet, t *testing.T, skLength, vkLength, sigLength int) {
	dsa := mldsa.Init(params)

	vk, sk, err := dsa.KeyGen()
	if err != nil {
		t.Fatalf("error generating keypair: %v", err)
	}

	if len(sk) != skLength {
		t.Fatalf("unexpected sk length: %d != %d", len(sk), skLength)
	}

	if len(vk) != vkLength {
		t.Fatalf("unexpected sk length: %d != %d", len(vk), vkLength)
	}

	message, _ := hex.DecodeString("89b0c4b23019af3498a27da290892d981dd59fa08993bc05da21e1d72503664c98cadefc061d176d0b44bcab049bb540e0680a58bdad0d16316f772d44d47281")
	ctx, _ := hex.DecodeString("09764e76473cc969442691dd0574afdd")
	sig, err := dsa.Sign(sk, message, ctx)
	if err != nil {
		t.Fatalf("error signing: %v", err)
	}

	if len(sig) != sigLength {
		t.Fatalf("unexpected sig length: %d != %d", len(sig), sigLength)
	}

	valid, err := dsa.Verify(vk, message, sig, ctx)
	if err != nil {
		t.Fatalf("error verifying: %v", err)
	}

	if !valid {
		t.Fatalf("signature not valid!")
	}

	sigPrime := copyAndMutate(sig)
	valid, err = dsa.Verify(vk, message, sigPrime, ctx)
	if err != nil {
		t.Fatalf("error verifying: %v", err)
	}

	if valid {
		t.Fatalf("mutated signature is still valid!")
	}

	messagePrime := copyAndMutate(message)
	valid, err = dsa.Verify(vk, messagePrime, sig, ctx)
	if err != nil {
		t.Fatalf("error verifying: %v", err)
	}

	if valid {
		t.Fatalf("mutated message is still valid!")
	}

	ctxPrime := copyAndMutate(ctx)
	valid, err = dsa.Verify(vk, message, sig, ctxPrime)
	if err != nil {
		t.Fatalf("error verifying: %v", err)
	}

	if valid {
		t.Fatalf("mutated context is still valid!")
	}
}

func copyAndMutate(bytes []byte) []byte {
	result := make([]byte, len(bytes))
	copy(result, bytes)

	offset := rand.Intn(len(bytes))
	// flip a bit in some byte
	result[offset] = byte(int32(result[offset]) ^ 0x01)

	return result
}

func TestMLDSA44RoundTrip(t *testing.T) {
	testParamsRoundtrip(mldsa.ML_DSA_44_Parameters, t, 2560, 1312, 2420)
}

func TestMLDSA65RoundTrip(t *testing.T) {
	testParamsRoundtrip(mldsa.ML_DSA_65_Parameters, t, 4032, 1952, 3309)
}

func TestMLDSA87RoundTrip(t *testing.T) {
	testParamsRoundtrip(mldsa.ML_DSA_87_Parameters, t, 4896, 2592, 4627)
}
