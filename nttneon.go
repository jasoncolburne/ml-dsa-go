//go:build neon

package mldsa

/*
#cgo CFLAGS: -march=armv8-a+simd -O3
#cgo LDFLAGS: -lm
#include "nttneon.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

func ntt(parameters ParameterSet, w []int32) []int32 {
	wHat := make([]int32, 256)
	copy(wHat, w)

	m := 0
	q := parameters.Q

	for len := 128; len >= 1; len /= 2 {
		for start := 0; start < 256; start += 2 * len {
			m += 1
			z := zetas[m]

			for j := start; j < start+len; j++ {
				t := modMultiply(z, wHat[j+len], q)
				wHat[j+len] = modQ(wHat[j]-t, q)
				wHat[j] = modQ(wHat[j]+t, q)
			}
		}
	}

	return wHat
}

func nttInverse(parameters ParameterSet, wHat []int32) []int32 {
	w := make([]int32, 256)

	for j := range 256 {
		w[j] = wHat[j]
	}

	m := 256
	len := 1
	q := parameters.Q

	for len < 256 {
		start := 0
		for start < 256 {
			m -= 1
			z := -zetas[m]
			for j := start; j < start+len; j++ {
				t := w[j]
				w[j] = modQ(t+w[j+len], q)
				w[j+len] = modQ(t-w[j+len], q)
				w[j+len] = modMultiply(z, w[j+len], q)
			}
			start += 2 * len
		}
		len = 2 * len
	}

	f := int32(8347681)
	// modMultiply(f, 256, q) == 1

	for j := range 256 {
		w[j] = modCentered(modMultiply(f, w[j], q), q)
	}

	return w
}

func addNtt(parameters ParameterSet, aHat, bHat []int32) []int32 {
	if len(aHat) != 256 || len(bHat) != 256 {
		panic("input arrays must have length 256")
	}

	cHat := make([]int32, 256)

	C.add_ntt(
		(*C.int)(unsafe.Pointer(&aHat[0])),
		(*C.int)(unsafe.Pointer(&bHat[0])),
		(*C.int)(unsafe.Pointer(&cHat[0])),
		C.int(parameters.Q),
	)

	return cHat
}

func subtractNtt(parameters ParameterSet, aHat, bHat []int32) []int32 {
	if len(aHat) != 256 || len(bHat) != 256 {
		panic("input arrays must have length 256")
	}

	cHat := make([]int32, 256)

	C.subtract_ntt(
		(*C.int)(unsafe.Pointer(&aHat[0])),
		(*C.int)(unsafe.Pointer(&bHat[0])),
		(*C.int)(unsafe.Pointer(&cHat[0])),
		C.int(parameters.Q),
	)

	return cHat
}

func multiplyNtt(parameters ParameterSet, aHat, bHat []int32) []int32 {
	cHat := make([]int32, 256)

	for i := range 256 {
		cHat[i] = modMultiply(aHat[i], bHat[i], parameters.Q)
	}

	return cHat
}

func vectorNtt(parameters ParameterSet, v [][]int32) [][]int32 {
	length := len(v)

	vHat := make([][]int32, length)
	for j := range length {
		vHat[j] = ntt(parameters, v[j])
	}

	return vHat
}

func subtractVectorNtt(parameters ParameterSet, vHat, wHat [][]int32) [][]int32 {
	length := len(vHat)
	uHat := make([][]int32, length)

	for i := range length {
		uHat[i] = subtractNtt(parameters, vHat[i], wHat[i])
	}

	return uHat
}

func vectorNttInverse(parameters ParameterSet, vHat [][]int32) [][]int32 {
	length := len(vHat)

	v := make([][]int32, length)
	for j := range length {
		v[j] = nttInverse(parameters, vHat[j])
	}

	return v
}

func scalarVectorNtt(parameters ParameterSet, cHat []int32, vHat [][]int32) [][]int32 {
	length := len(vHat)
	wHat := make([][]int32, length)

	for i := range length {
		wHat[i] = multiplyNtt(parameters, cHat, vHat[i])
	}

	return wHat
}

// func matrixVectorNtt(parameters ParameterSet, MHat [][][]int32, vHat [][]int32) [][]int32 {
// 	K := parameters.K
// 	L := parameters.L
// 	Q := parameters.Q

// 	wHat := make([][]int32, K)
// 	for i := range K {
// 		wHat[i] = make([]int32, 256)
// 	}

// 	// Ensure correct sizes
// 	if len(MHat) != K || len(wHat) != K || len(vHat) != L {
// 		panic("Invalid dimensions for MHat, vHat, or wHat")
// 	}

// 	// Flatten MHat
// 	MHatFlat := make([]C.int32, K*L*256)
// 	for i := 0; i < K; i++ {
// 		for j := 0; j < L; j++ {
// 			for k := 0; k < 256; k++ {
// 				MHatFlat[i*L*256+j*256+k] = C.int32(MHat[i][j][k])
// 			}
// 		}
// 	}

// 	fmt.Printf("mhatflat %v\n", MHatFlat)

// 	// Flatten vHat
// 	vHatFlat := make([]C.int32, L*256)
// 	for i := 0; i < L; i++ {
// 		for j := 0; j < 256; j++ {
// 			vHatFlat[i*256+j] = C.int32(vHat[i][j])
// 		}
// 	}

// 	// Allocate flat wHat for results
// 	wHatFlat := make([]C.int32, K*256)

// 	// Call the C function
// 	C.matrix_vector_ntt(
// 		(***C.int32)(unsafe.Pointer(&MHatFlat[0])),
// 		(**C.int32)(unsafe.Pointer(&vHatFlat[0])),
// 		(**C.int32)(unsafe.Pointer(&wHatFlat[0])),
// 		C.int32(K),
// 		C.int32(L),
// 		C.int32(Q),
// 	)

// 	// Copy results back into wHat
// 	for i := 0; i < K; i++ {
// 		for j := 0; j < 256; j++ {
// 			wHat[i][j] = int32(wHatFlat[i*256+j])
// 		}
// 	}

// 	// fmt.Printf("wHat: %v\n", wHat)

// 	return wHat
// }

func matrixVectorNtt(parameters ParameterSet, MHat [][][]int32, vHat [][]int32) [][]int32 {
	wHat := make([][]int32, parameters.K)

	for i := range parameters.K {
		wHat[i] = make([]int32, 256)
		for j := range parameters.L {
			wHat[i] = addNtt(parameters, wHat[i], multiplyNtt(parameters, MHat[i][j], vHat[j]))
		}
	}

	// fmt.Printf("wHat: %v\n", wHat)

	return wHat
}
