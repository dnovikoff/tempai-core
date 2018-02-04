package base

type Opponent int

//go:generate stringer -type=Opponent
const (
	Self Opponent = iota
	Right
	Front
	Left
	OpponentCount
)

func (o Opponent) Advance(x int) Opponent {
	if x < 0 {
		x += (((-x) / 4) + 1) * 4
	}
	return (o + Opponent(x)) % OpponentCount
}

func (o Opponent) Next() Opponent {
	return o.Advance(int(Right))
}

func (o Opponent) Prev() Opponent {
	return o.Advance(int(Left))
}
