package compact

import "github.com/dnovikoff/tempai-core/tile"

// Unique tiles
type Tiles uint64

func NewTiles(x ...tile.Tile) Tiles {
	r := Tiles(0)
	for _, v := range x {
		r = r.Set(v)
	}
	return r
}

func FromTile(t tile.Tile) Tiles {
	return Tiles(0).Set(t)
}

func FromTiles(tiles ...tile.Tile) Tiles {
	x := Tiles(0)
	for _, v := range tiles {
		x = x.Set(v)
	}
	return x
}

func shift(t tile.Tile) uint {
	return uint(t - tile.TileBegin)
}

func (ts Tiles) Check(t tile.Tile) bool {
	return ((1 << shift(t)) & ts) != 0
}

func (ts Tiles) Set(t tile.Tile) Tiles {
	return ((1 << shift(t)) | ts)
}

func (ts Tiles) Sub(other Tiles) Tiles {
	return (ts | other) ^ other
}

func (ts Tiles) Unset(t tile.Tile) Tiles {
	mask := Tiles(1 << shift(t))
	return ts &^ mask
}

func (ts Tiles) IsEmpty() bool {
	return ts == 0
}

func (ts Tiles) Merge(other Tiles) Tiles {
	return ts | other
}

func (ts Tiles) Count() int {
	ts = ts.Normalize()
	count := 0
	for ts != 0 {
		count += int(ts & 1)
		ts >>= 1
	}
	return count
}

func (ts Tiles) Invert() Tiles {
	return (^ts).Normalize()
}

func (ts Tiles) Normalize() Tiles {
	return ((1 << shift(tile.TileEnd)) - 1) & ts
}

func (ts Tiles) EachRange(begin, end tile.Tile, f func(tile.Tile) bool) bool {
	ts >>= shift(begin)
	for i := begin; ts != 0 && i < end; i++ {
		if ts&1 == 1 {
			if !f(i) {
				return false
			}
		}
		ts >>= 1
	}
	return true
}

func (ts Tiles) Each(f func(tile.Tile) bool) bool {
	return ts.EachRange(tile.TileBegin, tile.TileEnd, f)
}

func (ts Tiles) Tiles() tile.Tiles {
	ret := make(tile.Tiles, 0, tile.TileCount)
	for i := tile.TileBegin; i < tile.TileEnd; i++ {
		if ts&1 == 1 {
			ret = append(ret, i)
		}
		ts >>= 1
	}
	return ret
}

func (ts Tiles) SetAll(t tile.Tiles) Tiles {
	for _, v := range t {
		ts = ts.Set(v)
	}
	return ts
}
