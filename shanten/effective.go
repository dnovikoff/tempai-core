package shanten

import (
	"sort"

	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/meld"
	"github.com/dnovikoff/tempai-core/tile"
)

type EffectivityResults map[tile.Tile]*ShantenResults

type EffectivityResult struct {
	Tile    tile.Tile
	Shanten *ShantenResults
}

type EffectivityResultsSorted []*EffectivityResult

func CalculateEffectivityMelds(closed compact.Instances, melds meld.Melds) EffectivityResults {
	used := compact.NewInstances()
	melds.AddTo(used)
	return CalculateEffectivity(closed, len(melds), used)
}

func CalculateEffectivity(closed compact.Instances, opened int, used compact.Instances) EffectivityResults {
	res, _ := CalculateEffectivityBoth(closed, opened, used)
	return res
}

type effectiveChanResult struct {
	tile tile.Tile
	all  *ShantenResults
	reg  *ShantenResults
}

func CalculateEffectivityBoth(closed compact.Instances, opened int, used compact.Instances) (results EffectivityResults, regular EffectivityResults) {
	results = make(EffectivityResults)
	regular = make(EffectivityResults)
	closed.Each(func(mask compact.Mask) bool {
		k := mask.Tile()
		first := closed.RemoveTile(k)

		all, reg := CalculateShantenBoth(closed, opened, used)
		results[k] = all
		regular[k] = reg

		closed.Set(first)
		return true
	})
	return
}

func CalculateEffectivitySpecial(calc *Calculator) (results EffectivityResults) {
	results = make(EffectivityResults, 13)

	closed := calc.tiles
	closed.Each(func(mask compact.Mask) bool {
		k := mask.Tile()
		first := closed.RemoveTile(k)
		res := NewShantenResults(0)
		calc.ResetResult(res)
		calc.Calculate()
		res.CalculatePairs(calc.tiles)
		res.Calculate13(calc.tiles)

		results[k] = res.Recount(calc.totals)
		closed.Set(first)
		return true
	})
	return
}

func CalculateEffectivity13(calc *Calculator) (results EffectivityResults) {
	results = make(EffectivityResults, 13)

	closed := calc.tiles
	closed.Each(func(mask compact.Mask) bool {
		k := mask.Tile()
		first := closed.RemoveTile(k)

		res := NewShantenResults(0)
		res.Calculate13(closed)

		results[k] = res.Recount(calc.totals)
		closed.Set(first)
		return true
	})
	return
}

func (this EffectivityResultsSorted) Len() int {
	return len(this)
}

func (this EffectivityResultsSorted) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

//wd98765
func tilePriority(t tile.Tile) int {
	if t.IsHonor() {
		return 0
	}

	switch t.NumberInSequence() {
	case 1, 9:
		return 1
	case 2, 8:
		return 2
	case 3, 7:
		return 3
	case 4, 6:
		return 5
	case 5:
		return 6
	}
	// Unreachable
	return 7
}

func tileSortId(t tile.Tile) int {
	return (tilePriority(t) << 8) | int(t)
}

func tileLess(l, r tile.Tile) bool {
	return tileSortId(l) < tileSortId(r)
}

func effLess(l, r *EffectivityResult) bool {
	lhs, rhs := l.Shanten.SortId(), r.Shanten.SortId()
	if lhs.BetterThan(rhs) {
		return true
	} else if rhs.BetterThan(lhs) {
		return false
	}
	return tileLess(l.Tile, r.Tile)
}

func (this EffectivityResultsSorted) Less(i, j int) bool {
	return effLess(this[i], this[j])
}

func (this EffectivityResultsSorted) First() (ret *EffectivityResult) {
	if len(this) == 0 {
		return
	}
	return this[0]
}

func (this EffectivityResults) Sorted() EffectivityResultsSorted {
	ret := make(EffectivityResultsSorted, 0, len(this))
	for k, v := range this {
		ret = append(ret, &EffectivityResult{Tile: k, Shanten: v})
	}

	sort.Sort(ret)
	return ret
}

func (this EffectivityResults) Best() *EffectivityResult {
	if len(this) == 0 {
		return nil
	}
	var best *EffectivityResult

	for k, v := range this {
		next := &EffectivityResult{Tile: k, Shanten: v}
		if best == nil || effLess(next, best) {
			best = next
		}
	}

	return best
}
