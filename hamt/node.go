package hamt

type node struct {
	id       uint64
	mask     uint32
	children []interface{}
}

func newNode(id uint64, mask uint32, size int) *node {
	cap := size + 2
	if cap > 32 {
		cap = 32
	}
	return &node{id, mask, make([](interface{}), size, cap)}
}

func (self *node) has(index uint) bool {
	if index >= 0 && index < 32 {
		return ((self.mask >> index) & 1) == 1
	}
	return false
}

func (self *node) childAt(index uint) interface{} {
	if self.has(index) {
		m := self.mask & uint32((1<<index)-1)
		pos := popcount32(m)

		return self.children[pos]
	}
	return nil
}

func (self *node) putChildAt(id uint64, index uint, child interface{}) *node {
	if id == self.id {
		// Within the same transaction
		m := self.mask & uint32((1<<index)-1)
		pos := popcount32(m)

		if !self.has(index) {
			children := append(self.children, nil)
			for i := len(children) - 1; i > pos; i-- {
				children[i] = children[i-1]
			}
		}
		self.children[pos] = child
		self.mask |= (1 << index)
		return self
	} else {
		// Not within the same transaction, copy the node
		cloned := newNode(id, self.mask, len(self.children))
		copy(cloned.children, self.children)
		return cloned.putChildAt(id, index, child)
	}
}
