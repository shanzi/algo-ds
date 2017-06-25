package vector_trie

type trieNode struct {
	id       uint64
	children []interface{}
}

func newTrieNode(id uint64) *trieNode {
	return &trieNode{id, make([]interface{}, NODE_SIZE)}
}

func (node *trieNode) getChildNode(index int) *trieNode {
	if child := node.children[index]; child != nil {
		return child.(*trieNode)
	} else {
		return nil
	}
}

func (node *trieNode) getChildValue(index int) interface{} {
	return node.children[index]
}

func (node *trieNode) clone(id uint64) *trieNode {
	children := make([]interface{}, NODE_SIZE)
	copy(children, node.children)
	return &trieNode{id, children}
}

func (node *trieNode) setChild(id uint64, n int, child interface{}) *trieNode {
	if node.id == id {
		node.children[n] = child
		return node
	} else {
		newNode := node.clone(id)
		newNode.children[n] = child
		return newNode
	}
}
