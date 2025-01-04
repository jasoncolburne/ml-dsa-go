//go:build neon

#include <arm_neon.h>
#include <stdint.h>
#include <stdio.h>

int32x4_t reduce_mod_q_vec(int32x4_t v, int32_t q) {
    int32_t reduced[4];
    vst1q_s32(reduced, v);
    for (int i = 0; i < 4; i++) {
        // printf("%d:", reduced[i]);
        reduced[i] = ((reduced[i] % q) + q) % q;
        // printf("%d\n", reduced[i]);
    }
    return vld1q_s32(reduced);
}

void add_ntt(const int32_t *aHat, const int32_t *bHat, int32_t *cHat, int32_t q) {
    for (int i = 0; i < 256; i += 4) {
        int32x4_t a_vec = vld1q_s32(&aHat[i]);
        int32x4_t b_vec = vld1q_s32(&bHat[i]);
        int32x4_t sum_vec = vaddq_s32(a_vec, b_vec);
        sum_vec = reduce_mod_q_vec(sum_vec, q);
        vst1q_s32(&cHat[i], sum_vec);
    }
}

// void matrix_vector_ntt(const int32 *MHat, const int32 *vHat, int32 *wHat, int32 K, int32 L, int32 q) {
//     for (int32 i = 0; i < K; i++) {
//         int32x4_t sum_vec[64];
//         for (int32 v = 0; v < 64; v++) {
//             sum_vec[v] = vdupq_n_s32(0);
//         }

//         for (int32 j = 0; j < L; j++) {
//             for (int32 k = 0; k < 256; k += 4) {
//                 int32x4_t mat_row_vec = vld1q_s32(&MHat[(i * L * 256) + (j * 256) + k]);
//                 int32x4_t vec_elem_vec = vld1q_s32(&vHat[(j * 256) + k]);

//                 int32x4_t product_vec = vmulq_s32(mat_row_vec, vec_elem_vec);
//                 reduce_mod_q_vec(product_vec, q);

//                 sum_vec[k / 4] = vaddq_s32(sum_vec[k / 4], product_vec);
//                 sum_vec[k / 4] = reduce_mod_q_vec(sum_vec[k / 4], q);
//               }
//         }

//         for (int32 k = 0; k < 256; k += 4) {
//             vst1q_s32(&wHat[(i * 256) + k], sum_vec[k / 4]);
//         }
//     }
// }
