package hamt

type Map interface {
	Get(key string) (interface{}, bool)
	Put(key string, value interface{}) (Map, bool)
	Remove(key string) (Map, interface{})
	TransientMap() TransientMap
	Size() int
}
