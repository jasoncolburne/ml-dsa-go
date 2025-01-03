package mldsa

import (
	"math/bits"
)

func coeffFromHalfByte(parameters ParameterSet, b int) *int {
	if parameters.Eta == 2 && b < 15 {
		result := 2 - modQ(b, 5)
		return &result
	}

	if parameters.Eta == 4 && b < 9 {
		result := 4 - b
		return &result
	}

	return nil
}

func coeffFromThreeBytes(parameters ParameterSet, b0, b1, b2 byte) *int {
	b2Prime := int(b2)
	if b2Prime > 127 {
		b2Prime -= 128
	}

	z := 65536*b2Prime + 256*int(b1) + int(b0)
	if z < parameters.Q {
		return &z
	}

	return nil
}

func bitsToBytes(y []bool) []byte {
	alpha := len(y)
	z := make([]byte, (alpha+7)/8)

	for i := range alpha {
		// TODO: evaluate attacks (we optimized out a computation using this bool)
		if y[i] {
			z[i/8] += (1 << modQ(i, 8))
		}
	}

	return z
}

func bytesToBits(z []byte) []bool {
	zLength := len(z)
	zPrime := make([]uint8, zLength)
	for i := range zLength {
		zPrime[i] = uint8(z[i])
	}

	y := make([]bool, 8*zLength)

	for i := range zLength {
		for j := range 8 {
			y[8*i+j] = modQ(int(zPrime[i]), 2) == 1
			zPrime[i] /= 2
		}
	}

	return y
}

func bitsToInteger(y []bool, alpha int) int {
	x := 0

	for i := 1; i <= alpha; i++ {
		x <<= 1
		if y[alpha-i] {
			x += 1
		}
	}

	return x
}

func integerToBits(x, alpha int) []bool {
	y := make([]bool, alpha)

	xPrime := x

	for i := range alpha {
		y[i] = modQ(xPrime, 2) == 1
		xPrime /= 2
	}

	return y
}

func integerToBytes(x int, alpha int) []byte {
	y := make([]byte, alpha)

	xPrime := x
	for i := range alpha {
		y[i] = byte(modQ(xPrime, 256))
		xPrime /= 256
	}

	return y
}

func pkEncode(parameters ParameterSet, rho []byte, t [][]int) []byte {
	width := bits.Len(uint(parameters.Q-1)) - parameters.D
	pk := make([]byte, len(rho))
	copy(pk, rho)

	for i := range parameters.K {
		pk = append(pk, simpleBitPack(t[i], (1<<width)-1)...)
	}

	return pk
}

func pkDecode(parameters ParameterSet, pk []byte) ([]byte, [][]int) {
	rho := pk[:32]
	z := pk[32:]
	toShift := bits.Len(uint(parameters.Q-1)) - parameters.D
	width := 32 * toShift

	t := make([][]int, parameters.K)
	for i := range parameters.K {
		offset := i * width
		limit := offset + width
		t[i] = simpleBitUnpack(z[offset:limit], (1<<toShift)-1)
	}

	return rho, t
}

func skEncode(
	parameters ParameterSet,
	rho, kappa, tr []byte,
	s1, s2, t [][]int,
) []byte {
	sk := make([]byte, 128)
	copy(sk[:32], rho)
	copy(sk[32:64], kappa)
	copy(sk[64:], tr)

	eta := parameters.Eta

	for i := range parameters.L {
		sk = append(sk, bitPack(s1[i], eta, eta)...)
	}

	for i := range parameters.K {
		sk = append(sk, bitPack(s2[i], eta, eta)...)
	}

	x := 1 << (parameters.D - 1)
	y := x - 1
	for i := range parameters.K {
		sk = append(sk, bitPack(t[i], y, x)...)
	}

	return sk
}

// this function uses named returns, brace yourself
func skDecode(parameters ParameterSet, sk []byte) (
	rho []byte,
	kappa []byte,
	tr []byte,
	s1 [][]int,
	s2 [][]int,
	t [][]int,
) {
	rho = sk[:32]
	kappa = sk[32:64]
	tr = sk[64:128]

	baseOffset := 128

	eta := parameters.Eta
	width := 32 * bits.Len(uint(2*eta))

	s1 = make([][]int, parameters.L)
	for i := range parameters.L {
		offset := baseOffset + width*i
		limit := offset + width
		y := sk[offset:limit]
		s1[i] = bitUnpack(y, eta, eta)
	}

	baseOffset += width * parameters.L

	s2 = make([][]int, parameters.K)
	for i := range parameters.K {
		offset := baseOffset + width*i
		limit := offset + width
		z := sk[offset:limit]
		s2[i] = bitUnpack(z, eta, eta)
	}

	baseOffset += width * parameters.K
	wWidth := 32 * parameters.D
	x := 1 << (parameters.D - 1)
	y := x - 1

	t = make([][]int, parameters.K)
	for i := range parameters.K {
		offset := baseOffset + wWidth*i
		limit := offset + wWidth
		w := sk[offset:limit]
		t[i] = bitUnpack(w, y, x)
	}

	return
}

func sigEncode(parameters ParameterSet, cTilde []byte, z [][]int, h [][]bool) []byte {
	sigma := make([]byte, len(cTilde))
	copy(sigma, cTilde)

	gamma1 := parameters.Gamma1
	for i := range parameters.L {
		sigma = append(sigma, bitPack(z[i], gamma1-1, gamma1)...)
	}
	hints := hintBitPack(parameters, h)
	sigma = append(sigma, hints...)

	return sigma
}

func sigDecode(parameters ParameterSet, sigma []byte) ([]byte, [][]int, [][]bool) {
	width := 32 * (1 + bits.Len(uint(parameters.Gamma1-1)))

	cTilde := sigma[:parameters.Lambda/4]
	x := sigma[parameters.Lambda/4 : parameters.Lambda/4+parameters.L*width]
	y := sigma[parameters.Lambda/4+parameters.L*width:]

	z := make([][]int, parameters.L)
	for i := range parameters.L {
		offset := i * width
		limit := offset + width
		z[i] = bitUnpack(x[offset:limit], parameters.Gamma1-1, parameters.Gamma1)
	}

	h := hintBitUnpack(parameters, y)

	return cTilde, z, h
}

func w1Encode(parameters ParameterSet, w [][]int) []byte {
	w1Tilde := []byte{}
	length := len(w)

	for i := range length {
		b := (parameters.Q-1)/(2*parameters.Gamma2) - 1
		w1Tilde = append(w1Tilde, simpleBitPack(w[i], b)...)
	}

	return w1Tilde
}

func bitPack(w []int, a, b int) []byte {
	z := []bool{}

	for i := range 256 {
		z = append(z, integerToBits(b-w[i], bits.Len(uint(a+b)))...)
	}

	return bitsToBytes(z)
}

func bitUnpack(v []byte, a, b int) []int {
	c := bits.Len(uint(a + b))
	z := bytesToBits(v)

	w := make([]int, 256)

	for i := range 256 {
		offset := i * c
		limit := offset + c
		w[i] = b - bitsToInteger(z[offset:limit], c)
	}

	return w
}

func simpleBitPack(w []int, b int) []byte {
	z := []bool{}
	for i := range 256 {
		z = append(z, integerToBits(w[i], bits.Len(uint(b)))...)
	}

	return bitsToBytes(z)
}

func simpleBitUnpack(v []byte, b int) []int {
	c := bits.Len(uint(b))
	z := bytesToBits(v)

	w := make([]int, 256)
	for i := range 256 {
		offset := i * c
		limit := offset + c
		w[i] = bitsToInteger(z[offset:limit], c)
	}

	return w
}

func hintBitPack(parameters ParameterSet, h [][]bool) []byte {
	y := make([]byte, parameters.Omega+parameters.K)
	index := 0

	for i := range parameters.K {
		count := 0

		for j := range 256 {
			if h[i][j] {
				count += 1
				y[index] = byte(j)
				index += 1
			}
		}

		y[parameters.Omega+i] = byte(index)
	}

	return y
}

func hintBitUnpack(parameters ParameterSet, y []byte) [][]bool {
	h := make([][]bool, parameters.K)
	index := 0
	omega := parameters.Omega

	for i := range parameters.K {
		h[i] = make([]bool, 256)

		yOmegaI := int(y[omega+i])

		if yOmegaI < index || yOmegaI > omega {
			return nil
		}

		first := index
		for index < yOmegaI {
			if index > first {
				if y[index-1] >= y[index] {
					return nil
				}
			}

			h[i][int(y[index])] = true
			index += 1
		}
	}

	for i := index; i < omega; i++ {
		if y[i] != byte(0) {
			return nil
		}
	}

	return h
}
