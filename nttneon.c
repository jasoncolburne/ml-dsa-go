//go:build neon

#include <arm_neon.h>
#include <stdint.h>

inline int64x2_t vmulq_s64_apple(const int64x2_t a, const int64x2_t b) {
   const int32x2_t ac = vmovn_s64(a);
   const int32x2_t pr = vmovn_s64(b);

   const int32x4_t hi = vmulq_s32(b, vrev64q_s32(a));

   return vmlal_u32(vshlq_n_s64(vpaddlq_u32(hi), 32), ac, pr);
}

int64x2_t vmodq_s64(int64x2_t v, int q) {
    int64_t reduced[2];
    vst1q_s64(reduced, v);

    reduced[0] = ((reduced[0] % q) + q) % q;
    reduced[1] = ((reduced[1] % q) + q) % q;

    return vld1q_s64(reduced);
}

void add_ntt(const int32_t *aHat, const int32_t *bHat, int32_t *cHat, int32_t q) {
    for (int i = 0; i < 256; i += 4) {
        int32x4_t a_vec = vld1q_s32(&aHat[i]);
        int32x4_t b_vec = vld1q_s32(&bHat[i]);
        int32x4_t sum_vec = vaddq_s32(a_vec, b_vec);
        vst1q_s32(&cHat[i], sum_vec);
    }
}

void subtract_ntt(const int32_t *aHat, const int32_t *bHat, int32_t *cHat, int32_t q) {
    for (int i = 0; i < 256; i += 4) {
        int32x4_t a_vec = vld1q_s32(&aHat[i]);
        int32x4_t b_vec = vld1q_s32(&bHat[i]);
        int32x4_t diff_vec = vsubq_s32(a_vec, b_vec);
        vst1q_s32(&cHat[i], diff_vec);
    }
}

void multiply_ntt(const int64_t *aHat, const int64_t *bHat, int64_t *cHat, int32_t q) {
    for (int i = 0; i < 256; i += 2) {
        int64x2_t a_vec = vld1q_s64(&aHat[i]);
        int64x2_t b_vec = vld1q_s64(&bHat[i]);
#ifdef __APPLE__
        int64x2_t prod_vec = vmodq_s64(vmulq_s64_apple(a_vec, b_vec), q);
#else
        int64x2_t prod_vec = vmodq_s64(vmulq_s64(a_vec, b_vec), q);
#endif
        vst1q_s64(&cHat[i], prod_vec);
    }
}
