package murmur3

import (
	"hash"
	"unsafe"
)

type digest32 struct {
	len  int
	seed uint32
	hash uint32
	tail []byte
}

func New32(seed uint32) hash.Hash32 {
	return &digest32{0, seed, seed, nil}
}

func (self *digest32) Size() int {
	return 4
}

func (self *digest32) BlockSize() int {
	return 1
}

func (self *digest32) Reset() {
	self.len = 0
	self.tail = nil
	self.hash = self.seed
}

func (self *digest32) Write(b []byte) (n int, err error) {
	if len(b) <= 0 {
		return 0, nil
	}

	if self.tail != nil {
		if len(self.tail)+len(b) <= 3 {
			self.tail = append(self.tail, b...)
			self.len += len(b)
		} else {
			r := 4 - len(self.tail)
			bb := append(self.tail, b[:r]...)

			self.len -= len(self.tail)
			self.tail = nil

			self.Write(bb)
			self.Write(b[r:])
		}
		return len(b), nil
	}

	if len(b) <= 3 {
		self.tail = make([]byte, len(b))
		copy(self.tail, b)
		self.len += len(b)
		return len(b), nil
	}

	nblocks := len(b) >> 2
	h := self.hash
	p := unsafe.Pointer(&b[0])
	for i := 0; i < nblocks; i++ {
		k := *(*uint32)(p)

		k *= c1_32
		k = rol32(k, r1_32)
		k *= c2_32

		h ^= k
		h = rol32(h, r2_32)
		h = h*m_32 + c5_32

		p = unsafe.Pointer(uintptr(p) + 4)
	}

	if len(b)&3 != 0 {
		self.tail = make([]byte, len(b)&3)
		copy(self.tail, b[nblocks*4:])
	}

	self.hash = h
	self.len += len(b)

	return len(b), nil
}

func (self *digest32) Sum(b []byte) []byte {
	h := self.Sum32()
	return append(b, byte(h>>24), byte(h>>16), byte(h>>8), byte(h))
}

func (self *digest32) Sum32() uint32 {
	h := self.hash

	if self.tail != nil {
		k := uint32(0)
		for i, b := range self.tail {
			k |= uint32(b) << uint(i*8)
		}

		k *= c1_32
		k = rol32(k, r1_32)
		k *= c2_32

		h ^= k
	}

	h ^= uint32(self.len)
	h ^= (h >> 16)
	h *= c3_32
	h ^= (h >> 13)
	h *= c4_32
	h ^= (h >> 16)
	return h
}
