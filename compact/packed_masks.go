package compact

import "bitbucket.org/dnovikoff/tempai-core/tile"

// Using 32 implementation, because gopherjs supports only 53 bits
type PackedMasks uint32

func SinglePackedMasks(mask Mask, index uint) PackedMasks {
	return PackedMasks(mask.Mask() << (4 * index))
}

func (this PackedMasks) Set(mask Mask, index uint) PackedMasks {
	erase := ^SinglePackedMasks(15, index)
	return (this & erase) | SinglePackedMasks(mask, index)
}

func (this PackedMasks) Get(index uint, tile tile.Tile) Mask {
	return NewMask(uint(this)>>(4*index), tile)
}

func (this PackedMasks) Each(start, end tile.Tile, skipEmpty bool, f func(mask Mask) bool) bool {
	for i := 0; i < tilesPerPack; i++ {
		mask := NewMask(uint(this), start)
		if !skipEmpty || !mask.IsEmpty() {
			if !f(mask) {
				return false
			}
		}
		start++
		if start >= end {
			return true
		}
		this >>= 4
	}
	return true
}

func (this PackedMasks) CountBits() int {
	cnt := 0

	for this > 0 {
		cnt += int(this & 1)
		this >>= 1
	}
	return cnt
}
