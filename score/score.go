package score

import (
	"math"

	"bitbucket.org/dnovikoff/tempai-core/base"
	"bitbucket.org/dnovikoff/tempai-core/yaku"
)

type Honba int
type MoneyBase int
type RiichiSticks int

func (this RiichiSticks) Money() Money {
	return Money(this) * 1000
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

func (this ScoreChanges) ArrayFrom(wnd base.Wind, max int) []Money {
	ret := make([]Money, max)
	for k, v := range this {
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

func (this ScoreChanges) TotalWin() Money {
	for _, v := range this {
		if v > 0 {
			return v
		}
	}
	return 0
}

func (this ScoreChanges) TotalPayed() (total Money) {
	for _, v := range this {
		if v < 0 {
			total -= v
		}
	}
	return
}

func (this Score) GetChanges(selfWind, otherWind base.Wind, sticks RiichiSticks) ScoreChanges {
	changes := make(ScoreChanges, 4)
	isDealer := (selfWind == base.WindEast)
	isTsumo := (selfWind == otherWind)
	// Tsumo
	if isTsumo {
		if isDealer {
			for w := base.WindSouth; w != base.WindEnd; w++ {
				changes[w] = -this.PayTsumoDealer
			}
		} else {
			changes[base.WindEast] = -this.PayTsumoDealer
			for w := base.WindSouth; w != base.WindEnd; w++ {
				if w != selfWind {
					changes[w] = -this.PayTsumo
				}
			}
		}
	} else {
		if isDealer {
			changes[otherWind] = -this.PayRonDealer
		} else {
			changes[otherWind] = -this.PayRon
		}
	}
	changes[selfWind] = changes.TotalPayed() + sticks.Money()
	return changes
}

func (this MoneyBase) Round100(mul int) Money {
	b := MoneyBase(mul) * this
	left := b % 100
	if left > 0 {
		b = b - left + 100
	}
	return Money(b)
}

func (this MoneyBase) Money(rules Rules, mul int, honba Honba) Money {
	return this.Round100(mul) + rules.GetHonbaMoney(honba)
}

func (this Rules) calculateScoreBase(han yaku.HanPoints, fu yaku.FuPoints) (MoneyBase, yaku.Limit) {
	switch {
	case han < 5:
		score := MoneyBase(fu.Round()) * MoneyBase(math.Pow(2.0, 2.0+float64(han)))
		if score > 2000 || (this.ManganRound && score > 1900) {
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
	case this.KazoeYakuman:
		return 8000, yaku.LimitYakuman
	}
	return 6000, yaku.LimitSanbaiman
}

// TODO: yakuman split for pao?
func (this Rules) GetScoreByResult(res *yaku.YakuResult, honba Honba) Score {
	if len(res.Yakuman) > 0 {
		mul := 1
		if this.YakumanSum {
			mul = len(res.Yakuman)
		}
		if this.DoubleYakuman {
			for _, v := range res.Yakuman {
				if v > 1 {
					if this.YakumanSum {
						mul++
					} else {
						mul = 2
						break
					}
				}
			}
		}
		return this.GetYakumanScore(mul, honba)
	}
	return this.GetScore(res.Sum(), res.Fus.Sum(), honba)
}

func (this Rules) GetScoreByBase(base MoneyBase, special yaku.Limit, han yaku.HanPoints, fu yaku.FuPoints, honba Honba) (score Score) {
	score.Special = special
	score.PayTsumo = base.Money(this, 1, honba)
	score.PayTsumoDealer = base.Money(this, 2, honba)
	score.PayRon = base.Money(this, 4, honba*3)
	score.PayRonDealer = base.Money(this, 6, honba*3)
	score.Han = han
	score.Fu = fu
	return score
}

func (this Rules) GetYakumanScore(mul int, honba Honba) (score Score) {
	return this.GetScoreByBase(MoneyBase(mul)*8000, yaku.LimitYakuman, 0, 0, honba)
}

func (this Rules) GetScore(han yaku.HanPoints, fu yaku.FuPoints, honba Honba) (score Score) {
	base, special := this.calculateScoreBase(han, fu)
	return this.GetScoreByBase(base, special, han, fu, honba)
}
