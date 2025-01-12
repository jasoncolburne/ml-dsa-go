## This repo is in its infancy

But, the implementation should be complete as per [NIST](https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.204.pdf).

It's possible there are a few timing vulnerabilities, and I've only just learned these ML-DSA concepts.

I may have done some silly things to try and prevent optimization without understanding the go
compiler when I can likely use a flag or something.

## Makefile

Roundtrip happy/sad and KAT (Known Answer Tests)
```
make test
```

Using NEON instructions
```
make test-neon
```

Performance benchmarks
```
make benchmarks
```

Formatting
```
make format
```

## Future improvements
- [x] Stop using Sha256 to provide entropy, use Sha3-256 or Blake3 or something approved and better (chose SHA3-512, I dunno I don't think it's necessary to use 512 but meh)
- [ ] Use hardware for NTT math
  - [x] NEON (arm) update: partial implementation sees negative performance gain, but maybe I'm just doing it wrong. if you want to try building it, just build with `-tags=neon`. not planning on investing more time here.
  - [ ] AVX (x86)
- [ ] Audit for side channel attacks
  - I've done a bit of this now, I added some else clauses and tried to make everything constant time but there are a couple TODOs remaining that require some thought
- [ ] Zero sensitive data containers
