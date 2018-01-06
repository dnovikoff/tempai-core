package compact

import "bitbucket.org/dnovikoff/tempai-core/tile"

type Totals []int

func NewTotals() Totals {
	t := make(Totals, int(tile.End))
	return t
}

func (this Totals) Merge(in Instances) Totals {
	if in == nil {
		return this
	}
	in.Each(func(mask Mask) bool {
		this.Add(mask.Tile(), mask.Count())
		return true
	})
	return this
}

func (this Totals) Count() int {
	c := 0
	for _, v := range this {
		c += v
	}
	return c
}

func (this Totals) Get(t tile.Tile) int {
	return this[int(t)]
}

func (this Totals) Clone() Totals {
	x := NewTotals()
	for k, v := range this {
		x[k] = v
	}
	return x
}

func (this Totals) UniqueCount() int {
	c := 0
	for _, v := range this {
		if v > 0 {
			c++
		}
	}
	return c
}

func (this Totals) IsFull(t tile.Tile) bool {
	return this[int(t)] > 3
}

func (this Totals) Set(t tile.Tile, d int) {
	this[int(t)] = d
}

func (this Totals) Add(t tile.Tile, d int) {
	this[int(t)] += d
}

func (this Totals) UniqueTiles() Tiles {
	ret := Tiles(0)
	for t, c := range this {
		if c > 0 {
			ret = ret.Set(tile.Tile(t))
		}
	}
	return ret
}

func (this Totals) FullTiles() Tiles {
	ret := Tiles(0)
	for t, c := range this {
		if c > 3 {
			ret = ret.Set(tile.Tile(t))
		}
	}
	return ret
}

func (this Totals) FreeTiles() Tiles {
	return ^this.FullTiles()
}
