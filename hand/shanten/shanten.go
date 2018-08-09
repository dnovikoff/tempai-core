package shanten

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/calc"
)

func Calculate(tiles compact.Instances, options ...calc.Option) Results {
	opts := calc.GetOptions(options...)
	res := Results{}
	if opts.Forms.Check(calc.Regular) {
		res.Regular = calculateRegular(tiles, opts)
	}
	if opts.Opened == 0 {
		if opts.Forms.Check(calc.Pairs) {
			res.Pairs = CalculatePairs(tiles)
		}
		if opts.Forms.Check(calc.Kokushi) {
			res.Kokushi = CalculateKokushi(tiles)
		}
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
	if tiles.CountBits() != 13 {
		return nil
	}
	improves := compact.Tiles(0)
	pairs := 0
	single := 0
	allTiles := compact.AllTiles
	tiles.Each(func(m compact.Mask) bool {
		c := m.Count()
		if c == 1 {
			single++
			improves = improves.Set(m.Tile())
		} else {
			allTiles = allTiles.Unset(m.Tile())
			pairs++
		}
		return true
	})
	pairsLeft := 7 - pairs
	if single > pairsLeft {
		single = pairsLeft
	} else if single < pairsLeft {
		improves = allTiles
	}
	return &Result{
		Value:    13 - pairs*2 - single,
		Improves: improves,
	}
}
