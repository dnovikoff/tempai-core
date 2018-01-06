package yaku

type HanPoints int
type FuPoints int
type FuPointsRounded int

func (this FuPoints) Round() FuPointsRounded {
	if this == 25 {
		return FuPointsRounded(this)
	}
	left := this % 10
	if left == 0 {
		return FuPointsRounded(this)
	}
	return FuPointsRounded(this - left + 10)
}
