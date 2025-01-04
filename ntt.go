//go:build !hardware

package mldsa

func ntt(parameters ParameterSet, w []int) []int {
	wHat := make([]int, 256)
	for j := range 256 {
		wHat[j] = w[j]
	}

	m := 0
	len := 128
	q := parameters.Q

	for len >= 1 {
		start := 0
		for start < 256 {
			m += 1
			z := zetas[m]

			for j := start; j < start+len; j++ {
				t := modQ(z*wHat[j+len], q)
				wHat[j+len] = modQ(wHat[j]-t, q)
				wHat[j] = modQ(wHat[j]+t, q)
			}

			start += 2 * len
		}

		len /= 2
	}

	return wHat
}

func nttInverse(parameters ParameterSet, wHat []int) []int {
	w := make([]int, 256)

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
				w[j+len] = modQ(z*w[j+len], q)
			}
			start += 2 * len
		}
		len = 2 * len
	}

	f := 8347681

	for j := range 256 {
		w[j] = modCentered(f*w[j], q)
	}

	return w
}

func addNtt(parameters ParameterSet, aHat, bHat []int) []int {
	cHat := make([]int, 256)

	for i := range 256 {
		cHat[i] = modQ(aHat[i]+bHat[i], parameters.Q)
	}

	return cHat
}

func subtractNtt(parameters ParameterSet, aHat, bHat []int) []int {
	cHat := make([]int, 256)

	for i := range 256 {
		cHat[i] = modQ(aHat[i]-bHat[i], parameters.Q)
	}

	return cHat
}

func multiplyNtt(parameters ParameterSet, aHat, bHat []int) []int {
	cHat := make([]int, 256)

	for i := range 256 {
		cHat[i] = modQ(aHat[i]*bHat[i], parameters.Q)
	}

	return cHat
}

func vectorNtt(parameters ParameterSet, v [][]int) [][]int {
	length := len(v)

	vHat := make([][]int, length)
	for j := range length {
		vHat[j] = ntt(parameters, v[j])
	}

	return vHat
}

func subtractVectorNtt(parameters ParameterSet, vHat, wHat [][]int) [][]int {
	length := len(vHat)
	uHat := make([][]int, length)

	for i := range length {
		uHat[i] = subtractNtt(parameters, vHat[i], wHat[i])
	}

	return uHat
}

func vectorNttInverse(parameters ParameterSet, vHat [][]int) [][]int {
	length := len(vHat)

	v := make([][]int, length)
	for j := range length {
		v[j] = nttInverse(parameters, vHat[j])
	}

	return v
}

func scalarVectorNtt(parameters ParameterSet, cHat []int, vHat [][]int) [][]int {
	length := len(vHat)
	wHat := make([][]int, length)

	for i := range length {
		wHat[i] = multiplyNtt(parameters, cHat, vHat[i])
	}

	return wHat
}

func matrixVectorNtt(parameters ParameterSet, MHat [][][]int, vHat [][]int) [][]int {
	wHat := make([][]int, parameters.K)

	for i := range parameters.K {
		wHat[i] = make([]int, 256)
		for j := range parameters.L {
			wHat[i] = addNtt(parameters, wHat[i], multiplyNtt(parameters, MHat[i][j], vHat[j]))
		}
	}

	return wHat
}
