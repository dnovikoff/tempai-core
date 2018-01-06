package base

type RoundID int

const FirstRound RoundID = 0

func NewRound(r, w Wind) RoundID {
	return RoundID(r*4 + (w % 4))
}

func (this RoundID) Wind() Wind {
	return Wind(this % 4)
}

func (this RoundID) Round() Wind {
	return Wind(this / 4)
}

func (this RoundID) IsLastWind() bool {
	return this.Wind() == WindNorth
}

func (this RoundID) Next() RoundID {
	return this + 1
}
