package suffix_automaton

type suffix_automaton struct {
	length      int
	states      []State
	constructor StateConstructor
}

func New(length int, constructor StateConstructor) *suffix_automaton {
	sam := &suffix_automaton{
		0,
		make([]State, (length<<1)+5),
		constructor,
	}

	init_state := constructor.NewState()
	sam.appendState(init_state)
	return sam
}

func (self *suffix_automaton) appendState(s State) {
	self.states[self.length] = s
	self.length += 1
}

func (self *suffix_automaton) Add(ch int) {
	last := self.states[self.length-1]
	next := self.constructor.NewState()

	next.SetMaxLength(last.MaxLength() + 1)

	for last != nil && last.TransitionAt(ch) == nil {
		last.SetTransitionAt(ch, next)
		last = last.Suffix()
	}

	if last == nil {
		next.SetSuffix(self.states[0])
		self.appendState(next)
		return
	}

	lastNext := last.TransitionAt(ch)
	if lastNext.MaxLength() == last.MaxLength()+1 {
		next.SetSuffix(lastNext)
		return
	}

	splitLastNext := self.constructor.CloneState(lastNext)

	lastNext.SetSuffix(splitLastNext)
	next.SetSuffix(splitLastNext)

	for last != nil && last.TransitionAt(ch) == lastNext {
		last.SetTransitionAt(ch, splitLastNext)
		last = last.Suffix()
	}

	splitLastNext.SetMinLength(splitLastNext.Suffix().MaxLength() + 1)

	self.appendState(splitLastNext)
	self.appendState(next)
}
