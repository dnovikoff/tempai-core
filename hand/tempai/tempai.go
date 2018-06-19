package tempai

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/calc"
	"github.com/dnovikoff/tempai-core/tile"
)

func calculatePairs(tiles compact.Instances) *TempaiResult {
	var wait calc.Meld
	pairs := make(calc.Melds, 0, 7)
	if !tiles.Each(func(m compact.Mask) bool {
		switch m.Count() {
		case 1:
			wait = calc.Tanki(m.Tile())
		case 2:
			head := calc.Pair(m.Tile())
			pairs = append(pairs, head)
		default:
			return false
		}
		return true
	}) {
		return nil
	}
	if wait == nil || len(pairs) != 6 {
		return nil
	}

	return &TempaiResult{
		Type:  TypePairs,
		Melds: pairs,
		Last:  wait,
		Waits: wait.CompactWaits(),
	}
}

func calculateKokushi(tiles compact.Instances) *TempaiResult {
	pair := tile.TileNull
	mask := compact.TerminalOrHonor
	if !tiles.Each(func(m compact.Mask) bool {
		t := m.Tile()
		if !compact.TerminalOrHonor.Check(t) {
			return false
		}
		mask = mask.Unset(t)
		switch m.Count() {
		case 1:
		case 2:
			if pair != tile.TileNull {
				return false
			}
			pair = t
		default:
			return false
		}
		return true
	}) {
		return nil
	}
	res := &TempaiResult{
		Type: TypeKokushi,
	}
	if mask.IsEmpty() {
		res.Last = calc.Kokushi13()
	} else {
		res.Last = calc.KokushiMeld(pair, mask.Tiles()[0])
	}
	res.Waits = res.Last.CompactWaits()
	return res
}

func StartMelds(tiles compact.Instances) calc.Option {
	return calc.StartMelds(startMelds(tiles))
}

var tempaiMelds = calc.CreateComplete()

func startMelds(tiles compact.Instances) calc.Melds {
	return calc.FilterMelds(tiles, tempaiMelds)
}

func calculateRegular(ret *TempaiResults, opts *calc.Options) bool {
	if !opts.Forms.Check(calc.Regular) {
		return false
	}
	if opts.Opened*3+ret.Hand.CountBits() != 13 {
		return false
	}
	x := &result{
		original: ret.Hand,
		tmp:      compact.NewInstances(),
	}
	opts.Results = x
	melds := opts.StartMelds
	if melds == nil {
		melds = startMelds(ret.Hand)
	}
	calc.Calculate(melds, ret.Hand, opts)
	if len(x.results) == 0 {
		return false
	}
	ret.Results = append(ret.Results, x.results...)
	return true
}

func Calculate(closed compact.Instances, options ...calc.Option) *TempaiResults {
	opts := calc.GetOptions(options...)
	ret := calculate(closed, opts)
	if len(ret.Results) == 0 {
		return nil
	}
	return ret
}

func calculate(closed compact.Instances, opts *calc.Options) *TempaiResults {
	ret := &TempaiResults{
		Hand:     closed,
		Declared: opts.Declared,
	}
	calculateRegular(ret, opts)
	if opts.Opened != 0 {
		return ret
	}
	// Pairs should be calcluated after regular hand
	// because of ryanpeiko

	if opts.Forms.Check(calc.Pairs) {
		res := calculatePairs(closed)
		if res != nil {
			ret.Results = append(ret.Results, res)
		}
	}
	if opts.Forms.Check(calc.Kokushi) {
		res := calculateKokushi(closed)
		if res != nil {
			ret.Results = append(ret.Results, res)
		}
	}
	return ret
}

func CheckTempai(closed compact.Instances, options ...calc.Option) bool {
	opts := calc.GetOptions(options...)
	return checkTempai(closed, opts)
}

func checkTempai(closed compact.Instances, opts *calc.Options) bool {
	// TODO: optimize
	x := calculate(closed, opts)
	return len(x.Results) > 0
}

// TODO: solve with effectivity
func GetTempaiTiles(closed compact.Instances, options ...calc.Option) compact.Tiles {
	opts := calc.GetOptions(options...)
	if opts.Opened*3+closed.CountBits() != 14 {
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
