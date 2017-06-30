package murmur3

import "testing"

var test_cases_32 = []string{
	"",
	"hello",
	"hello, world",
	"19 Jan 2038 at 3:14:07 AM",
	"The quick brown fox jumps over the lazy dog.",
}

func TestMurmur32(t *testing.T) {
	digest := New32()

	for _, s := range test_cases_32 {
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

func TestMurmur32Incremental(t *testing.T) {
	digest := New32()

	data := []byte{}
	for _, s := range test_cases_32 {
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
