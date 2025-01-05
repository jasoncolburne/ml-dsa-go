## This repo is in its infancy

But, the implementation should be complete as per [NIST](https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.204.pdf).

It's possible there are a few timing vulnerabilities, and I've only just learned these ML-DSA concepts.

However, I did my best to implement as instructed.

Future improvements:
- [ ] Stop using Sha256 to provide entropy, use Sha3-256 or Blake3 or something approved and better
- [ ] Use hardware for NTT math
  - [ ] NEON (arm) update: partial implementation sees negative performance gain, but maybe I'm just doing it wrong
  - [ ] AVX (x86)
- [ ] Audit for side channel attacks
- [ ] Zero sensitive data containers
