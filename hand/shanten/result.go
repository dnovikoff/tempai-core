package shanten

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/meld"
	"github.com/dnovikoff/tempai-core/tile"
)

type Result struct {
	Value    int
	Improves compact.Tiles
}

type Results struct {
	Regular *Result
	Pairs   *Result
	Kokushi *Result
	Total   *Result
}

type calcResult struct {
	Result
	checked compact.Tiles
	opened  int
}

func (r *Result) clone() *Result {
	x := *r
	return &x
}

func (r *Result) merge(other *Result) *Result {
	if other == nil {
		return r
	}
	if r.Value < other.Value {
		return r
	}
	if r.Value > other.Value {
		return other
	}
	r.Improves = r.Improves | other.Improves
	return r
}

func (r *Result) CalculateUkeIre(total compact.Totals) compact.Totals {
	uke := compact.NewTotals()
	uCount := 0
	r.Improves.Each(func(t tile.Tile) bool {
		c := 4 - total.Get(t)
		if c <= 0 {
			return true
		}
		uke.Set(t, c)
		uCount++
		return true
	})
	return uke
}

func (r *calcResult) CheckMinuses(minuses int) bool {
	// A set without 1 tile is still tempai
	return minuses <= r.Value+1
}

func (r *calcResult) Record(melds meld.Melds, tiles compact.Instances, totals compact.Totals) {
	sets := r.opened
	value := 8 - sets*2
	var improves compact.Tiles
	havePair := false
	for _, v := range melds {
		if v.Type() == meld.TypePair {
			havePair = true
		} else {
			sets++
		}
		value -= reducesShantenBy(v)
		improves |= v.Waits()
	}

	if value > r.Value {
		// avoid useless calculations
		return
	} else if value < r.Value {
		r.checked = 0
	}

	fullSets := sets > 3

	if !havePair {
		tiles.Each(func(m compact.Mask) bool {
			if !totals.IsFull(m.Tile()) {
				improves = improves.Set(m.Tile())
			}
			return true
		})
	}

	if !fullSets {
		toCheck := tiles.UniqueTiles() & (^r.checked)
		try := func(central, improve, wait tile.Tile) {
			if (central.Type() != improve.Type()) || (central.Type() != wait.Type()) {
				return
			}

			if totals.IsFull(wait) || totals.IsFull(improve) {
				return
			}
			improves = improves.Set(improve)
		}
		toCheck.EachRange(tile.TileBegin, tile.SequenceEnd, func(t tile.Tile) bool {
			try(t, t-2, t-1)
			try(t, t-1, t-2)
			try(t, t-1, t+1)
			try(t, t+1, t-1)
			try(t, t+1, t+2)
			try(t, t+2, t+1)
			return true
		})
		r.checked |= toCheck
	}

	r.add(value, improves)
}

func (r *calcResult) add(value int, improves compact.Tiles) {
	if improves == 0 {
		return
	}
	if value > r.Value {
		return
	} else if value < r.Value {
		r.Improves = 0
		r.Value = value
	}
	r.Improves |= improves
}

func reducesShantenBy(m meld.Meld) int {
	if m == 0 {
		return 0
	}
	if m.Type() == meld.TypePair {
		if meld.Pair(m).IsComplete() {
			return 1
		} else {
			return 0
		}
	}
	if m.IsComplete() {
		return 2
	}
	return 1
}
