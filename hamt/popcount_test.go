package hamt

import (
	"math/rand"
	"testing"
)

func popcount_iter(x uint32) int {
	count := 0
	for x != 0 {
		x &= x - 1
		count++
	}
	return count
}

func TestSmallIntegers(t *testing.T) {
	for i := 0; i < (1 << 12); i++ {
		a := popcount32(uint32(i))
		b := popcount_iter(uint32(i))
		if a != b {
			t.Errorf("Expect %d, got %d", b, a)
		}
	}
}

func TestRandomLargeIntegers(t *testing.T) {
	for i := 0; i < (1 << 12); i++ {
		x := rand.Uint32()
		a := popcount32(x)
		b := popcount_iter(x)
		if a != b {
			t.Errorf("Expect %d, got %d", b, a)
		}
	}
}
