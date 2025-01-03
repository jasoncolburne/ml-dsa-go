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
