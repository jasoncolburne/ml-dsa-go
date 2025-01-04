package mldsa

import (
	"math/rand"
	"testing"
)

func TestNTTRoundTrip(t *testing.T) {
	params := ML_DSA_44_Parameters
	input := make([]int32, 256)
	for i := range 256 {
		input[i] = int32(rand.Intn(int(params.Q))) - params.Q/2
	}
	result := nttInverse(params, ntt(params, input))

	for i := range 256 {
		if result[i] != input[i] {
			// fmt.Printf("%v", result)
			t.Fatalf("mismatch at %d: %d != %d", i, result[i], input[i])
		}
	}
}
