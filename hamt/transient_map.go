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
	var p interface{}

	// Round 0
	if p, ok = getWithHash(self.root, key, keyHash32(key, seed0)); ok || p == nil {
		return p, ok
	}

	// Hash0 drained, still have conflict, perform round 1
	if p, ok = getWithHash(p.(*node), key, keyHash32(key, seed1)); ok || p == nil {
		return p, ok
	}

	// Hash1 drained, still have conflict, perform round 2
	if p, ok = getWithHash(p.(*node), key, keyHash32(key, seed2)); ok || p == nil {
		return p, ok
	}

	// All inserted item must be within 3 rounds of getting.
	assert_unreachable()
}

func (self *tMapHead) Put(key string, value interface{}) bool {
	p := self.root
	e := &entry{key, value}
}

func (self *tMapHead) getWithHash(root *node, key string, hash uint32) (interface{}, bool) {
	p := root
	for i := 0; i < 32; i += 5 {
		h := (hash >> i) & 0x1f

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

	// Nothing found after drained hash code, return current node but
	// false indicating not found
	return p, false
}

func (self *tMapHead) putWithHash(root *node, e *entry, depth int) *node {
	d := depth & 0x1f
	h := (hash >> d) & 0x1f

	if !root.has(h) {
		// found a position to put new item in

	}
}
