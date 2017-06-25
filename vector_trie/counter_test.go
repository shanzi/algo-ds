package vector_trie

import "testing"

func TestCounter(t *testing.T) {
	for i := 1; i <= 1000; i++ {
		if n := nextId(); n != uint64(i) {
			t.Errorf("Expect next id to be %d, but got %d", i, n)
		}
	}
}
