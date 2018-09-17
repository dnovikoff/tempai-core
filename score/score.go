package score

import (
	"math"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/yaku"
)

type Honba int
type MoneyBase int
type RiichiSticks int

func (rs RiichiSticks) Money() Money {
	return Money(rs) * 1000
}

type Score struct {
	PayRon         Money
	PayRonDealer   Money
	PayTsumo       Money
	PayTsumoDealer Money
	Special        yaku.Limit

	Han yaku.HanPoints
	Fu  yaku.FuPoints
}

type ScoreChanges map[base.Wind]Money

func MergeChanges(args ...ScoreChanges) ScoreChanges {
	ret := make(ScoreChanges, 4)
	for _, arg := range args {
		for k, v := range arg {
			ret[k] += v
		}
	}
	return ret
}

func (sc ScoreChanges) ArrayFrom(wnd base.Wind, max int) []Money {
	ret := make([]Money, max)
	for k, v := range sc {
		ret[k.Advance(-int(wnd))] = v
	}
	return ret
}

func MergeChangeArrays(args ...[]Money) []Money {
	max := 0
	for _, arg := range args {
		if len(arg) > max {
			max = len(arg)
		}
	}
	ret := make([]Money, max)
	for _, arg := range args {
		for k, v := range arg {
			ret[k] += v
		}
	}
	return ret
}

func (sc ScoreChanges) TotalWin() Money {
	for _, v := range sc {
		if v > 0 {
			return v
		}
	}
	return 0
}

func (sc ScoreChanges) TotalPayed() (total Money) {
	for _, v := range sc {
		if v < 0 {
			total -= v
		}
	}
	return
}

func (s Score) GetChanges(selfWind, otherWind base.Wind, sticks RiichiSticks) ScoreChanges {
	changes := make(ScoreChanges, 4)
	isDealer := (selfWind == base.WindEast)
	isTsumo := (selfWind == otherWind)
	// Tsumo
	if isTsumo {
		if isDealer {
			for w := base.WindSouth; w != base.WindEnd; w++ {
				changes[w] = -s.PayTsumoDealer
			}
		} else {
			changes[base.WindEast] = -s.PayTsumoDealer
			for w := base.WindSouth; w != base.WindEnd; w++ {
				if w != selfWind {
					changes[w] = -s.PayTsumo
				}
			}
		}
	} else {
		if isDealer {
			changes[otherWind] = -s.PayRonDealer
		} else {
			changes[otherWind] = -s.PayRon
		}
	}
	changes[selfWind] = changes.TotalPayed() + sticks.Money()
	return changes
}

func (mb MoneyBase) Round100(mul int) Money {
	b := MoneyBase(mul) * mb
	left := b % 100
	if left > 0 {
		b = b - left + 100
	}
	return Money(b)
}

func (mb MoneyBase) Money(rules Rules, mul int, honba Honba) Money {
	return mb.Round100(mul) + GetHonbaMoney(rules, honba)
}

func calculateScoreBase(r Rules, han yaku.HanPoints, fu yaku.FuPoints) (MoneyBase, yaku.Limit) {
	switch {
	case han < 5:
		score := MoneyBase(fu.Round()) * MoneyBase(math.Pow(2.0, 2.0+float64(han)))
		if score > 2000 || (r.ManganRound() && score > 1900) {
			return 2000, yaku.LimitMangan
		}
		return score, yaku.LimitNone
	case han < 6:
		return 2000, yaku.LimitMangan
	case han < 8:
		return 3000, yaku.LimitHaneman
	case han < 11:
		return 4000, yaku.LimitBaiman
	case han < 13:
		return 6000, yaku.LimitSanbaiman
	case r.KazoeYakuman():
		return 8000, yaku.LimitYakuman
	}
	return 6000, yaku.LimitSanbaiman
}

// TODO: yakuman split for pao?
func GetScoreByResult(r Rules, res *yaku.Result, honba Honba) Score {
	if len(res.Yakuman) > 0 {
		mul := 1
		if r.YakumanSum() {
			mul = len(res.Yakuman)
		}
		if r.YakumanDouble() {
			for _, v := range res.Yakuman {
				if v > 1 {
					if r.YakumanSum() {
						mul++
					} else {
						mul = 2
						break
					}
				}
			}
		}
		return GetYakumanScore(r, mul, honba)
	}
	return GetScore(r, res.Sum(), res.Fus.Sum(), honba)
}

func GetScoreByBase(r Rules, base MoneyBase, special yaku.Limit, han yaku.HanPoints, fu yaku.FuPoints, honba Honba) (score Score) {
	score.Special = special
	score.PayTsumo = base.Money(r, 1, honba)
	score.PayTsumoDealer = base.Money(r, 2, honba)
	score.PayRon = base.Money(r, 4, honba*3)
	score.PayRonDealer = base.Money(r, 6, honba*3)
	score.Han = han
	score.Fu = fu
	return score
}

func GetYakumanScore(r Rules, mul int, honba Honba) Score {
	return GetScoreByBase(r, MoneyBase(mul)*8000, yaku.LimitYakuman, 0, 0, honba)
}

func GetScore(r Rules, han yaku.HanPoints, fu yaku.FuPoints, honba Honba) Score {
	base, special := calculateScoreBase(r, han, fu)
	return GetScoreByBase(r, base, special, han, fu, honba)
}
