package shanten

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/calc"
)

func Calculate(tiles compact.Instances, options ...calc.Option) Results {
	opts := calc.GetOptions(options...)
	res := Results{}
	res.Regular = calculateRegular(tiles, opts)
	if opts.Opened == 0 {
		res.Pairs = CalculatePairs(tiles)
		res.Kokushi = CalculateKokushi(tiles)
	}
	res.Total = res.Regular.clone().merge(res.Pairs).merge(res.Kokushi)
	return res
}

func CalculateRegular(tiles compact.Instances, options ...calc.Option) *Result {
	opts := calc.GetOptions(options...)
	return calculateRegular(tiles, opts)
}

func StartMelds(tiles compact.Instances) calc.Option {
	return calc.StartMelds(startMelds(tiles))
}

var shantenMelds = calc.CreateAll()

func startMelds(tiles compact.Instances) calc.Melds {
	melds := shantenMelds.Clone()
	return calc.FilterMelds(tiles, melds)
}

func calculateRegular(tiles compact.Instances, opts *calc.Options) *Result {
	results := calcResult{
		opened: opts.Opened,
		Result: Result{Value: 8},
	}
	opts.Results = &results
	melds := opts.StartMelds
	if melds == nil {
		melds = startMelds(tiles)
	}
	calc.Calculate(melds, tiles, opts)
	return &results.Result
}

func CalculateKokushi(tiles compact.Instances) *Result {
	havePair := false
	count := 0
	missing := compact.Tiles(0)
	for _, t := range compact.TerminalOrHonor.Tiles() {
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
		missing = compact.TerminalOrHonor
	} else {
		count++
	}
	return &Result{
		Value:    13 - count,
		Improves: missing,
	}
}

func CalculatePairs(tiles compact.Instances) *Result {
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
	return &Result{
		Value:    value,
		Improves: improves,
	}
}
