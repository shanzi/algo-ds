package hamt

type node struct {
	id       uint64
	mask     uint32
	children []interface{}
}

func newNode(id uint64, mask uint32, size int) *node {
	cap := size + (size & 1)
	return &node{id, mask, make([](interface{}), size, cap)}
}

func (self *node) has(index uint) bool {
	if index >= 0 && index < 32 {
		return ((self.mask >> index) & 1) == 1
	}
	return false
}

func (self *node) pos(index uint) int {
	m := self.mask & uint32((1<<index)-1)
	return popcount32(m)
}

func (self *node) size() int {
	return len(self.children)
}

func (self *node) trim() {
	// This method removes unnecessary empty slots in children
	// once the difference between cap and len of child are larger than
	// desired as to save memory
	if cap(self.children)-len(self.children) > 4 {
		size := self.size()
		children := make([](interface{}), size, size+(size&1))
		copy(children, self.children)
		self.children = children
	}
}

func (self *node) childAt(index uint) interface{} {
	if self.has(index) {
		return self.children[self.pos(index)]
	}
	return nil
}

func (self *node) putChildAt(id uint64, index uint, child interface{}) *node {
	pos := self.pos(index)
	if id == self.id {
		// Within the same transaction
		if !self.has(index) {
			children := append(self.children, nil)
			for i := len(children) - 1; i > pos; i-- {
				children[i] = children[i-1]
			}
			self.children = children
		}
		self.children[pos] = child
		self.mask |= (1 << index)
		return self
	} else {
		// Not within the same transaction, copy the node
		cloned := newNode(id, self.mask|(1<<index), len(self.children)+1)
		for i := 0; i < pos; i++ {
			cloned.children[i] = self.children[i]
		}
		cloned.children[pos] = child
		for i := pos; i < len(self.children); i++ {
			cloned.children[i+1] = self.children[i]
		}
		return cloned
	}
}

func (self *node) removeChildAt(id uint64, index uint) *node {
	if !self.has(index) {
		return self
	}

	pos := self.pos(index)

	if id == self.id {
		// Within the same transaction
		children := self.children
		for i := pos + 1; i < len(children); i++ {
			children[i-1] = children[i]
		}
		// Clear reference to be gc friendly
		children[len(children)-1] = nil
		self.mask ^= (1 << index)
		self.children = children
		self.trim()
		return self
	} else {
		// Not within the same transaction, copy the node
		cloned := newNode(id, self.mask^(1<<index), len(self.children)-1)
		for i := 0; i < pos; i++ {
			cloned.children[i] = self.children[i]
		}

		for i := pos + 1; i < len(self.children); i++ {
			cloned.children[i-1] = self.children[i]
		}
		return cloned
	}
}
