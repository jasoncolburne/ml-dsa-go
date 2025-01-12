package mldsa

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"

	"golang.org/x/crypto/sha3"
)

type DRBG struct {
	seedLength    int32
	reseedCounter int64

	V []byte
	C []byte
}

func (drbg *DRBG) Init(entropy []byte, personalizationString []byte) {
	drbg.seedLength = 888

	seedMaterial := append(entropy, personalizationString...)
	drbg.reseed(seedMaterial)
}

func (drbg *DRBG) reseed(seedMaterial []byte) {
	drbg.V = drbg.derive(seedMaterial, drbg.seedLength)
	drbg.C = drbg.derive(append([]byte{1}, drbg.V...), drbg.seedLength)
	drbg.reseedCounter += 1
}

func (drbg *DRBG) derive(input []byte, numberOfBits int32) []byte {
	numberOfBytes := (numberOfBits + 7) / 8

	output := []byte{}
	temp := input[:]

	offset := int32(0)

	for offset < numberOfBytes {
		temp = append([]byte{1}, temp...)
		dig := sha3.Sum512(temp)
		digLength := int32(64)

		var bytesToWrite int32
		if digLength < numberOfBytes-offset {
			bytesToWrite = digLength
		} else {
			bytesToWrite = numberOfBytes - offset
		}

		output = append(output, dig[:bytesToWrite]...)
		offset += bytesToWrite
	}

	return output
}

func (drbg *DRBG) Reseed(entropy []byte) {
	seedMaterial := append([]byte{1}, drbg.V...)
	seedMaterial = append(seedMaterial, entropy...)
	drbg.reseed(seedMaterial)
}

func (drbg *DRBG) Generate(numberOfBits int32) ([]byte, error) {
	if drbg.reseedCounter > int64(1<<48) {
		return nil, fmt.Errorf("must reseed")
	}

	numberOfBytes := (numberOfBits + 7) / 8
	output := []byte{}

	offset := int32(0)
	for offset < numberOfBytes {
		sum512 := sha3.Sum512(drbg.V)
		drbg.V = sum512[:]
		digLength := int32(64)

		var bytesToWrite int32
		if digLength < numberOfBytes-offset {
			bytesToWrite = digLength
		} else {
			bytesToWrite = numberOfBytes - offset
		}

		output = append(output, drbg.V[:bytesToWrite]...)
		offset += bytesToWrite
	}

	V := bytesToBigInt(drbg.V)
	C := bytesToBigInt(drbg.C)
	sum := &big.Int{}

	sum.Add(V, C)
	X := *sum
	V = &X
	sum.Add(V, big.NewInt(int64(drbg.reseedCounter)))
	drbg.V = bigIntToBytes(sum, int(drbg.seedLength/8))

	drbg.reseedCounter += 1

	return output, nil
}

func bytesToBigInt(bytes []byte) *big.Int {
	hexStr := hex.EncodeToString(bytes)
	bigInt := new(big.Int)
	bigInt.SetString(hexStr, 16)
	return bigInt
}

func bigIntToBytes(number *big.Int, length int) []byte {
	hexStr := number.Text(16)
	hexStr = padLeft(hexStr, length*2, '0')

	bytes, _ := hex.DecodeString(hexStr)

	if len(bytes) < length {
		paddedBytes := make([]byte, length)
		copy(paddedBytes[length-len(bytes):], bytes)
		return paddedBytes
	}
	return bytes[:length]
}

func padLeft(str string, length int, pad byte) string {
	if len(str) >= length {
		return str
	}
	padding := make([]byte, length-len(str))
	for i := range padding {
		padding[i] = pad
	}
	return string(padding) + str
}

func rbg(len int32) ([]byte, error) {
	entropy := make([]byte, 32)
	_, err := rand.Read(entropy)
	if err != nil {
		return nil, err
	}

	hashDRBG := &DRBG{}
	hashDRBG.Init(entropy, []byte("jasoncolburne/ml-dsa-go"))

	bytes, err := hashDRBG.Generate(len * 8)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
