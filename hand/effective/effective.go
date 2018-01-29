package effective

import (
	"sort"

	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/shanten"
	"github.com/dnovikoff/tempai-core/meld"
	"github.com/dnovikoff/tempai-core/tile"
)

type Results map[tile.Tile]shanten.Results

type Result struct {
	Tile    tile.Tile
	Shanten shanten.Results
	sortId  sSortId
}

func CalculateByMelds(closed compact.Instances, melds meld.Melds) Results {
	used := compact.NewInstances()
	melds.AddTo(used)
	return Calculate(closed, len(melds), used)
}

func Calculate(closed compact.Instances, opened int, used compact.Instances) (results Results) {
	results = make(Results)
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

func (this Results) Sorted(used compact.Instances) ResultsSorted {
	ret := make(ResultsSorted, 0, len(this))
	uq := used.UniqueTiles().Invert()
	for k, v := range this {
		t := (uq & v.Total.Improves).Count()
		u := used.CountFree(v.Total.Improves)
		id := newSortId(u, t, v.Total.Value)
		ret = append(ret, &Result{
			Tile:    k,
			Shanten: v,
			sortId:  id,
		})
	}

	sort.Sort(ret)
	return ret
}
