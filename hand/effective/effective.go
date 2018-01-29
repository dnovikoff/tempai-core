package effective

import (
	"sort"

	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/shanten"
	"github.com/dnovikoff/tempai-core/meld"
	"github.com/dnovikoff/tempai-core/tile"
)

type EffectivityResults map[tile.Tile]shanten.Results

type EffectivityResult struct {
	Tile    tile.Tile
	Shanten shanten.Results
	sortId  sSortId
}

type EffectivityResultsSorted []*EffectivityResult

func CalculateByMelds(closed compact.Instances, melds meld.Melds) EffectivityResults {
	used := compact.NewInstances()
	melds.AddTo(used)
	return Calculate(closed, len(melds), used)
}

func Calculate(closed compact.Instances, opened int, used compact.Instances) (results EffectivityResults) {
	results = make(EffectivityResults)
	cp := closed.Clone()
	closed.Each(func(mask compact.Mask) bool {
		k := mask.Tile()
		first := cp.RemoveTile(k)
		results[k] = shanten.Calculate(cp, opened, used)
		cp.Set(first)
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
	lhs, rhs := l.sortId, r.sortId
	if lhs.betterThan(rhs) {
		return true
	} else if rhs.betterThan(lhs) {
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

func (this EffectivityResults) Sorted(used compact.Instances) EffectivityResultsSorted {
	ret := make(EffectivityResultsSorted, 0, len(this))
	uq := used.UniqueTiles().Invert()
	for k, v := range this {
		t := (uq & v.Total.Improves).Count()
		u := used.CountFree(v.Total.Improves)
		id := newSortId(u, t, v.Total.Value)
		ret = append(ret, &EffectivityResult{
			Tile:    k,
			Shanten: v,
			sortId:  id,
		})
	}

	sort.Sort(ret)
	return ret
}

func (this EffectivityResultsSorted) Best() *EffectivityResult {
	if len(this) == 0 {
		return nil
	}
	return this[0]
}
