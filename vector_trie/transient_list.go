package vector_trie

type TransientList interface {
	Get(n int) (interface{}, bool)
	Set(n int, value interface{}) bool
	PushBack(value interface{})
	RemoveBack() interface{}
	Persist() List
	Len() int
}

type tListHead struct {
	id     uint64
	len    int
	level  int
	offset int
	root   *trieNode
	tail   *trieNode
}

func (head *tListHead) Get(n int) (interface{}, bool) {
	if n < 0 || n >= head.len {
		return nil, false
	}

	if n >= head.offset {
		return head.tail.getChildValue(n - head.offset), true
	}

	root := head.root
	for lv := head.level - 1; ; lv-- {
		index := (n >> uint(lv*SHIFT)) & MASK
		if lv <= 0 {
			// Arrived at leaves node, return value
			return root.getChildValue(index), true
		} else {
			// Update root node
			root = root.getChildNode(index)
		}
	}
}

func (head *tListHead) Set(n int, value interface{}) bool {
	if n < 0 || n >= head.len {
		panic("Index out of bound")
	}

	ok := false

	if n >= head.offset {
		head.tail, ok = setTail(head.id, head.tail, n-head.offset, value)
	} else {
		head.root, ok = setInNode(head.id, head.root, n, head.level, value)
	}
	return ok
}

func (head *tListHead) PushBack(value interface{}) {
	// Increase the depth of tree while the capacity is not enough
	if head.len-head.offset < NODE_SIZE {
		// Tail node has free space
		head.tail, _ = setTail(head.id, head.tail, head.len-head.offset, value)
	} else {
		// Tail node is full
		n := head.offset
		lv := head.level
		root := head.root

		for lv == 0 || (n>>uint(lv*SHIFT)) > 0 {
			parent := newTrieNode(head.id)
			parent.children[0] = root
			root = parent
			lv++
		}

		head.root = putTail(head.id, root, head.tail, n, lv)
		head.tail, _ = setTail(head.id, nil, 0, value)

		head.level = lv
		head.offset += NODE_SIZE
	}
	head.len++
}

func (head *tListHead) RemoveBack() interface{} {
	if head.len == 0 {
		panic("Remove from empty list")
	}

	value := head.tail.getChildValue(head.len - head.offset - 1)
	head.tail, _ = setTail(head.id, head.tail, head.len-head.offset-1, nil) // clear reference to release memory

	head.len--

	if head.len == 0 {
		head.level = 0
		head.offset = 0
		head.root = nil
		head.tail = nil
	} else {
		if head.len <= head.offset {
			// tail is empty, retrieve new tail from root
			head.root, head.tail = getTail(head.id, head.root, head.len-1, head.level)
			head.offset -= NODE_SIZE
		}

		// Reduce the depth of tree if root only have one child
		n := head.offset - 1
		lv := head.level
		root := head.root

		for lv > 1 && (n>>uint((lv-1)*SHIFT)) == 0 {
			root = root.getChildNode(0)
			lv--
		}

		head.root = root
		head.level = lv
	}

	return value
}

func (head *tListHead) Persist() List {
	perisitHead := (*listHead)(head)
	perisitHead.id = 0
	return perisitHead
}

func (head *tListHead) Len() int {
	return head.len
}

func setInNode(id uint64, root *trieNode, n int, level int, value interface{}) (*trieNode, bool) {
	index := (n >> uint((level-1)*SHIFT)) & MASK

	if level == 1 {
		if root.getChildValue(index) == value {
			return root, false
		}
		return root.setChild(id, index, value), true
	} else {
		child := root.getChildNode(index)
		newChild, ok := setInNode(id, child, n, level-1, value)
		if ok {
			return root.setChild(id, index, newChild), true
		}
		return root, false
	}
}

func setTail(id uint64, tail *trieNode, n int, value interface{}) (*trieNode, bool) {
	if tail == nil {
		tail = newTrieNode(id)
	}

	if tail.getChildValue(n) == value {
		return tail, false
	}
	return tail.setChild(id, n, value), true
}

func getTail(id uint64, root *trieNode, n int, level int) (*trieNode, *trieNode) {
	index := (n >> uint((level-1)*SHIFT)) & MASK

	if level == 1 {
		return nil, root
	} else {
		child, tail := getTail(id, root.getChildNode(index), n, level-1)

		if index == 0 && child == nil {
			// The first element has been removed, which means current node
			// becomes empty, remove current node by returning nil
			return nil, tail
		} else {
			// Current node is not empty
			root = root.setChild(id, index, child)
			return root, tail
		}
	}
}

func putTail(id uint64, root *trieNode, tail *trieNode, n int, level int) *trieNode {
	index := (n >> uint((level-1)*SHIFT)) & MASK

	if level == 1 {
		return tail
	} else {
		if root == nil {
			root = newTrieNode(id)
		}
		return root.setChild(id, index, putTail(id, root.getChildNode(index), tail, n, level-1))
	}
}
