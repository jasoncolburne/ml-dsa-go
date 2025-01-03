package mldsa

import (
	"crypto/subtle"
	"fmt"

	"golang.org/x/crypto/sha3"
)

const SEEDLENGTH = 32

func KeyGen(parameters ParameterSet) (public []byte, private []byte, err error) {
	rnd, err := rbg(SEEDLENGTH)
	if err != nil {
		return nil, nil, err
	}

	return keyGen(parameters, rnd)
}

func Sign(parameters ParameterSet, sk, message, ctx []byte) ([]byte, error) {
	if len(ctx) > 255 {
		return nil, fmt.Errorf("ctx length > 255")
	}

	rnd, err := rbg(SEEDLENGTH)
	if err != nil {
		return nil, err
	}

	// for KAT testing
	// rnd = make([]byte, SEEDLENGTH)

	mPrime := integerToBytes(0, 1)
	mPrime = append(mPrime, integerToBytes(len(ctx), 1)...)
	mPrime = append(mPrime, ctx...)
	mPrime = append(mPrime, message...)

	sigma := sign(parameters, sk, mPrime, rnd)
	return sigma, nil
}

func Verify(parameters ParameterSet, pk, message, signature, ctx []byte) (bool, error) {
	if len(ctx) > 255 {
		return false, fmt.Errorf("ctx length > 255")
	}

	mPrime := integerToBytes(0, 1)
	mPrime = append(mPrime, integerToBytes(len(ctx), 1)...)
	mPrime = append(mPrime, ctx...)
	mPrime = append(mPrime, message...)

	return verify(parameters, pk, mPrime, signature), nil
}

func keyGen(parameters ParameterSet, rnd []byte) (public []byte, private []byte, err error) {
	// for KAT testing, we fix the seed
	// _ = rnd
	// // hexStr := "f696484048ec21f96cf50a56d0759c448f3779752f0383d37449690694cf7a68"
	// // hexStr := "6de62e3465a55c9c78a07d265be8540b3e58b0801a124d07ff12b438d5202ea0"
	// // hexStr := "1eaae6bb91b27cd748c402c4111140d5a942cf3c95ff7977f88d2ef515bb26d0"
	// hexStr := "b585d4eb01085111a172a87688d0032e3381a9e9a35fdd6ef2f8aeb3b40eb5ce"
	// input, _ := hex.DecodeString(hexStr)

	input := rnd
	input = append(input, integerToBytes(parameters.K, 1)...)
	input = append(input, integerToBytes(parameters.L, 1)...)

	inputHash := make([]byte, 128)
	sha3.ShakeSum256(inputHash, input)

	rho := inputHash[:32]
	rhoPrime := inputHash[32:96]
	kappa := inputHash[96:]

	AHat := expandA(parameters, rho)
	// fmt.Printf("AHat: %v\n", AHat)
	s1, s2 := expandS(parameters, rhoPrime)
	// fmt.Printf("s1: %v\n", s1)
	// fmt.Printf("s2: %v\n", s2)
	s1Hat := vectorNtt(parameters, s1)

	product := matrixVectorNtt(parameters, AHat, s1Hat)

	t := make([][]int, parameters.K)
	for j := range parameters.K {
		// fmt.Printf("product[%d]: %v\n", j, product[j])
		t[j] = addPolynomials(nttInverse(parameters, product[j]), s2[j])
	}

	t0 := make([][]int, parameters.K)
	t1 := make([][]int, parameters.K)
	for j := range parameters.K {
		t0[j] = make([]int, 256)
		t1[j] = make([]int, 256)
		for i := range 256 {
			t1[j][i], t0[j][i] = power2Round(parameters, t[j][i])
		}
	}

	pk := pkEncode(parameters, rho, t1)
	// fmt.Printf("rho: %v\n", rho)
	// fmt.Printf("t1: %v\n", t1)

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

	k := 0
	var z [][]int
	var h [][]bool
	var cTilde []byte

	for z == nil && h == nil {
		y := expandMask(parameters, rhoPrimePrime, k)

		yHat := vectorNtt(parameters, y)
		product := matrixVectorNtt(parameters, AHat, yHat)

		w := make([][]int, parameters.K)
		for j, polynomial := range product {
			w[j] = nttInverse(parameters, polynomial)
		}

		w1 := make([][]int, parameters.K)
		for j, row := range w {
			w1[j] = make([]int, 256)
			for i, value := range row {
				w1[j][i] = highBits(parameters, value)
			}
		}

		inputHash = make([]byte, 64)
		copy(inputHash, mu)
		inputHash = append(inputHash, w1Encode(parameters, w1)...)
		// fmt.Printf("w1: %v\n", w1)

		cTilde = make([]byte, parameters.Lambda/4)
		sha3.ShakeSum256(cTilde, inputHash)

		c := sampleInBall(parameters, cTilde)
		// fmt.Printf("c: %v\n", c)
		cHat := ntt(parameters, c)

		cs1 := vectorNttInverse(parameters, scalarVectorNtt(parameters, cHat, s1Hat))
		cs2 := vectorNttInverse(parameters, scalarVectorNtt(parameters, cHat, s2Hat))
		for i, row := range cs1 {
			for j, value := range row {
				cs1[i][j] = modCentered(value, parameters.Q)
			}
		}

		z = vectorAddPolynomials(y, cs1)
		// fmt.Printf("z: %v\n", z)
		// fmt.Printf("y: %v\n", y)
		// fmt.Printf("c: %v\n", c)
		// fmt.Printf("cs1: %v\n", cs1)
		r := vectorSubtractPolynomials(w, cs2)

		r0Max := 0
		for _, polynomial := range r {
			for _, value := range polynomial {
				r0 := lowBits(parameters, value)
				if r0 < 0 {
					r0 *= -1
				}

				if r0Max < r0 {
					r0Max = r0
				}
			}
		}

		zMax := 0
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
			// fmt.Printf("%d:%d %d:%d\n", zMax, parameters.Gamma1-parameters.Beta, r0Max, parameters.Gamma2-parameters.Beta)
			z = nil
			h = nil
		} else {
			ct0 := vectorNttInverse(parameters, scalarVectorNtt(parameters, cHat, t0Hat))
			ct0Neg := scalarVectorMultiply(-1, ct0)

			wPrime := vectorAddPolynomials(vectorSubtractPolynomials(w, cs2), ct0)
			h = make([][]bool, len(ct0Neg))

			for i, ct0NegValues := range ct0Neg {
				h[i] = make([]bool, len(ct0NegValues))
				for j, value := range ct0NegValues {
					h[i][j] = makeHint(parameters, value, wPrime[i][j])
				}
			}

			ct0Max := 0
			for _, row := range ct0 {
				for _, value := range row {
					if ct0Max < value {
						ct0Max = value
					}
				}
			}

			onesInH := 0
			for _, row := range h {
				for _, value := range row {
					if value {
						onesInH += 1
					}
				}
			}

			if ct0Max >= parameters.Gamma2 || onesInH > parameters.Omega {
				// fmt.Printf("%d:%d %d:%d\n", ct0Max, parameters.Gamma2, onesInH, parameters.Omega)
				z = nil
				h = nil
			}
		}

		k += parameters.L
	}

	zModQCentered := make([][]int, len(z))
	for i, row := range z {
		zModQCentered[i] = make([]int, len(row))

		for j, value := range row {
			zModQCentered[i][j] = modCentered(value, parameters.Q)
		}
	}

	sigma := sigEncode(parameters, cTilde, zModQCentered, h)
	// fmt.Printf("cTilde: %v\n", cTilde)
	// fmt.Printf("zModQCentered: %v\n", zModQCentered)
	// fmt.Printf("h: %v\n", h)

	return sigma
}

func verify(parameters ParameterSet, pk, mPrime, sigma []byte) bool {
	rho, t1 := pkDecode(parameters, pk)
	// fmt.Printf("t1: %v\n", t1)
	// fmt.Printf("rho: %v\n", rho)
	cTilde, z, h := sigDecode(parameters, sigma)
	// fmt.Printf("cTilde: %v\n", cTilde)
	// fmt.Printf("z: %v\n", z)
	// fmt.Printf("h: %v\n", h)

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
	// fmt.Printf("c: %v\n", c)

	ct := scalarVectorNtt(parameters, cHat, vectorNtt(parameters, scalarVectorMultiply(1<<parameters.D, t1)))
	Az := matrixVectorNtt(parameters, AHat, vectorNtt(parameters, z))
	Azct := subtractVectorNtt(parameters, Az, ct)

	wApproxPrime := make([][]int, parameters.K)
	for i, value := range Azct {
		wApproxPrime[i] = nttInverse(parameters, value)
	}
	// fmt.Printf("wApproxPrime: %v\n", wApproxPrime)

	w1Prime := make([][]int, parameters.K)
	for i, row := range wApproxPrime {
		w1Prime[i] = make([]int, len(row))
		for j, value := range row {
			w1Prime[i][j] = useHint(parameters, h[i][j], value)
		}
	}

	inputHash = make([]byte, 64)
	copy(inputHash, mu)
	inputHash = append(inputHash, w1Encode(parameters, w1Prime)...)
	// fmt.Printf("w1Prime: %v\n", w1Prime)

	cTildePrime := make([]byte, parameters.Lambda/4)
	sha3.ShakeSum256(cTildePrime, inputHash)

	zMax := 0
	for _, row := range z {
		for _, value := range row {
			if zMax < value {
				zMax = value
			}
		}
	}

	// fmt.Printf("%d, %d, %v, %v\n", zMax, parameters.Gamma1-parameters.Beta, cTilde, cTildePrime)

	return zMax < (parameters.Gamma1-parameters.Beta) && subtle.ConstantTimeCompare(cTilde, cTildePrime) == 1
}
