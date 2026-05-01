package game

type State int

const (
	StateIntro State = iota
	StateQuestion
	StateBonus
	StatePaused
	StateFinished
)
