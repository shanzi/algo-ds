package hamt

type Map interface {
	Get(key string) (interface{}, bool)
	Put(key string, value interface{}) Map
	Remove(key string) (Map, interface{})
	TransientMap() TransientMap
	Size() int
}

type mapHead tMapHead

var empty = &mapHead{0, 0, nil}

func New() Map {
	return empty
}

func (self *mapHead) Get(key string) (interface{}, bool) {
	return (*tMapHead)(self).Get(key)
}

func (self *mapHead) Put(key string, value interface{}) Map {
	t := self.TransientMap()
	if ok := t.Put(key, value); ok {
		return t.Map()
	} else {
		return self
	}
}

func (self *mapHead) Remove(key string) (Map, interface{}) {
	t := self.TransientMap()
	if value, ok := t.Remove(key); ok {
		return t.Map(), value
	} else {
		return self, nil
	}
}

func (self *mapHead) TransientMap() TransientMap {
	return &tMapHead{nextId(), self.size, self.root}
}

func (self *mapHead) Size() int {
	return self.size
}
