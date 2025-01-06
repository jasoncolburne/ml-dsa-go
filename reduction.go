package mldsa

func modMultiply(a, b, q int32) int32 {
	return int32((int64(a) * int64(b)) % int64(q))
}

func modQ(n, q int32) int32 {
	return (n%q + q) % q
}

func modQSymmetric(n, q int32) int32 {
	result := modQ(n, q)

	if result > q/2 {
		result -= q
	}

	return result
}

func vectorModQSymmetric(z [][]int32, q int32) [][]int32 {
	zModQSymmetric := make([][]int32, len(z))

	for i, row := range z {
		zModQSymmetric[i] = make([]int32, len(row))
		for j, value := range row {
			zModQSymmetric[i][j] = modQSymmetric(value, q)
		}
	}

	return zModQSymmetric
}

func power2Round(parameters ParameterSet, r int32) (int32, int32) {
	rPlus := modQ(r, parameters.Q)
	bound := int32(1) << parameters.D
	r0 := modQSymmetric(rPlus, bound)

	return (rPlus - r0) / bound, r0
}

func vectorPower2Round(parameters ParameterSet, t [][]int32) ([][]int32, [][]int32) {
	t0 := make([][]int32, parameters.K)
	t1 := make([][]int32, parameters.K)

	for j := range parameters.K {
		t0[j] = make([]int32, 256)
		t1[j] = make([]int32, 256)
		for i := range 256 {
			t1[j][i], t0[j][i] = power2Round(parameters, t[j][i])
		}
	}

	return t1, t0
}

func decompose(parameters ParameterSet, r int32) (int32, int32) {
	rPlus := modQ(r, parameters.Q)
	r0 := modQSymmetric(rPlus, 2*parameters.Gamma2)
	r1 := int32(0)

	if rPlus-r0 == parameters.Q-1 {
		r0 -= 1
	} else {
		r1 = (rPlus - r0) / (2 * parameters.Gamma2)
	}

	return r1, r0
}

func highBits(parameters ParameterSet, r int32) int32 {
	r1, _ := decompose(parameters, r)
	return r1
}

func lowBits(parameters ParameterSet, r int32) int32 {
	_, r0 := decompose(parameters, r)
	return r0
}

func makeHint(parameters ParameterSet, z, r int32) bool {
	r1 := highBits(parameters, r)
	v1 := highBits(parameters, r+z)

	return r1 != v1
}

func useHint(parameters ParameterSet, h bool, r int32) int32 {
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
