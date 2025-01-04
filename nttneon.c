//go:build neon

#include <arm_neon.h>
#include <stdint.h>

inline int64x2_t arm_vmulq_s64(const int64x2_t a, const int64x2_t b) {
   const int32x2_t ac = vmovn_s64(a);
   const int32x2_t pr = vmovn_s64(b);

   const int32x4_t hi = vmulq_s32(b, vrev64q_s32(a));

   return vmlal_u32(vshlq_n_s64(vpaddlq_u32(hi), 32), ac, pr);
}

int64x2_t reduce_mod_q_vec(int64x2_t v, int64_t q) {
    int64_t reduced[2];
    vst1q_s64(reduced, v);
    for (int i = 0; i < 2; i++) {
        reduced[i] = ((reduced[i] % q) + q) % q;
    }
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

        int32x4_t prod_vec = arm_vmulq_s64(a_vec, b_vec);
        prod_vec = reduce_mod_q_vec(prod_vec, q);

        vst1q_s64(&cHat[i], prod_vec);
    }
}
