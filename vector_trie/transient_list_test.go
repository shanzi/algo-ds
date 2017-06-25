package vector_trie

import "testing"

func generateList(start, end, step int) TransientList {
	list := New().TransientList()
	for i := start; i < end; i += step {
		list.PushBack(i)
	}
	return list
}

func TestTransientGet(t *testing.T) {
	list := generateList(0, 100, 1)

	for i := 0; i < 100; i++ {
		if v, _ := list.Get(i); v != i {
			t.Errorf("Expect value at index %d to be %d, but got %d\n", i, i, v)
		}
	}
}

func TestTransientSet(t *testing.T) {
	list := generateList(0, 100, 1)

	for i := 0; i < 100; i++ {
		list.Set(i, i*2)
	}

	for i := 0; i < 100; i++ {
		if v, _ := list.Get(i); v != i*2 {
			t.Errorf("Expect value at index %d to be %d, but got %d\n", i, i*2, v)
		}
	}
}

func TestTransientPushBack(t *testing.T) {
	list := New().TransientList()

	if l := list.Len(); l != 0 {
		t.Errorf("Expect length of list to be %d, but got %d", 0, l)
	}

	for i := 0; i < 100; i++ {
		list.PushBack(i)
		if l := list.Len(); l != i+1 {
			t.Errorf("Expect length of list to be %d, but got %d", i+1, l)
		}
	}
}

func TestTransientRemoveBack(t *testing.T) {
	list := generateList(0, 100, 1)

	if l := list.Len(); l != 100 {
		t.Errorf("Expect length of list to be %d, but got %d", 100, l)
	}

	for i := 99; i >= 0; i-- {
		if v := list.RemoveBack(); v != i {
			t.Errorf("Expect %d to be removed from back, but got %d", i, v)
		}

		if l := list.Len(); l != i {
			t.Errorf("Expect length of list to be %d, but got %d", i, l)
		}
	}
}

func TestTransientLongList(t *testing.T) {
	list := New().TransientList()
	for i := 0; i < (1 << 20); i++ {
		list.PushBack(i)

		if l := list.Len(); l != i+1 {
			t.Errorf("Expect length of list to be %d, but got %d", i+1, l)
		}
	}

	for i := 0; i < (1 << 20); i++ {
		if v, _ := list.Get(i); v != i {
			t.Errorf("Expect value at index %d to be %d, but got %d\n", i, i, v)
		}

		list.Set(i, -i)

		if v, _ := list.Get(i); v != -i {
			t.Errorf("Expect value at index %d to be %d, but got %d\n", i, -i, v)
		}
	}

	for i := (1 << 20) - 1; i >= 0; i-- {
		if v := list.RemoveBack(); v != -i {
			t.Errorf("Expect %d to be removed from back, but got %d", -i, v)
		}

		if l := list.Len(); l != i {
			t.Errorf("Expect length of list to be %d, but got %d", i, l)
		}
	}
}

func TestTransientGetIndexOutOfBound(t *testing.T) {
	list := generateList(0, 100, 1)
	_, ok := list.Get(-1)
	if ok {
		t.Error("Expect ok to be false when getting negative index")
	}
	_, ok = list.Get(1000)
	if ok {
		t.Error("Expect ok to be false when getting index beyound size")
	}
}

func TestTransientSetNegativeIndexPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != "Index out of bound" {
			t.Errorf("Expect panic with `Index out of bound`")
		}
	}()

	list := generateList(0, 100, 1)
	list.Set(-1, 123)
}

func TestTransientSetOverIndexPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != "Index out of bound" {
			t.Errorf("Expect panic with `Index out of bound`")
		}
	}()

	list := generateList(0, 100, 1)
	list.Set(100, -1)
}

func TestTransientEmptyRemovePanic(t *testing.T) {
	defer func() {
		if r := recover(); r != "Remove from empty list" {
			t.Errorf("Expect panic with `Remove from empty list`")
		}
	}()

	list := New().TransientList()
	list.RemoveBack()
}

func BenchmarkPushBack(b *testing.B) {
	list := New().TransientList()
	var v interface{}
	for i := 0; i < b.N; i++ {
		list.PushBack(v)
	}
}

func BenchmarkGet(b *testing.B) {
	mask := (1 << 10) - 1
	list := generateList(0, mask+1, 1)

	for i := 0; i < b.N; i++ {
		list.Get(i & mask)
	}
}

func BenchmarkSet(b *testing.B) {
	mask := (1 << 10) - 1
	list := generateList(0, mask+1, 1)

	var v interface{} = 0

	for i := 0; i < b.N; i++ {
		list.Set(i&mask, v)
	}
}

func BenchmarkRemoveBack(b *testing.B) {
	mask := (1 << 20) - 1
	list := generateList(0, mask+1, 1)

	for i := 0; i < b.N; i++ {
		if list.Len() > 0 {
			list.RemoveBack()
		} else {
			b.StopTimer()
			list = generateList(0, mask+1, 1)
			b.StartTimer()
		}
	}
}
