package aho_corasick

type Match interface {
	Start() int
	End() int
}

type match struct {
	start int
	end   int
}

func (self *match) Start() int {
	return self.start
}

func (self *match) End() int {
	return self.end
}
