package hamt

type Entry interface {
	Key() string
	Value() interface{}
}

type entry struct {
	key   string
	hash  uint32
	value interface{}
}

func (self *entry) Key() string {
	return self.key
}

func (self *entry) Value() interface{} {
	return self.value
}
