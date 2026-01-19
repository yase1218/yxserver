package utils

import "time"

var (
	FastRandUtil = NewFastRand(uint64(time.Now().UnixNano()))
)

type FastRand struct {
	seed uint64
}

func NewFastRand(seed uint64) *FastRand {
	return &FastRand{seed: seed}
}

func (r *FastRand) Intn(n int) int {
	v := r.Uint32()
	return int((uint64(v) * uint64(n)) >> 32)
}

func (r *FastRand) Uint32() uint32 {
	r.seed ^= r.seed << 13
	r.seed ^= r.seed >> 17
	r.seed ^= r.seed << 5
	return uint32(r.seed)
}
