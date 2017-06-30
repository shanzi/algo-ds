package murmur3

import "testing"

var test_cases_128 = []string{
	"",
	"hello",
	"hello, world",
	"19 Jan 2038 at 3:14:07 AM",
	"The quick brown fox jumps over the lazy dog.",
}

var test_cases_128res = []uint64{
	0x0000000000000000, 0x0000000000000000,
	0xcbd8a7b341bd9b02, 0x5b1e906a48ae1d19,
	0x342fac623a5ebc8e, 0x4cdcbc079642414d,
	0xb89e5988b737affc, 0x664fc2950231b2cb,
	0xcd99481f9ee902c9, 0x695da1a38987b6e7,
}

func TestMurmur128(t *testing.T) {
	digest := New128()
	for i, v := range test_cases_128 {
		digest.Reset()
		digest.Write(([]byte)(v))

		h1, h2 := digest.Sum128()
		ch1 := test_cases_128res[i*2]
		ch2 := test_cases_128res[i*2+1]
		if !(h1 == ch1 && h2 == ch2) {
			t.Errorf("Expect hash of \"%s\" to be %x%x, but got %x%x", v, ch2, ch1, h2, h1)
		}
	}
}

func TestMurmur128Incremental(t *testing.T) {
	digest1 := New128()
	digest2 := New128()

	data := []byte{}
	for _, v := range test_cases_128 {
		digest2.Write(([]byte)(v))
		data = append(data, ([]byte)(v)...)
	}

	digest1.Write(data)
	ch1, ch2 := digest1.Sum128()
	h1, h2 := digest2.Sum128()

	if !(h1 == ch1 && h2 == ch2) {
		t.Errorf("Expect hash to be %x%x, but got %x%x", ch2, ch1, h2, h1)
	}
}
