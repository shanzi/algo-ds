package vector_trie

const (
	SHIFT     = 5
	NODE_SIZE = (1 << SHIFT)
	MASK      = NODE_SIZE - 1
)
