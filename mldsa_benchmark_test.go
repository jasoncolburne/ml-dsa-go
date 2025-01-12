package mldsa_test

import (
	"testing"

	mldsa "github.com/jasoncolburne/ml-dsa-go"
)

func benchmarkMLDSAGenerate(parameters mldsa.ParameterSet, b *testing.B) {
	dsa := mldsa.Init(parameters)

	b.ResetTimer()
	for range b.N {
		dsa.KeyGen()
	}
	b.StopTimer()

	opsPerSec := float64(b.N) / b.Elapsed().Seconds()
	b.ReportMetric(opsPerSec, "ops/s")
}

func benchmarkMLDSASign(parameters mldsa.ParameterSet, b *testing.B) {
	dsa := mldsa.Init(parameters)
	_, sk, _ := dsa.KeyGen()
	message := []byte("fabulous message")
	ctx := []byte("context")

	b.ResetTimer()
	for range b.N {
		dsa.Sign(sk, message, ctx)
	}
	b.StopTimer()

	opsPerSec := float64(b.N) / b.Elapsed().Seconds()
	b.ReportMetric(opsPerSec, "ops/s")
}

func benchmarkMLDSAVerify(parameters mldsa.ParameterSet, b *testing.B) {
	dsa := mldsa.Init(parameters)
	pk, sk, _ := dsa.KeyGen()
	message := []byte("fabulous message")
	ctx := []byte("context")
	sig, _ := dsa.Sign(sk, message, ctx)

	b.ResetTimer()
	for range b.N {
		dsa.Verify(pk, message, sig, ctx)
	}
	b.StopTimer()

	opsPerSec := float64(b.N) / b.Elapsed().Seconds()
	b.ReportMetric(opsPerSec, "ops/s")
}

func BenchmarkMLDSA44Generate(b *testing.B) {
	benchmarkMLDSAGenerate(mldsa.ML_DSA_44_Parameters, b)
}

func BenchmarkMLDSA65Generate(b *testing.B) {
	benchmarkMLDSAGenerate(mldsa.ML_DSA_65_Parameters, b)
}

func BenchmarkMLDSA87Generate(b *testing.B) {
	benchmarkMLDSAGenerate(mldsa.ML_DSA_87_Parameters, b)
}

func BenchmarkMLDSA44Sign(b *testing.B) {
	benchmarkMLDSASign(mldsa.ML_DSA_44_Parameters, b)
}

func BenchmarkMLDSA65Sign(b *testing.B) {
	benchmarkMLDSASign(mldsa.ML_DSA_65_Parameters, b)
}

func BenchmarkMLDSA87Sign(b *testing.B) {
	benchmarkMLDSASign(mldsa.ML_DSA_87_Parameters, b)
}

func BenchmarkMLDSA44Verify(b *testing.B) {
	benchmarkMLDSAVerify(mldsa.ML_DSA_44_Parameters, b)
}

func BenchmarkMLDSA65Verify(b *testing.B) {
	benchmarkMLDSAVerify(mldsa.ML_DSA_65_Parameters, b)
}

func BenchmarkMLDSA87Verify(b *testing.B) {
	benchmarkMLDSAVerify(mldsa.ML_DSA_87_Parameters, b)
}
