package yaku

type Limit int

//go:generate stringer -type=Limit
const (
	LimitNone Limit = iota
	LimitMangan
	LimitHaneman
	LimitBaiman
	LimitSanbaiman
	LimitYakuman
)

func (this Limit) ShortString() string {
	return this.String()[len("Limit"):]
}

func (this Limit) BaseHans() HanPoints {
	switch this {
	case LimitMangan:
		return 5
	case LimitHaneman:
		return 6
	case LimitBaiman:
		return 8
	case LimitSanbaiman:
		return 11
	case LimitYakuman:
		return 13
	}
	return 0
}
