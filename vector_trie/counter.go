// Package counter provides an atomic global counter for generating unique id
package vector_trie

import "sync/atomic"

var id uint64 = 0

func nextId() uint64 {
	return atomic.AddUint64(&id, 1)
}
