package mldsa

import (
	"math/rand"
	"testing"
)

func TestNTTRoundTrip(t *testing.T) {
	params := ML_DSA_44_Parameters
	input := make([]int, 256)
	for i := range 256 {
		input[i] = rand.Intn(params.Q) - params.Q/2
	}
	result := nttInverse(params, ntt(params, input))

	for i := range 256 {
		if result[i] != input[i] {
			t.Fatalf("mismatch at %d: %d != %d", i, result[i], input[i])
		}
	}
}
