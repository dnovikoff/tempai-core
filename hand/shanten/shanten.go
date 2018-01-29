package shanten

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/calculator"
	"github.com/dnovikoff/tempai-core/meld"
	"github.com/dnovikoff/tempai-core/tile"
)

func NewCalculator(tiles compact.Instances, used compact.Instances, opened int) *calculator.Calculator {
	return calculator.NewCalculator(meld.AllShantenMelds, tiles, used, opened)
}

func CalculateByMelds(tiles compact.Instances, melds meld.Melds) Results {
	used := compact.NewInstances()
	melds.AddTo(used)
	return Calculate(tiles, len(melds), used)
}

func Calculate(tiles compact.Instances, opened int, used compact.Instances) (ret Results) {
	ret.Regular = CalculateRegular(tiles, opened, used)
	if opened == 0 {
		ret.Pairs = CalculatePairs(tiles)
		ret.Kokushi = CalculateKokushi(tiles)
		ret.Total = ret.Regular.merge(ret.Pairs).merge(ret.Kokushi)
	} else {
		ret.Total = ret.Regular
	}
	return
}

func CalculateRegular(tiles compact.Instances, opened int, used compact.Instances) Result {
	calc := NewCalculator(tiles, used, opened)
	results := calcResult{opened: opened, Result: Result{Value: 8}}
	calc.ResetResult(&results)
	calc.Calculate()
	return results.Result
}

func CalculateKokushi(tiles compact.Instances) Result {
	havePair := false
	count := 0
	missing := compact.Tiles(0)
	for _, t := range tile.KokushiTiles {
		switch tiles.GetMask(t).Count() {
		case 0:
			missing = missing.Set(t)
		case 1:
			count++
		default:
			count++
			havePair = true
		}
	}

	if !havePair {
		missing = compact.KokushiTiles
	} else {
		count++
	}
	return Result{
		Value:    13 - count,
		Improves: missing,
	}
}

func CalculatePairs(tiles compact.Instances) Result {
	value := 6
	improves := compact.Tiles(0)
	tiles.Each(func(m compact.Mask) bool {
		c := m.Count()
		if c == 1 {
			improves = improves.Set(m.Tile())
		} else if c > 1 {
			value--
		}
		return true
	})
	return Result{
		Value:    value,
		Improves: improves,
	}
}
