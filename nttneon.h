#include <stdint.h>

void add_ntt(const int32_t *aHat, const int32_t *bHat, int32_t *cHat, int32_t q);
void subtract_ntt(const int32_t *aHat, const int32_t *bHat, int32_t *cHat, int32_t q);
