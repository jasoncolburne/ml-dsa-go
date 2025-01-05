package mldsa

import (
	cryptorand "crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"math/bits"
	"math/rand"
	"testing"
)

func TestCoeffFromThreeBytes(t *testing.T) {
	result := coeffFromThreeBytes(ML_DSA_44_Parameters, byte(0x77), byte(0xe0), byte(0xff))
	if result != nil {
		t.Fatalf("non-nil result")
	}
}

func TestSimpleBitPackRoundtrip(t *testing.T) {
	b := int32((1 << 11) - 1)
	input := make([]int32, 256)
	for i := range 256 {
		input[i] = int32(rand.Intn(int(b)))
	}

	result := simpleBitUnpack(simpleBitPack(input, b), b)

	for i := range 256 {
		if input[i] != result[i] {
			t.Fatalf("mismatch at %d: %d != %d", i, input[i], result[i])
		}
	}
}

func TestIntegerToBitsRoundtrip(t *testing.T) {
	q := ML_DSA_44_Parameters.Q
	alpha := int32(bits.Len(uint(q - 1)))
	for i := range 100 {
		x := int32(rand.Intn(int(q)))
		y := bitsToInteger(integerToBits(x, alpha), alpha)
		if x != y {
			t.Fatalf("test[%d] failed, %d != %d", i, x, y)
		}
	}
}

func TestBytesToBitsRoundtrip(t *testing.T) {
	bytes := make([]byte, 128)
	cryptorand.Read(bytes)
	result := bitsToBytes(bytesToBits(bytes))

	if subtle.ConstantTimeCompare(bytes, result) != 1 {
		t.Fatalf("mismatch bytes %s vs result %s", hex.EncodeToString(bytes), hex.EncodeToString(result))
	}
}
