package murmur3

import (
	"testing"
	"unsafe"
)

var test_cases = []string{
	"",
	"hello",
	"hello, world",
	"19 Jan 2038 at 3:14:07 AM",
	"The quick brown fox jumps over the lazy dog.",
}

var test_seeds = []uint32{
	0x2bfab2f3,
	0xabeabfa7,
	0xffabbcc3,
}

func murmurSum32(data []byte, seed uint32) uint32 {

	var h1 = seed

	nblocks := len(data) / 4
	var p uintptr
	if len(data) > 0 {
		p = uintptr(unsafe.Pointer(&data[0]))
	}
	p1 := p + uintptr(4*nblocks)
	for ; p < p1; p += 4 {
		k1 := *(*uint32)(unsafe.Pointer(p))

		k1 *= c1_32
		k1 = (k1 << 15) | (k1 >> 17) // rotl32(k1, 15)
		k1 *= c2_32

		h1 ^= k1
		h1 = (h1 << 13) | (h1 >> 19) // rotl32(h1, 13)
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
		k1 *= c1_32
		k1 = (k1 << 15) | (k1 >> 17) // rotl32(k1, 15)
		k1 *= c2_32
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

func TestMurmur32(t *testing.T) {
	digest := New32(0)

	for _, s := range test_cases {
		data := []byte(s)

		digest.Reset()
		digest.Write(data)

		hash := digest.Sum32()
		chash := murmurSum32(data, 0)

		if hash != chash {
			t.Errorf("Expect hash value to be %x, got %x", chash, hash)
		}
	}
}

func TestMurmur32Partial(t *testing.T) {
	digest := New32(0)

	data := []byte{}
	for _, s := range test_cases {
		pdata := []byte(s)
		data = append(data, pdata...)

		digest.Write(pdata)

		hash := digest.Sum32()
		chash := murmurSum32(data, 0)

		if hash != chash {
			t.Errorf("Expect hash value to be %x, got %x", chash, hash)
		}
	}
}

func TestMurmur32Seed(t *testing.T) {
	for _, seed := range test_seeds {
		digest := New32(seed)

		for _, s := range test_cases {
			data := []byte(s)

			digest.Reset()
			digest.Write(data)

			hash := digest.Sum32()
			chash := murmurSum32(data, seed)

			if hash != chash {
				t.Errorf("Expect hash value to be %x, got %x", chash, hash)
			}
		}
	}
}

func TestMurmur32SeedPartial(t *testing.T) {
	for _, seed := range test_seeds {
		digest := New32(seed)

		data := []byte{}
		for _, s := range test_cases {
			pdata := []byte(s)
			data = append(data, pdata...)

			digest.Write(pdata)

			hash := digest.Sum32()
			chash := murmurSum32(data, seed)

			if hash != chash {
				t.Errorf("Expect hash value to be %x, got %x", chash, hash)
			}
		}
	}
}
