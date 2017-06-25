package murmur3

func rol32(v uint32, l uint) uint32 {
	l &= 31
	return (v << l) | (v >> (32 - l))
}

func ror32(v uint32, r uint) uint32 {
	r &= 31
	return (v >> r) | (v << (32 - r))
}
