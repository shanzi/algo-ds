package suffix_automaton

type State interface {
	MaxLength() int
	SetMaxLength(len int)

	MinLength() int
	SetMinLength(len int)

	TransitionAt(charpoint int) State
	TransitionCharPoints() []int
	SetTransitionAt(charpoint int, s State) bool

	Suffix() State
	SetSuffix(suffix State)
}

type StateConstructor interface {
	NewState() State
	CloneState(s State) State
}
