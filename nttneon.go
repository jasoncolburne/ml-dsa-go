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
	if len(aHat) != 256 || len(bHat) != 256 {
		panic("input arrays must have length 256")
	}

	a := make([]int64, 256)
	b := make([]int64, 256)

	for i := range 256 {
		a[i] = int64(aHat[i])
		b[i] = int64(bHat[i])
	}

	c := make([]int64, 256)

	C.multiply_ntt(
		(*C.int)(unsafe.Pointer(&a[0])),
		(*C.int)(unsafe.Pointer(&b[0])),
		(*C.int)(unsafe.Pointer(&c[0])),
		C.int(parameters.Q),
	)

	cHat := make([]int32, 256)
	for i := range 256 {
		cHat[i] = int32(c[i])
	}

	return cHat
}
