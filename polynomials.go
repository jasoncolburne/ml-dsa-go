package mldsa

import (
	"golang.org/x/crypto/sha3"
)

func sampleInBall(parameters ParameterSet, rho []byte) []int32 {
	c := make([]int32, 256)

	hasher := sha3.NewShake256()
	hasher.Write(rho)

	s := make([]byte, 8)
	hasher.Read(s)

	h := bytesToBits(s)
	for i := 256 - parameters.Tau; i < 256; i++ {
		jSlice := make([]byte, 1)
		hasher.Read(jSlice)

		for int32(jSlice[0]) > i {
			hasher.Read(jSlice)
		}

		j := int32(jSlice[0])
		c[i] = c[j]

		if h[i+parameters.Tau-256] {
			c[j] = -1
		} else {
			c[j] = 1
		}
	}

	return c
}

func rejNttPoly(parameters ParameterSet, rho []byte) []int32 {
	a := make([]int32, 256)

	hasher := sha3.NewShake128()
	hasher.Write(rho)

	j := 0
	for j < 256 {
		s := make([]byte, 3)
		hasher.Read(s)

		coefficient := coeffFromThreeBytes(parameters, s[0], s[1], s[2])

		// this pattern will prevent side channel attacks
		var newA int32
		delta := -1

		if coefficient != nil {
			newA = *coefficient
			delta = 1
		} else {
			newA = a[j]
			delta = 0
		}

		a[j] = newA
		j += delta
	}

	return a
}

func rejBoundedPoly(parameters ParameterSet, rho []byte) []int32 {
	a := make([]int32, 256)

	hasher := sha3.NewShake256()
	hasher.Write(rho)

	j := 0
	for j < 256 {
		zArray := make([]byte, 1)
		hasher.Read(zArray)

		z := int32(zArray[0])
		z0 := coeffFromHalfByte(parameters, modQ(z, 16))
		z1 := coeffFromHalfByte(parameters, z/16)

		// these patterns ensure no timing attack is introduced
		var newA int32
		delta := -1

		if z0 != nil {
			newA = *z0
			delta = 1
		} else {
			newA = a[j]
			delta = 0
		}

		a[j] = newA
		j += delta

		if j < 256 {
			if z1 != nil {
				newA = *z1
				delta = 1
			} else {
				newA = a[j]
				delta = 0
			}

			a[j] = newA
			j += delta
		}
	}

	return a
}

func addPolynomials(parameters ParameterSet, a, b []int32) []int32 {
	result := make([]int32, 256)

	for i := range 256 {
		result[i] = modQSymmetric(a[i]+b[i], parameters.Q)
	}

	return result
}

func subtractPolynomials(parameters ParameterSet, a, b []int32) []int32 {
	result := make([]int32, 256)

	for i := range 256 {
		result[i] = modQSymmetric(a[i]-b[i], parameters.Q)
	}

	return result
}

func vectorAddPolynomials(parameters ParameterSet, a, b [][]int32) [][]int32 {
	length := len(a)

	result := make([][]int32, length)

	for i := range length {
		result[i] = addPolynomials(parameters, a[i], b[i])
	}

	return result
}

func vectorSubtractPolynomials(parameters ParameterSet, a, b [][]int32) [][]int32 {
	length := len(a)

	result := make([][]int32, length)

	for i := range length {
		result[i] = subtractPolynomials(parameters, a[i], b[i])
	}

	return result
}

func scalarVectorMultiply(parameters ParameterSet, c int32, v [][]int32) [][]int32 {
	w := make([][]int32, len(v))

	for i, row := range v {
		w[i] = make([]int32, len(row))
		for j, value := range row {
			w[i][j] = modQSymmetric(value*c, parameters.Q)
		}
	}

	return w
}
