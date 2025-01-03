package mldsa

import (
	"golang.org/x/crypto/sha3"
)

func sampleInBall(parameters ParameterSet, rho []byte) []int {
	c := make([]int, 256)

	hasher := sha3.NewShake256()
	hasher.Write(rho)

	s := make([]byte, 8)
	hasher.Read(s)

	h := bytesToBits(s)
	for i := 256 - parameters.Tau; i < 256; i++ {
		jSlice := make([]byte, 1)
		hasher.Read(jSlice)

		for int(jSlice[0]) > i {
			hasher.Read(jSlice)
		}

		j := int(jSlice[0])
		c[i] = c[j]

		if h[i+parameters.Tau-256] {
			c[j] = -1
		} else {
			c[j] = 1
		}
	}

	return c
}

func rejNttPoly(parameters ParameterSet, rho []byte) []int {
	a := make([]int, 256)

	hasher := sha3.NewShake128()
	hasher.Write(rho)

	j := 0
	for j < 256 {
		s := make([]byte, 3)
		hasher.Read(s)

		coefficient := coeffFromThreeBytes(parameters, s[0], s[1], s[2])
		if coefficient == nil {
			continue
		}

		a[j] = *coefficient
		j += 1
	}

	return a
}

func rejBoundedPoly(parameters ParameterSet, rho []byte) []int {
	a := make([]int, 256)

	hasher := sha3.NewShake256()
	hasher.Write(rho)

	j := 0
	for j < 256 {
		zArray := make([]byte, 1)
		hasher.Read(zArray)

		z := int(zArray[0])
		z0 := coeffFromHalfByte(parameters, modQ(z, 16))
		z1 := coeffFromHalfByte(parameters, z/16)

		if z0 != nil {
			a[j] = *z0
			j += 1
		}

		if z1 != nil && j < 256 {
			a[j] = *z1
			j += 1
		}
	}

	return a
}

var q = ML_DSA_44_Parameters.Q

func addPolynomials(a, b []int) []int {
	result := make([]int, 256)

	for i := range 256 {
		result[i] = modCentered(a[i]+b[i], q)
	}

	return result
}

func subtractPolynomials(a, b []int) []int {
	result := make([]int, 256)

	for i := range 256 {
		result[i] = modCentered(a[i]-b[i], q)
	}

	return result
}

func vectorAddPolynomials(a, b [][]int) [][]int {
	length := len(a)

	result := make([][]int, length)

	for i := range length {
		result[i] = addPolynomials(a[i], b[i])
	}

	return result
}

func vectorSubtractPolynomials(a, b [][]int) [][]int {
	length := len(a)

	result := make([][]int, length)

	for i := range length {
		result[i] = subtractPolynomials(a[i], b[i])
	}

	return result
}

func scalarVectorMultiply(c int, v [][]int) [][]int {
	w := make([][]int, len(v))

	for i, row := range v {
		w[i] = make([]int, len(row))
		for j, value := range row {
			w[i][j] = modCentered(value*c, q)
		}
	}

	return w
}
