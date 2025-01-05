//go:build !neon

package mldsa

func addNtt(parameters ParameterSet, aHat, bHat []int32) []int32 {
	cHat := make([]int32, 256)

	for i := range 256 {
		cHat[i] = modQ(aHat[i]+bHat[i], parameters.Q)
	}

	return cHat
}

func subtractNtt(parameters ParameterSet, aHat, bHat []int32) []int32 {
	cHat := make([]int32, 256)

	for i := range 256 {
		cHat[i] = modQ(aHat[i]-bHat[i], parameters.Q)
	}

	return cHat
}

func multiplyNtt(parameters ParameterSet, aHat, bHat []int32) []int32 {
	cHat := make([]int32, 256)

	for i := range 256 {
		cHat[i] = modMultiply(aHat[i], bHat[i], parameters.Q)
	}

	return cHat
}
