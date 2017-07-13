package suffix_automaton

type fixedsizestate struct {
	max_len       int
	min_len       int
	trans         []State
	suffixed_link State
}

func (self *fixedsizestate) MaxLength() int {
	return self.max_len
}

func (self *fixedsizestate) SetMaxLength(len int) {
	self.max_len = len
}

func (self *fixedsizestate) MinLength() int {
	return self.min_len
}

func (self *fixedsizestate) SetMinLength(len int) {
	self.min_len = len
}

func (self *fixedsizestate) TransitionAt(charpoint int) State {
	if charpoint >= 0 && charpoint < len(self.trans) {
		return self.trans[charpoint]
	}
	return nil
}

func (self *fixedsizestate) TransitionCharPoints() []int {
	ret := []int{}
	for i := 0; i < len(self.trans); i++ {
		if self.trans[i] != nil {
			ret = append(ret, i)
		}
	}
	return ret
}

func (self *fixedsizestate) SetTransitionAt(charpoint int, s State) bool {
	if charpoint >= 0 && charpoint < len(self.trans) {
		self.trans[charpoint] = s
		return true
	}
	return false
}

func (self *fixedsizestate) Suffix() State {
	return self.suffixed_link
}

func (self *fixedsizestate) SetSuffix(s State) {
	if s == nil {
		self.suffixed_link = nil
		self.min_len = 0
	} else {
		self.suffixed_link = s
		self.min_len = s.MaxLength() + 1
	}
}

type FixedSizeStateConstructor int

func (self FixedSizeStateConstructor) NewState() State {
	size := int(self)
	return &fixedsizestate{0, 0, make([]State, size), nil}
}

func (self FixedSizeStateConstructor) CloneState(s State) State {
	size := int(self)
	new := &fixedsizestate{
		s.MaxLength(),
		s.MinLength(),
		make([]State, size),
		s.Suffix(),
	}

	for i := 0; i < size; i++ {
		new.SetTransitionAt(i, s.TransitionAt(i))
	}
	return new
}
