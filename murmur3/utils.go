package murmur3

func rol32(v uint32, l uint) uint32 {
	l &= 31
	return (v << l) | (v >> (32 - l))
}

func rol64(v uint64, r uint) uint64 {
	r &= 63
	return (v << r) | (v >> (64 - r))
}

func fetch64(b []byte) uint64 {
	k := uint64(0)
	for i := len(b) - 1; i >= 0; i-- {
		k |= uint64(b[i]) << uint(i*8)
	}
	return k
}
