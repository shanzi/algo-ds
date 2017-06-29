package hamt

import "sync/atomic"

var id uint64 = 0

func nextId() uint64 {
	return atomic.AddUint64(&id, 1)
}
