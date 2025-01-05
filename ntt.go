package mldsa

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
		w[j] = modQSymmetric(modMultiply(f, w[j], q), q)
	}

	return w
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

func matrixVectorNtt(parameters ParameterSet, MHat [][][]int32, vHat [][]int32) [][]int32 {
	wHat := make([][]int32, parameters.K)

	for i := range parameters.K {
		wHat[i] = make([]int32, 256)
		for j := range parameters.L {
			wHat[i] = addNtt(parameters, wHat[i], multiplyNtt(parameters, MHat[i][j], vHat[j]))
		}
	}

	return wHat
}
