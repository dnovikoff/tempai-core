package shanten

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/calc"
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

func (r *calcResult) Record(in *calc.ResultData) {
	if !in.Validator.Validate(in.Closed) {
		return
	}
	value := 8 - r.opened*2
	var improves compact.Tiles
	for _, v := range in.Closed {
		value -= reducesShantenBy(v)
		improves |= v.CompactWaits()
	}
	if in.Pair != nil {
		value--
	}

	if value > r.Value {
		// avoid useless calculations
		return
	} else if value < r.Value {
		r.checked = 0
	}

	fullSets := in.Sets > 3

	if in.Pair == nil {
		for _, t := range in.Left {
			improves = improves.Set(t)
		}
	}

	if !fullSets {
		toCheck := compact.FromTiles(in.Left...) & (^r.checked)
		try := func(central, improve, wait tile.Tile) {
			if (central.Type() != improve.Type()) || (central.Type() != wait.Type()) {
				return
			}

			if in.Validator.Empty(wait) || in.Validator.Empty(improve) {
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

	// Cleanup
	for x := tile.TileBegin; x < tile.TileEnd; x++ {
		if in.Validator.Empty(x) {
			improves = improves.Unset(x)
		}
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

func reducesShantenBy(m calc.Meld) int {
	if m == nil {
		return 0
	}
	tags := m.Tags()
	if tags.CheckAny(calc.TagPair) {
		if tags.CheckAny(calc.TagComplete) {
			return 1
		}
		return 0
	}
	if tags.CheckAny(calc.TagComplete) {
		return 2
	}
	return 1
}
