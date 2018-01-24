package compact

import "github.com/dnovikoff/tempai-core/tile"

// Unique tiles
type Tiles uint64

// 53 is js limit for int64
const _ = uint(53 - tile.End)

const AllTiles = (^Tiles(0)) >> (64 - uint(tile.End))

var KokushiTiles = Tiles(0).SetAll(tile.KokushiTiles)

func NewFromTile(t tile.Tile) Tiles {
	return Tiles(0).Set(t)
}

func NewFromTiles(tiles ...tile.Tile) Tiles {
	x := Tiles(0)
	for _, v := range tiles {
		x = x.Set(v)
	}
	return x
}

func (this Tiles) Check(t tile.Tile) bool {
	return ((1 << uint(t)) & this) != 0
}

func (this Tiles) Set(t tile.Tile) Tiles {
	return ((1 << uint(t)) | this)
}

func (this Tiles) Sub(other Tiles) Tiles {
	return (this | other) ^ other
}

func (this Tiles) Unset(t tile.Tile) Tiles {
	mask := Tiles(1 << uint(t))
	return this &^ mask
}

func (this Tiles) IsEmpty() bool {
	return this == 0
}

func (this Tiles) Merge(other Tiles) Tiles {
	return this | other
}

func (this Tiles) Count() int {
	this = this.Normalize()
	count := 0
	for this != 0 {
		count += int(this & 1)
		this >>= 1
	}
	return count
}

func (this Tiles) Invert() Tiles {
	return (^this).Normalize()
}

func (this Tiles) Normalize() Tiles {
	return ((1 << uint(tile.End)) - 1) & this
}

func (this Tiles) EachRange(begin, end tile.Tile, f func(tile.Tile) bool) bool {
	this >>= uint(begin)

	for i := begin; i < end; i++ {
		if this&1 == 1 {
			if !f(i) {
				return false
			}
		}
		this >>= 1
	}
	return true
}

func (this Tiles) Each(f func(tile.Tile) bool) bool {
	return this.EachRange(tile.Begin, tile.End, f)
}

func (this Tiles) Tiles() tile.Tiles {
	ret := make(tile.Tiles, 0, int(tile.End))
	for i := tile.Begin; i < tile.End; i++ {
		if this&1 == 1 {
			ret = append(ret, i)
		}
		this >>= 1
	}
	return ret
}

func (this Tiles) SetAll(t tile.Tiles) Tiles {
	for _, v := range t {
		this = this.Set(v)
	}
	return this
}
