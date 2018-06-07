package compact

import (
	"github.com/dnovikoff/tempai-core/tile"
)

type Totals []int

func NewTotals() Totals {
	return make(Totals, tile.TileCount)
}

func (ts Totals) Merge(in Instances) Totals {
	if in == nil {
		return ts
	}
	in.Each(func(mask Mask) bool {
		ts.Add(mask.Tile(), mask.Count())
		return true
	})
	return ts
}

func (ts Totals) Count() int {
	c := 0
	for _, v := range ts {
		c += v
	}
	return c
}

func (ts Totals) Get(t tile.Tile) int {
	return ts[shift(t)]
}

func (ts Totals) Clone() Totals {
	x := NewTotals()
	for k, v := range ts {
		x[k] = v
	}
	return x
}

func (ts Totals) UniqueCount() int {
	c := 0
	for _, v := range ts {
		if v > 0 {
			c++
		}
	}
	return c
}

func (ts Totals) IsFull(t tile.Tile) bool {
	return ts[shift(t)] > 3
}

func (ts Totals) Set(t tile.Tile, d int) {
	ts[shift(t)] = d
}

func (ts Totals) Add(t tile.Tile, d int) {
	ts[shift(t)] += d
}

func (ts Totals) UniqueTiles() Tiles {
	ret := Tiles(0)
	for t, c := range ts {
		if c > 0 {
			ret = ret.Set(tile.Tile(t) + tile.TileBegin)
		}
	}
	return ret
}

func (ts Totals) FullTiles() Tiles {
	ret := Tiles(0)
	for t, c := range ts {
		if c > 3 {
			ret = ret.Set(tile.Tile(t) + tile.TileBegin)
		}
	}
	return ret
}

func (ts Totals) FreeTiles() Tiles {
	return (^ts.FullTiles()).Normalize()
}
