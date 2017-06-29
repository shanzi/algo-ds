package hamt

import (
	"fmt"
	"testing"
)

func TestPutAndGet(t *testing.T) {
	tmap := &tMapHead{0, 0, nil}

	for i := 0; i < 2048; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		if !tmap.Put(key, value) {
			t.Errorf("Should put %s successfully\n", key)
		}
		if tmap.Size() != i+1 {
			t.Errorf("Expect size to be %d, but got %d\n", i+1, tmap.Size())
		}
	}

	for i := 0; i < 2048; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)

		if v, ok := tmap.Get(key); !ok {
			t.Errorf("Should get %s successfully\n", key)
		} else if v != value {
			t.Errorf("Expect value to be %s, but got %s\n", value, v)
		}
	}
}

func TestGetNonExistKey(t *testing.T) {
	tmap := &tMapHead{0, 0, nil}

	// Empty map
	for i := 0; i < 2048; i++ {
		key := fmt.Sprintf("key-%d", i)

		if v, ok := tmap.Get(key); ok || v != nil {
			t.Errorf("Should get nil for key %s\n", key)
		}
	}

	// Fill in some items
	for i := 0; i < 2048; i += 2 {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		if !tmap.Put(key, value) {
			t.Errorf("Should put %s successfully\n", key)
		}
	}

	for i := 0; i < 2048; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)

		v, ok := tmap.Get(key)

		if (i&1) == 1 && ok {
			t.Errorf("Should get nil for key %s\n", key)
		}

		if (i & 1) == 0 {
			if !ok {
				t.Errorf("Should get %s successfully\n", key)
			}

			if v != value {
				t.Errorf("Expect value to be %s, but got %s\n", value, v)
			}
		}
	}
}

func TestPutSameKeys(t *testing.T) {
	tmap := &tMapHead{0, 0, nil}
	count := 0

	for i := 0; i < 2048; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		if !tmap.Put(key, value) {
			t.Errorf("Should put %s successfully\n", key)
		}
		count++
	}

	for i := 0; i < 2048; i++ {
		key := fmt.Sprintf("key-%d", i)
		if (i & 1) == 0 {
			// For even numbers, put original object
			value := fmt.Sprintf("value-%d", i)
			if tmap.Put(key, value) {
				t.Errorf("Key %s with value %s should not be put ok twice\n", key, value)
			}
		} else {
			// For odd numbers put different object
			value := fmt.Sprintf("changed-value-%d", i)
			if !tmap.Put(key, value) {
				t.Errorf("Key %s with value %s should be put ok\n", key, value)
			}
			count++
		}

		if tmap.Size() != count {
			t.Errorf("Expect size to be %d, but got %d\n", count, tmap.Size())
		}
	}
}

func BenchmarkGetMethod(b *testing.B) {
	tmap := &tMapHead{0, 0, nil}
	keys := make([]string, 4096)
	value := "samevalue"

	for i := 0; i < 4096; i++ {
		key := fmt.Sprintf("key-%d", i)
		keys[i] = key
		tmap.Put(key, value)
	}

	for i := 0; i < b.N; i++ {
		key := keys[i%4096]
		tmap.Get(key)
	}
}

func BenchmarkPutMethod(b *testing.B) {
	tmap := &tMapHead{0, 0, nil}
	keys := make([]string, (1 << 20))
	value := "samevalue"

	for i := 0; i < (1 << 16); i++ {
		keys[i] = fmt.Sprintf("key-%d", i)
	}

	for i := 0; i < b.N; i++ {
		key := keys[i&((1<<16)-1)]
		tmap.Put(key, value)
	}
}
