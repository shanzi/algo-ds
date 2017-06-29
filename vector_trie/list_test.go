package vector_trie

import "testing"

func TestListGet(t *testing.T) {
	list := generateList(0, 1000, 1).List()

	for i := 0; i < 1000; i++ {
		if n, _ := list.Get(i); n != i {
			t.Error("Expect value at index %d to be %d, but got %d\n", i, i, n)
		}
	}
}

func TestListSet(t *testing.T) {
	list := generateList(0, 1000, 1).List()

	for i := 0; i < 1000; i++ {
		nlist, ok := list.Set(i, i*2+1)
		if !ok {
			t.Errorf("Setting value at index %d should be ok\n", i)
		}

		if n, _ := nlist.Get(i); n != i*2+1 {
			t.Errorf("Expect value at index %d to be %d, but got %d\n", i, i*2+1, n)
		}

		if n, _ := list.Get(i); n != i {
			t.Errorf("Expect value of old list at index %d to be %d, but got %d\n", i, i, n)
		}
	}

	for i := 0; i < 1000; i++ {
		nlist, ok := list.Set(i, i)
		if ok {
			t.Errorf("Setting same value at index %d should not be ok\n", i)
		}
		if list != nlist {
			t.Error("The new list should be equal to the old")
		}
	}
}

func TestListPushBack(t *testing.T) {
	oldlist := New()
	for i := 0; i < 1000; i++ {
		newlist := oldlist.PushBack(i)

		if newlist.Len() != i+1 {
			t.Errorf("Expect new list to have length %d\n", i+1)
		}

		if n, _ := newlist.Get(i); n != i {
			t.Errorf("Expect value at index %d to be %d, but got %d\n", i, i*2, n)
		}

		if oldlist.Len() != i {
			t.Errorf("Expect old list to have length %d\n", i)
		}

		oldlist = newlist
	}
}

func TestListRemoveBack(t *testing.T) {
	oldlist := generateList(0, 1000, 1).List()
	for i := 999; i >= 0; i-- {
		newlist, value := oldlist.RemoveBack()

		if newlist.Len() != i {
			t.Errorf("Expect new list to have length %d\n", i)
		}

		if value != i {
			t.Errorf("Expect value at index %d to be %d, but got %d\n", i, i, value)
		}

		if oldlist.Len() != i+1 {
			t.Errorf("Expect old list to have length %d\n", i+1)
		}

		if n, _ := oldlist.Get(i); n != i {
			t.Errorf("Expect value of old list at index %d to be %d, but got %d\n", i, i, n)
		}

		oldlist = newlist
	}

	if oldlist != New() {
		t.Error("List of length zero should be equal to Empty List")
	}
}
