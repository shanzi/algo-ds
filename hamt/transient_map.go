package hamt

type TransientMap interface {
	Get(key string) (interface{}, bool)
	Put(key string, value interface{}) bool
	Remove(key string) interface{}
	Map() Map
	Size() int
}

type tMapHead struct {
	id   uint64
	size int
	root *node
}

func (self *tMapHead) Get(key string) (interface{}, bool) {
	if self.root == nil {
		return nil, false
	}

	if value, ok := self.getWithHash(self.root, key); ok {
		return value, true
	}

	return nil, false
}

func (self *tMapHead) Put(key string, value interface{}) bool {
	e := &entry{key: key, hash: 0, value: value}

	if self.root == nil {
		self.root = newNode(0, 0, 0)
	}

	if root, ok := self.putEntry(self.root, e, 0); ok {
		self.root = root
		self.size += 1
		return true
	}
	return false
}

func (self *tMapHead) getWithHash(root *node, key string) (interface{}, bool) {
	p := root
	var hash uint32
	for i := 0; i < 18; i++ {
		// At some specific depth, hash need to be recalculate
		switch i {
		case 0:
			hash = keyHash32(key, seed0)
		case 6:
			hash = keyHash32(key, seed1)
		case 12:
			hash = keyHash32(key, seed2)
		}

		d := uint(i % 6)
		h := uint((hash >> (d * 5)) & 0x1f)

		c := p.childAt(h)
		if c == nil {
			return nil, false
		}

		if e, ok := c.(*entry); ok {
			if e.key == key {
				// Find object
				return e.value, true
			} else {
				// No match, return nil
				return nil, false
			}
		}

		// c must be a node
		p = c.(*node)
	}

	// Nothing found after drained hash code,
	// return false indicating not found
	return nil, false
}

func (self *tMapHead) putEntry(root *node, e *entry, depth int) (*node, bool) {
	// At some specific depth, hash need to be recalculate
	switch depth {
	case 0:
		e.hash = keyHash32(e.key, seed0)
	case 6:
		e.hash = keyHash32(e.key, seed1)
	case 12:
		e.hash = keyHash32(e.key, seed2)
	case 18:
		panic("Inresolvable hash collision!")
	}

	d := uint(depth % 6)
	h := uint((e.hash >> (d * 5)) & 0x1f)

	if !root.has(h) {
		// Found a position to put new item in
		return root.putChildAt(self.id, h, e), true
	}

	child := root.childAt(h)
	if subnode, ok := child.(*node); ok {
		// Found a sub node, recursively put entry
		if child, ok = self.putEntry(subnode, e, depth+1); ok {
			return root.putChildAt(self.id, h, child), true
		} else {
			return root, false
		}
	}

	if olde, ok := child.(*entry); ok {
		// Found an entry
		if olde.key != e.key {
			// Collision. Create a new node and rehash current entry
			subnode := newNode(self.id, 0, 0)
			subnode, _ = self.putEntry(subnode, olde, depth+1)

			if child, ok = self.putEntry(subnode, e, depth+1); ok {
				return root.putChildAt(self.id, h, child), true
			} else {
				return root, false
			}
		}

		// Two keys are the same
		if olde.value == e.value {
			// Two values are the same, do nothing and return
			return root, false
		} else {
			// Two values are different, overwrite value
			return root.putChildAt(self.id, h, e), true
		}
	}

	assert_unreachable()
	return nil, false
}
