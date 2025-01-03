## This repo is in its infancy

But, the implementation should be complete as per [NIST](https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.204.pdf).

It's possible there are a few timing vulnerabilities, and I've only just learned these concepts.

However, I did my best to implement as instructed.

Future improvements:
- [ ] Stop using Sha256 to provide entropy, use Sha3-256 or Blake3 or something approved and better
- [ ] Use hardware for FFT math
