package mldsa

import (
	"math/bits"
)

func expandA(parameters ParameterSet, rho []byte) [][][]int32 {
	A := make([][][]int32, parameters.K)

	rhoLength := len(rho)
	rhoPrime := make([]byte, rhoLength+2)
	copy(rhoPrime, rho)

	for r := range parameters.K {
		A[r] = make([][]int32, parameters.L)

		for s := range parameters.L {
			rhoPrime[rhoLength] = integerToBytes(s, 1)[0]
			rhoPrime[rhoLength+1] = integerToBytes(r, 1)[0]

			A[r][s] = rejNttPoly(parameters, rhoPrime)
		}
	}

	return A
}

func expandS(parameters ParameterSet, rho []byte) ([][]int32, [][]int32) {
	rhoLength := len(rho)
	rhoPrime := make([]byte, rhoLength+2)
	copy(rhoPrime, rho)

	s1 := make([][]int32, parameters.L)
	s2 := make([][]int32, parameters.K)

	for r := range parameters.L {
		copy(rhoPrime[rhoLength:], integerToBytes(r, 2))
		s1[r] = rejBoundedPoly(parameters, rhoPrime)
	}

	for r := range parameters.K {
		copy(rhoPrime[rhoLength:], integerToBytes(r+parameters.L, 2))
		s2[r] = rejBoundedPoly(parameters, rhoPrime)
	}

	return s1, s2
}

func expandMask(parameters ParameterSet, rho []byte, mu int32) [][]int32 {
	c := 1 + bits.Len(uint(parameters.Gamma1-1))

	rhoPrime := make([]byte, 66)
	copy(rhoPrime[:64], rho)

	y := make([][]int32, parameters.L)
	for r := range parameters.L {
		copy(rhoPrime[64:], integerToBytes(mu+r, 2))
		v := concatenateBytesAndSHAKE(true, int32(32*c), rhoPrime)
		y[r] = bitUnpack(v, parameters.Gamma1-1, parameters.Gamma1)
	}

	return y
}
