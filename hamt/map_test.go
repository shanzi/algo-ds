package hamt

import (
	"fmt"
	"testing"
)

func TestMapPut(t *testing.T) {
	m := New()

	for i := 0; i < 1024; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)

		n := m.Put(key, value)

		if v, ok := n.Get(key); !ok || v != value {
			t.Errorf("Expect get %s for key %s, but got %s\n", value, key, v)
		}

		if n.Size() != i+1 {
			t.Errorf("Expect size of new map to be %d, but got %d\n", i+1, n.Size())
		}

		if _, ok := m.Get(key); ok {
			t.Errorf("Should not be able to get key %s in old map", key)
		}

		if m.Size() != i {
			t.Errorf("Expect size of old map to be %d, but got %d\n", i, m.Size())
		}

		m = n
	}
}

func TestMapPutOverwrite(t *testing.T) {
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

	if pm.Size() != 1024 {
		t.Errorf("Expect size of pn to be %d, but got %d\n", 1024, pm.Size())
	}

	if pn.Size() != 1024 {
		t.Errorf("Expect size of pn to be %d, but got %d\n", 1024, pn.Size())
	}

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

func TestMapRemove(t *testing.T) {
	m := New().TransientMap()
	for i := 0; i < 1024; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		m.Put(key, value)
	}

	pm := m.Map()

	n := pm.TransientMap()
	for i := 0; i < 1024; i += 2 {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		if v, ok := n.Remove(key); !ok || v != value {
			t.Errorf("Should be able to remove key %s, expect value %s, got %s\n", key, value, v)
		}
	}

	pn := n.Map()

	if pm.Size() != 1024 {
		t.Errorf("Expect size of pm to be %d, got %d\n", 1024, pm.Size())
	}

	if pn.Size() != 512 {
		t.Errorf("Expect size of pm to be %d, got %d\n", 512, pn.Size())
	}

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

		_, ok := pn.Get(key)
		if (i&1) == 1 && !ok {
			t.Errorf("Should be able to got key %s in new map\n", key)
		}

		if (i&1) == 0 && ok {
			t.Errorf("Should not be able to got key %s in new map\n", key)
		}
	}
}
