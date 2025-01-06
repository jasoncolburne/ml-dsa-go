package mldsa

import (
	"crypto/subtle"

	"golang.org/x/crypto/sha3"
)

func keyGen(parameters ParameterSet, rnd []byte) (public []byte, private []byte, err error) {
	input := make([]byte, len(rnd))
	copy(input, rnd)
	input = append(input, integerToBytes(parameters.K, 1)...)
	input = append(input, integerToBytes(parameters.L, 1)...)

	inputHash := make([]byte, 128)
	sha3.ShakeSum256(inputHash, input)

	rho := inputHash[:32]
	rhoPrime := inputHash[32:96]
	kappa := inputHash[96:]

	AHat := expandA(parameters, rho)
	s1, s2 := expandS(parameters, rhoPrime)
	s1Hat := vectorNtt(parameters, s1)

	product := matrixVectorNtt(parameters, AHat, s1Hat)
	t := vectorAddPolynomials(parameters, vectorNttInverse(parameters, product), s2)
	t1, t0 := vectorPower2Round(parameters, t)

	pk := pkEncode(parameters, rho, t1)

	tr := make([]byte, 64)
	sha3.ShakeSum256(tr, pk)
	sk := skEncode(parameters, rho, kappa, tr, s1, s2, t0)

	return pk, sk, nil
}

func sign(parameters ParameterSet, sk, mPrime, rnd []byte) []byte {
	rho, kappa, tr, s1, s2, t0 := skDecode(parameters, sk)

	s1Hat := vectorNtt(parameters, s1)
	s2Hat := vectorNtt(parameters, s2)
	t0Hat := vectorNtt(parameters, t0)
	AHat := expandA(parameters, rho)

	inputHash := make([]byte, 64)
	copy(inputHash, tr)
	inputHash = append(inputHash, mPrime...)

	mu := make([]byte, 64)
	sha3.ShakeSum256(mu, inputHash)

	inputHash = make([]byte, 128)
	copy(inputHash[:32], kappa)
	copy(inputHash[32:64], rnd)
	copy(inputHash[64:], mu)

	rhoPrimePrime := make([]byte, 64)
	sha3.ShakeSum256(rhoPrimePrime, inputHash)

	k := int32(0)
	var z [][]int32
	var h [][]bool
	var cTilde []byte

	for z == nil && h == nil {
		y := expandMask(parameters, rhoPrimePrime, k)

		yHat := vectorNtt(parameters, y)
		product := matrixVectorNtt(parameters, AHat, yHat)

		w := vectorNttInverse(parameters, product)
		w1 := vectorHighBits(parameters, w)

		inputHash = make([]byte, 64)
		copy(inputHash, mu)
		inputHash = append(inputHash, w1Encode(parameters, w1)...)

		cTilde = make([]byte, parameters.Lambda/4)
		sha3.ShakeSum256(cTilde, inputHash)

		c := sampleInBall(parameters, cTilde)
		cHat := ntt(parameters, c)

		cs1 := vectorNttInverse(parameters, scalarVectorNtt(parameters, cHat, s1Hat))
		cs2 := vectorNttInverse(parameters, scalarVectorNtt(parameters, cHat, s2Hat))
		for i, row := range cs1 {
			for j, value := range row {
				cs1[i][j] = modQSymmetric(value, parameters.Q)
			}
		}

		z = vectorAddPolynomials(parameters, y, cs1)
		r := vectorSubtractPolynomials(parameters, w, cs2)

		zMax := vectorMaxAbsCoefficient(parameters, z, false)
		r0Max := vectorMaxAbsCoefficient(parameters, r, true)

		if zMax >= parameters.Gamma1-parameters.Beta || r0Max >= parameters.Gamma2-parameters.Beta {
			z = nil
			h = nil
		} else {
			ct0 := vectorNttInverse(parameters, scalarVectorNtt(parameters, cHat, t0Hat))

			ct0Neg := scalarVectorMultiply(parameters, -1, ct0)
			wPrime := vectorAddPolynomials(parameters, vectorSubtractPolynomials(parameters, w, cs2), ct0)

			h = vectorMakeHint(parameters, ct0Neg, wPrime)
			ct0Max := vectorMaxAbsCoefficient(parameters, ct0, false)
			if ct0Max >= parameters.Gamma2 || onesInH(h) > parameters.Omega {
				z = nil
				h = nil
			}
		}

		k += parameters.L
	}

	zModQSymmetric := vectorModQSymmetric(z, parameters.Q)
	sigma := sigEncode(parameters, cTilde, zModQSymmetric, h)
	return sigma
}

func verify(parameters ParameterSet, pk, mPrime, sigma []byte) bool {
	rho, t1 := pkDecode(parameters, pk)
	cTilde, z, h := sigDecode(parameters, sigma)

	if h == nil {
		return false
	}

	AHat := expandA(parameters, rho)
	tr := make([]byte, 64)
	sha3.ShakeSum256(tr, pk)

	inputHash := make([]byte, 64)
	copy(inputHash, tr)
	inputHash = append(inputHash, mPrime...)

	mu := make([]byte, 64)
	sha3.ShakeSum256(mu, inputHash)

	c := sampleInBall(parameters, cTilde)
	cHat := ntt(parameters, c)

	ct := scalarVectorNtt(parameters, cHat, vectorNtt(parameters, scalarVectorMultiply(parameters, 1<<parameters.D, t1)))
	Az := matrixVectorNtt(parameters, AHat, vectorNtt(parameters, z))
	Azct := subtractVectorNtt(parameters, Az, ct)

	wApproxPrime := vectorNttInverse(parameters, Azct)
	w1Prime := vectorUseHint(parameters, wApproxPrime, h)

	inputHash = make([]byte, 64)
	copy(inputHash, mu)
	inputHash = append(inputHash, w1Encode(parameters, w1Prime)...)

	cTildePrime := make([]byte, parameters.Lambda/4)
	sha3.ShakeSum256(cTildePrime, inputHash)

	zMax := vectorMaxAbsCoefficient(parameters, z, false)
	return zMax < (parameters.Gamma1-parameters.Beta) && subtle.ConstantTimeCompare(cTilde, cTildePrime) == 1
}
