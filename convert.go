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

func power2Round(parameters ParameterSet, r int) (int, int) {
	rPlus := modQ(r, parameters.Q)
	bound := 1 << parameters.D
	r0 := modCentered(rPlus, bound)

	return (rPlus - r0) / bound, r0
}

func decompose(parameters ParameterSet, r int) (int, int) {
	rPlus := modQ(r, parameters.Q)
	r0 := modCentered(rPlus, 2*parameters.Gamma2)
	r1 := 0

	if rPlus-r0 == parameters.Q-1 {
		r0 -= 1
	} else {
		r1 = (rPlus - r0) / (2 * parameters.Gamma2)
	}

	return r1, r0
}

func highBits(parameters ParameterSet, r int) int {
	r1, _ := decompose(parameters, r)
	return r1
}

func lowBits(parameters ParameterSet, r int) int {
	_, r0 := decompose(parameters, r)
	return r0
}

func makeHint(parameters ParameterSet, z, r int) bool {
	r1 := highBits(parameters, r)
	v1 := highBits(parameters, r+z)

	return r1 != v1
}

func useHint(parameters ParameterSet, h bool, r int) int {
	m := (parameters.Q - 1) / (2 * parameters.Gamma2)
	r1, r0 := decompose(parameters, r)

	if h {
		if r0 > 0 {
			return modQ(r1+1, m)
		} else {
			return modQ(r1-1, m)
		}
	}

	return r1
}

func pkEncode(parameters ParameterSet, rho []byte, t [][]int) []byte {
	width := bits.Len(uint(parameters.Q-1)) - parameters.D
	pk := make([]byte, len(rho))
	copy(pk, rho)

	for i := range parameters.K {
		// fmt.Printf("i: %d, t[i]: %v\n", i, t[i])
		// if i == 15 {
		// 	bytes, _ := hex.DecodeString("6d3bdd99d55d80e17d21163b61406ad4eaa70927e4fa74add922624d964725f11c9b7b52a5f9e3a6ec36e1f17d0ea61baf68ed8c04851a1a82730d39da1ad2e69e38288f55c13f75fc65dec5af6634ade84ee77459453a126f5a5902a806903c7914fbfb25515be9e57aebb8ca258d281e1a06109d85ea687de74a40f14235bd4d7541c05096800c47ad4d7f1554817c962d23840050c3f1c12966e586bcb6e71659168d96e6610ca391970581979aa40e6247b5c1661042468fa50e20e0435c7e7159b12fb3ec2d06dba6aa40030531f48071f645f7838d9faef5ed83ec5676cd4f5aa25e095cecceabc2df851488a5188ef9ef47b75ea42795d73b63800796331688fbf6e0c2fc0a6193c729209e013af51d52d1805b5ef72dda8e7827d38d92a70c4e09f6b0223dbc3e55c15ddb6aa5650d62078cfb6fe30668dd0c283ff3")
		// 	fmt.Printf("unpacked[1]: %v\n", simpleBitUnpack(bytes, (1<<width)-1))
		// 	pk = append(pk, bytes...)
		// } else {
		pk = append(pk, simpleBitPack(t[i], (1<<width)-1)...)
		// }
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
	// fmt.Printf("hints: %v\n", h)
	// fmt.Printf("packed hints: %v\n", hints)
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
		// fmt.Printf("%d\n", count)
		y[parameters.Omega+i] = byte(index)
	}

	// fmt.Printf("h: %v\n", h)
	// fmt.Printf("%d\n", len(y))
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

	// fmt.Printf("h: %v\n", h)
	return h
}
