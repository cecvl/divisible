package game

type State int

const (
	StateQuestion State = iota
	StateBonus
	StateFinished
)
