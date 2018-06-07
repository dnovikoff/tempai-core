package tempai

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/calc"
	"github.com/dnovikoff/tempai-core/meld"
	"github.com/dnovikoff/tempai-core/tile"
)

func CalculatePairs(tiles compact.Instances) meld.Melds {
	var wait meld.Meld
	pairs := make(meld.Melds, 0, 7)
	tiles.Each(func(m compact.Mask) bool {
		c := m.Count()
		if c == 1 {
			wait = meld.NewTanki(m.First()).Meld()
		} else if c > 1 {
			pairs = append(pairs, meld.NewPairFromMask(m).Meld())
		} else {
			return false
		}
		return true
	})
	if wait.IsNull() || len(pairs) != 6 {
		return nil
	}
	return append(pairs, wait)
}

func CalculateKokushi(tiles compact.Instances) meld.Melds {
	if tiles.Count() != 13 {
		return nil
	}
	melds := make(meld.Melds, 0, 14)
	tanki := make(meld.Melds, 0, 14)

	var hole meld.Meld
	var pair meld.Meld
	for _, t := range tile.KokushiTiles {
		mask := tiles.GetMask(t)
		first := mask.First()
		switch mask.Count() {
		case 0:
			hole = meld.NewHole(t).Meld()
		case 1:
			tanki = append(tanki, meld.NewTanki(first).Meld())
			melds = append(melds, meld.NewOne(first).Meld())
		case 2:
			pair = meld.NewPairFromMask(mask).Meld()
		default:
			return nil
		}
	}

	if hole.IsNull() {
		if !pair.IsNull() || len(tanki) != 13 {
			return nil
		}
		return tanki
	}
	if pair.IsNull() || len(melds) != 11 {
		return nil
	}
	return append(melds, pair, hole)
}

func NewTempai(closed compact.Instances, opts *calc.Options) *calc.Calculator {
	if opts.Opened*3+closed.Count() != 13 {
		return nil
	}
	return calc.NewCalculator(meld.AllTempaiMelds, closed, opts)
}

func CalculateRegular(closed compact.Instances, options ...calc.Option) (ret TempaiMelds) {
	opts := calc.GetOptions(options...)
	return calculateRegular(closed, opts)
}

func calculateRegular(closed compact.Instances, opts *calc.Options) (ret TempaiMelds) {
	x := &result{declared: opts.Melds}
	opts.Results = x
	t := NewTempai(closed, opts)
	if t == nil {
		return nil
	}
	t.Calculate()
	return x.Melds
}

func Calculate(closed compact.Instances, options ...calc.Option) (ret TempaiMelds) {
	opts := calc.GetOptions(options...)
	return calculate(closed, opts)
}

func calculate(closed compact.Instances, opts *calc.Options) (ret TempaiMelds) {
	x := calculateRegular(closed, opts)
	if len(x) > 0 {
		ret = x
	}
	if len(opts.Melds) > 0 {
		return
	}
	// Pairs should be calcluated after regular hand
	// because of ryanpeiko
	melds := CalculatePairs(closed)
	if len(melds) > 0 {
		return append(ret, melds)
	}
	melds = CalculateKokushi(closed)
	if len(melds) > 0 {
		return TempaiMelds{melds}
	}
	return
}

func CheckTempai(closed compact.Instances, options ...calc.Option) bool {
	opts := calc.GetOptions(options...)
	return checkTempai(closed, opts)
}

func checkTempai(closed compact.Instances, opts *calc.Options) bool {
	// TODO: optimize
	x := calculate(closed, opts)
	return len(x) > 0
}

// TODO: solve with effectivity
func GetTempaiTiles(closed compact.Instances, options ...calc.Option) compact.Tiles {
	opts := calc.GetOptions(options...)
	if len(opts.Melds)*3+closed.Count() != 14 {
		return 0
	}
	result := compact.Tiles(0)

	closed.Each(func(m compact.Mask) bool {
		i := m.First()
		if !closed.Remove(i) {
			return false
		}
		if checkTempai(closed, opts) {
			result = result.Set(i.Tile())
		}
		closed.Set(i)
		return true
	})
	return result
}
