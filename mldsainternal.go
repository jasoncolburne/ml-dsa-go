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

	t := make([][]int32, parameters.K)
	for j := range parameters.K {
		t[j] = addPolynomials(parameters, nttInverse(parameters, product[j]), s2[j])
	}

	t0 := make([][]int32, parameters.K)
	t1 := make([][]int32, parameters.K)
	for j := range parameters.K {
		t0[j] = make([]int32, 256)
		t1[j] = make([]int32, 256)
		for i := range 256 {
			t1[j][i], t0[j][i] = power2Round(parameters, t[j][i])
		}
	}

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

		w := make([][]int32, parameters.K)
		for j, polynomial := range product {
			w[j] = nttInverse(parameters, polynomial)
		}

		w1 := make([][]int32, parameters.K)
		for j, row := range w {
			w1[j] = make([]int32, 256)
			for i, value := range row {
				w1[j][i] = highBits(parameters, value)
			}
		}

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

		r0Max := int32(0)
		for _, polynomial := range r {
			for _, value := range polynomial {
				r0 := int32(lowBits(parameters, value))
				if r0 < 0 {
					r0 *= -1
				}

				if r0Max < r0 {
					r0Max = r0
				}
			}
		}

		zMax := int32(0)
		for _, polynomial := range z {
			for _, value := range polynomial {
				zValue := value

				if zValue < 0 {
					zValue *= -1
				}

				if zMax < zValue {
					zMax = zValue
				}
			}
		}

		if zMax >= parameters.Gamma1-parameters.Beta || r0Max >= parameters.Gamma2-parameters.Beta {
			z = nil
			h = nil
		} else {
			ct0 := vectorNttInverse(parameters, scalarVectorNtt(parameters, cHat, t0Hat))
			ct0Neg := scalarVectorMultiply(parameters, -1, ct0)

			wPrime := vectorAddPolynomials(parameters, vectorSubtractPolynomials(parameters, w, cs2), ct0)
			h = make([][]bool, len(ct0Neg))

			for i, ct0NegValues := range ct0Neg {
				h[i] = make([]bool, len(ct0NegValues))
				for j, value := range ct0NegValues {
					h[i][j] = makeHint(parameters, value, wPrime[i][j])
				}
			}

			ct0Max := int32(0)
			for _, row := range ct0 {
				for _, value := range row {
					ct0 := value

					if ct0 < 0 {
						ct0 *= -1
					}

					if ct0Max < ct0 {
						ct0Max = ct0
					}
				}
			}

			onesInH := int32(0)
			for _, row := range h {
				for _, value := range row {
					if value {
						onesInH += 1
					}
				}
			}

			if ct0Max >= parameters.Gamma2 || onesInH > parameters.Omega {
				z = nil
				h = nil
			}
		}

		k += parameters.L
	}

	zModQCentered := make([][]int32, len(z))
	for i, row := range z {
		zModQCentered[i] = make([]int32, len(row))

		for j, value := range row {
			zModQCentered[i][j] = modQSymmetric(value, parameters.Q)
		}
	}

	sigma := sigEncode(parameters, cTilde, zModQCentered, h)
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

	wApproxPrime := make([][]int32, parameters.K)
	for i, value := range Azct {
		wApproxPrime[i] = nttInverse(parameters, value)
	}

	w1Prime := make([][]int32, parameters.K)
	for i, row := range wApproxPrime {
		w1Prime[i] = make([]int32, len(row))
		for j, value := range row {
			w1Prime[i][j] = useHint(parameters, h[i][j], value)
		}
	}

	inputHash = make([]byte, 64)
	copy(inputHash, mu)
	inputHash = append(inputHash, w1Encode(parameters, w1Prime)...)

	cTildePrime := make([]byte, parameters.Lambda/4)
	sha3.ShakeSum256(cTildePrime, inputHash)

	zMax := int32(0)
	for _, row := range z {
		for _, value := range row {
			zValue := value

			if zValue < 0 {
				zValue *= -1
			}

			if zMax < zValue {
				zMax = zValue
			}
		}
	}

	return zMax < (parameters.Gamma1-parameters.Beta) && subtle.ConstantTimeCompare(cTilde, cTildePrime) == 1
}
