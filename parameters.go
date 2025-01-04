package mldsa

type ParameterSet struct {
	Q      int32 // modulus
	Zeta   int32 // a 512th root of unity in Zq
	D      int32 // # of dropped bits from t
	Tau    int32 // # of +/-1s in polynomial c
	Lambda int32 // collision strength of c~
	Gamma1 int32 // coefficient range of y
	Gamma2 int32 // low order rounding range
	K      int32 // rows in A
	L      int32 // columns in A
	Eta    int32 // private key range
	Beta   int32 // Tau * Eta
	Omega  int32 // max # of 1s in the hint h
}

var (
	ML_DSA_44_Parameters = ParameterSet{
		Q:      8380417,
		Zeta:   1753,
		D:      13,
		Tau:    39,
		Lambda: 128,
		Gamma1: 131072,
		Gamma2: 95232,
		K:      4,
		L:      4,
		Eta:    2,
		Beta:   78,
		Omega:  80,
	}

	ML_DSA_65_Parameters = ParameterSet{
		Q:      8380417,
		Zeta:   1753,
		D:      13,
		Tau:    49,
		Lambda: 192,
		Gamma1: 524288,
		Gamma2: 261888,
		K:      6,
		L:      5,
		Eta:    4,
		Beta:   196,
		Omega:  55,
	}

	ML_DSA_87_Parameters = ParameterSet{
		Q:      8380417,
		Zeta:   1753,
		D:      13,
		Tau:    60,
		Lambda: 256,
		Gamma1: 524288,
		Gamma2: 261888,
		K:      8,
		L:      7,
		Eta:    2,
		Beta:   120,
		Omega:  75,
	}
)
