package hamt

import "unsafe"

// Murmur3 32-bit version
func sum32(data []byte, seed uint32) uint32 {

	var h1 = seed

	nblocks := len(data) / 4
	var p uintptr
	if len(data) > 0 {
		p = uintptr(unsafe.Pointer(&data[0]))
	}
	p1 := p + uintptr(4*nblocks)
	for ; p < p1; p += 4 {
		k1 := *(*uint32)(unsafe.Pointer(p))

		k1 *= 0xcc9e2d51
		k1 = (k1 << 15) | (k1 >> 17)
		k1 *= 0x1b873593

		h1 ^= k1
		h1 = (h1 << 13) | (h1 >> 19)
		h1 = h1*5 + 0xe6546b64
	}

	tail := data[nblocks*4:]

	var k1 uint32
	switch len(tail) & 3 {
	case 3:
		k1 ^= uint32(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint32(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint32(tail[0])
		k1 *= 0xcc9e2d51
		k1 = (k1 << 15) | (k1 >> 17)
		k1 *= 0x1b873593
		h1 ^= k1
	}

	h1 ^= uint32(len(data))

	h1 ^= h1 >> 16
	h1 *= 0x85ebca6b
	h1 ^= h1 >> 13
	h1 *= 0xc2b2ae35
	h1 ^= h1 >> 16

	return h1
}

func keyHash32(s string, seed uint32) uint32 {
	b := ([]byte)(s)
	h := sum32(b, seed)

	// The most significant two bits will be abandoned during insert and lookup,
	// Thus we'd better mix them down with the least significant two bits
	return ((h >> 30) & 3) ^ h
}
