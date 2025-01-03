package mldsa

func modCentered(n, q int) int {
	result := modQ(n, q)

	if result > q/2 {
		result -= q
	}

	return result
}

func modQ(n, q int) int {
	return (n%q + q) % q
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
