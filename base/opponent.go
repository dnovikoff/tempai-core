package base

type Opponent int

//go:generate stringer -type=Opponent
const (
	Self Opponent = iota
	Right
	Front
	Left
)
