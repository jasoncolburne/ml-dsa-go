package mldsa

import (
	"math/rand"
	"testing"
)

func benchmarkNtt(parameters ParameterSet, b *testing.B) {
	w := make([]int32, 256)
	for i := range 256 {
		w[i] = int32(rand.Intn(int(parameters.Q)))
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Function under test
		ntt(parameters, w)
	}
}

func BenchmarkNTT(b *testing.B) {
	benchmarkNtt(ML_DSA_44_Parameters, b)
}
