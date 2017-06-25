// Package list provides an implementation of immutable list based on vector trie
package vector_trie

type List interface {
	Get(n int) (interface{}, bool)
	Set(n int, value interface{}) (List, bool)
	PushBack(value interface{}) List
	RemoveBack() (List, interface{})
	TransientList() TransientList
	Len() int
}

type listHead tListHead

var empty = &listHead{0, 0, 0, 0, nil, nil}

func New() List {
	return empty
}

func (head *listHead) Get(n int) (interface{}, bool) {
	return (*tListHead)(head).Get(n)
}

func (head *listHead) Set(n int, value interface{}) (List, bool) {
	t := head.TransientList()
	if t.Set(n, value) {
		return t.Persist(), true
	} else {
		return head, false
	}
}

func (head *listHead) PushBack(value interface{}) List {
	t := head.TransientList()
	t.PushBack(value)
	return t.Persist()
}

func (head *listHead) RemoveBack() (List, interface{}) {
	if head.len == 1 {
		value, _ := head.Get(0)
		return empty, value
	} else {
		t := head.TransientList()
		value := t.RemoveBack()
		return t.Persist(), value
	}
}

func (head *listHead) TransientList() TransientList {
	id := nextId()
	return &tListHead{id, head.len, head.level, head.offset, head.root, head.tail}
}

func (head *listHead) Len() int {
	return head.len
}
