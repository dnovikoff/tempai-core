package compact

import (
	"github.com/dnovikoff/tempai-core/tile"
)

// Using 32 implementation, because gopherjs supports only 53 bits
type PackedMasks uint32

func SinglePackedMasks(mask Mask, index uint) PackedMasks {
	return PackedMasks(mask.Mask() << (4 * index))
}

func (pm PackedMasks) Set(mask Mask, index uint) PackedMasks {
	erase := ^SinglePackedMasks(15, index)
	return (pm & erase) | SinglePackedMasks(mask, index)
}

func (pm PackedMasks) Get(index uint, tile tile.Tile) Mask {
	return NewMask(uint(pm)>>(4*index), tile)
}

func (pm PackedMasks) Each(start, end tile.Tile, skipEmpty bool, f func(mask Mask) bool) bool {
	for i := 0; i < tilesPerPack; i++ {
		mask := NewMask(uint(pm), start)
		if !skipEmpty || !mask.IsEmpty() {
			if !f(mask) {
				return false
			}
		}
		start++
		if start >= end {
			return true
		}
		pm >>= 4
	}
	return true
}

func (pm PackedMasks) CountBits() int {
	cnt := 0

	for pm > 0 {
		cnt += int(pm & 1)
		pm >>= 1
	}
	return cnt
}
