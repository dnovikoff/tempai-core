package yaku

type HanPoints int
type FuPoints int
type FuPointsRounded int

func (p FuPoints) Round() FuPointsRounded {
	if p == 25 {
		return FuPointsRounded(p)
	}
	left := p % 10
	if left == 0 {
		return FuPointsRounded(p)
	}
	return FuPointsRounded(p - left + 10)
}
