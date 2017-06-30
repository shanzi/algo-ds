package murmur3

import (
	"hash"
	"unsafe"
)

type Hash128 interface {
	hash.Hash

	Sum128() (uint64, uint64)
}

type digest128 struct {
	len   int
	hash1 uint64
	hash2 uint64
	tail  []byte
}

func New128() Hash128 {
	return &digest128{0, 0, 0, nil}
}

func (self *digest128) Size() int {
	return 8
}

func (self *digest128) BlockSize() int {
	return 1
}

func (self *digest128) Reset() {
	self.len = 0
	self.hash1 = 0
	self.hash2 = 0
	self.tail = nil
}

func (self *digest128) Write(b []byte) (int, error) {
	if len(b) <= 0 {
		return 0, nil
	}

	if self.tail != nil {
		if len(self.tail)+len(b) <= 15 {
			self.tail = append(self.tail, b...)
			self.len += len(b)
		} else {
			r := 16 - len(self.tail)
			bb := append(self.tail, b[:r]...)

			self.len -= len(self.tail)
			self.tail = nil

			self.Write(bb)
			self.Write(b[r:])
		}
		return len(b), nil
	}

	if len(b) <= 15 {
		self.tail = make([]byte, len(b))
		copy(self.tail, b)
		self.len += len(b)
		return len(b), nil
	}

	nblocks := len(b) >> 4
	h1, h2 := self.hash1, self.hash2
	p := unsafe.Pointer(&b[0])

	for i := 0; i < nblocks; i++ {
		t := (*[2]uint64)(p)
		k1, k2 := t[0], t[1]

		k1 *= c1_128
		k1 = rol64(k1, 31)
		k1 *= c2_128
		h1 ^= k1

		h1 = rol64(h1, 27)
		h1 += h2
		h1 = h1*5 + c5_128

		k2 *= c2_128
		k2 = rol64(k2, 33)
		k2 *= c1_128
		h2 ^= k2

		h2 = rol64(h2, 31)
		h2 += h1
		h2 = h2*5 + c6_128

		p = unsafe.Pointer(uintptr(p) + 16)
	}

	if len(b)&15 != 0 {
		self.tail = make([]byte, len(b)&15)
		copy(self.tail, b[nblocks*16:])
	}

	self.hash1 = h1
	self.hash2 = h2
	self.len += len(b)

	return len(b), nil
}

func (self *digest128) Sum(b []byte) []byte {
	h1, h2 := self.Sum128()

	for i := 56; i >= 0; i -= 8 {
		b = append(b, byte(h1>>uint(i)))
	}

	for i := 56; i >= 0; i -= 8 {
		b = append(b, byte(h2>>uint(i)))
	}

	return b
}

func (self *digest128) Sum128() (uint64, uint64) {
	h1, h2 := self.hash1, self.hash2

	if len(self.tail) > 8 {
		k2 := fetch64(self.tail[8:])
		k2 *= c2_128
		k2 = rol64(k2, 33)
		k2 *= c1_128
		h2 ^= k2
	}

	if len(self.tail) > 0 {
		l := len(self.tail)
		if l > 8 {
			l = 8
		}

		k1 := fetch64(self.tail[:l])
		k1 *= c1_128
		k1 = rol64(k1, 31)
		k1 *= c2_128
		h1 ^= k1
	}

	h1 ^= uint64(self.len)
	h2 ^= uint64(self.len)

	h1 += h2
	h2 += h1

	h1 = fmix64(h1)
	h2 = fmix64(h2)

	h1 += h2
	h2 += h1
	return h1, h2
}

func fmix64(h uint64) uint64 {
	h ^= h >> 33
	h *= c3_128
	h ^= h >> 33
	h *= c4_128
	h ^= h >> 33
	return h
}
