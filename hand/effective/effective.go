package effective

import (
	"sort"

	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/calc"
	"github.com/dnovikoff/tempai-core/hand/shanten"
	"github.com/dnovikoff/tempai-core/tile"
)

type Results map[tile.Tile]shanten.Results

type Result struct {
	Tile    tile.Tile
	Shanten shanten.Results
	sortId  sSortId
}

func Calculate(closed compact.Instances, options ...calc.Option) Results {
	results := make(Results)
	cp := closed.Clone()
	closed.Each(func(mask compact.Mask) bool {
		i := mask.First()
		if !cp.Remove(i) {
			return false
		}
		results[i.Tile()] = shanten.Calculate(cp, options...)
		cp.Set(i)
		return true
	})
	return results
}

//wd98765
func tilePriority(t tile.Tile) int {
	if compact.Honor.Check(t) {
		return 0
	}

	switch t.Number() {
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

func (r Results) Sorted(used compact.Instances) ResultsSorted {
	ret := make(ResultsSorted, 0, len(r))
	uq := used.UniqueTiles().Invert()
	for k, v := range r {
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
