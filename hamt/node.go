package hamt

type node struct {
	id       uint64
	mask     uint32
	children []interface{}
}

func newNode(id uint64, size int) *node {
	return &node{id, 0, make([]*node, size)}
}

func (self *node) has(index int) bool {
	if index >= 0 && index < 31 {
		return ((self.mask >> index) & 1) == 1
	}
	return false
}

func (self *node) childAt(index int) interface{} {
	if self.has(index) {
		m := self.mask & uint32((i<<index)-1)
		pos := popcount32(m)

		return self.children[pos]
	}
	return nil
}

func (self *node) putChildAt(id uint64, index int, child interface{}) {

}
