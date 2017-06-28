package hamt

const (
	m1_32 = 0x55555555
	m2_32 = 0x33333333
	m4_32 = 0x0f0f0f0f
)

func popcount32(x uint32) int {
	x -= (x >> 1) & m1_32
	x = (x & m2_32) + ((x >> 2) & m2_32)
	x = (x + (x >> 4)) & m4_32
	x += x >> 8
	x += x >> 16
	return int(x & 0x3f)
}
