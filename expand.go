package mldsa

import (
	"math/bits"

	"golang.org/x/crypto/sha3"
)

func expandA(parameters ParameterSet, rho []byte) [][][]int {
	A := make([][][]int, parameters.K)

	rhoLength := len(rho)
	rhoPrime := make([]byte, rhoLength+2)
	copy(rhoPrime, rho)

	for r := range parameters.K {
		A[r] = make([][]int, parameters.L)

		for s := range parameters.L {
			rhoPrime[rhoLength] = integerToBytes(s, 1)[0]
			rhoPrime[rhoLength+1] = integerToBytes(r, 1)[0]
			// fmt.Printf("rhoPrime: %v\n", rhoPrime)

			A[r][s] = rejNttPoly(parameters, rhoPrime)
		}
	}

	return A
}

func expandS(parameters ParameterSet, rho []byte) ([][]int, [][]int) {
	rhoLength := len(rho)
	rhoPrime := make([]byte, rhoLength+2)
	copy(rhoPrime, rho)

	s1 := make([][]int, parameters.L)
	s2 := make([][]int, parameters.K)

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

func expandMask(parameters ParameterSet, rho []byte, mu int) [][]int {
	c := 1 + bits.Len(uint(parameters.Gamma1-1))

	rhoPrime := make([]byte, 66)
	copy(rhoPrime[:64], rho)

	y := make([][]int, parameters.L)
	for r := range parameters.L {
		copy(rhoPrime[64:], integerToBytes(mu+r, 2))

		v := make([]byte, 32*c)
		sha3.ShakeSum256(v, rhoPrime)

		y[r] = bitUnpack(v, parameters.Gamma1-1, parameters.Gamma1)
	}

	return y
}
