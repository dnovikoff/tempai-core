package yaku

type Limit int

//go:generate stringer -type=Limit
// Limit numbers are now fixed and should not be changed
const (
	LimitNone      Limit = 0
	LimitMangan    Limit = 1
	LimitHaneman   Limit = 2
	LimitBaiman    Limit = 3
	LimitSanbaiman Limit = 4
	LimitYakuman   Limit = 5
)

func (lim Limit) ShortString() string {
	return lim.String()[len("Limit"):]
}

func (lim Limit) BaseHans() HanPoints {
	switch lim {
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
