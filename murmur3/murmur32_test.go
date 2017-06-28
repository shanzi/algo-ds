package murmur3

import "testing"

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

func TestMurmur32(t *testing.T) {
	digest := New32(0)

	for _, s := range test_cases {
		data := []byte(s)

		digest.Reset()
		digest.Write(data)

		hash := digest.Sum32()
		chash := Sum32(data, 0)

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
		chash := Sum32(data, 0)

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
			chash := Sum32(data, seed)

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
			chash := Sum32(data, seed)

			if hash != chash {
				t.Errorf("Expect hash value to be %x, got %x", chash, hash)
			}
		}
	}
}
