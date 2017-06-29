package hamt

import (
	"fmt"
	"testing"
)

func TestPut(t *testing.T) {
	m := New()

	for i := 0; i < 1024; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)

		n := m.Put(key, value)

		if v, ok := n.Get(key); !ok || v != value {
			t.Errorf("Expect get %s for key %s, but got %s\n", value, key, v)
		}

		if _, ok := m.Get(key); ok {
			t.Errorf("Should not be able to get key %s in old map", key)
		}

		m = n
	}
}

func TestPutOverwrite(t *testing.T) {
	m := New().TransientMap()
	for i := 0; i < 1024; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		m.Put(key, value)
	}

	pm := m.Map()

	n := pm.TransientMap()
	for i := 0; i < 1024; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("changed-%d", i)
		n.Put(key, value)
	}
	pn := n.Map()

	// old map should not changed
	for i := 0; i < 1024; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		if v, _ := pm.Get(key); v != value {
			t.Errorf("Expect old map not changed on key %s, expect value %s, got %s\n", key, value, v)
		}
	}

	// New map should be of new content
	for i := 0; i < 1024; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("changed-%d", i)
		if v, _ := pn.Get(key); v != value {
			t.Errorf("Expect new value %s, got %s\n", value, v)
		}
	}
}
