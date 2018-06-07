package compact

import (
	"github.com/dnovikoff/tempai-core/tile"
)

// Not more than 4 tiles per type by implementation
type Instances []PackedMasks

const (
	instancesBits = 32
	instancesInts = 5
	tilesPerPack  = 8
)

const (
	counterIndex = instancesInts + iota
	endIndex
)
const _ = uint(tilesPerPack*instancesInts*4 - tile.InstanceCount)

const counterInvalid = PackedMasks(1024)

func AllInstancesFromTo(from, to tile.Tile) Instances {
	x := NewInstances()
	for t := from; t < to; t++ {
		x.SetCount(t, 4)
	}
	return x
}

func AllInstances() Instances {
	return AllInstancesFromTo(tile.TileBegin, tile.TileEnd)
}

func NewInstances() Instances {
	return make(Instances, endIndex)
}

func (is Instances) Each(f func(mask Mask) bool) bool {
	start := tile.Man1
	for _, v := range is.packed() {
		cur := start
		for v != 0 {
			mask := uint(v & 15)
			if mask != 0 {
				if !f(NewMask(mask, cur)) {
					return false
				}
			}
			cur++
			v >>= 4
		}
		start += tilesPerPack
	}
	return true
}

func (is Instances) invalidateCounter() {
	is[counterIndex] = counterInvalid
}

func (is Instances) GetMask(t tile.Tile) Mask {
	block := is[int(shift(t)/tilesPerPack)]
	return block.Get(shift(t)%tilesPerPack, t)
}

func (is Instances) CountFree(in Tiles) int {
	result := 0
	in.Each(func(t tile.Tile) bool {
		result += 4 - is.GetCount(t)
		return true
	})
	return result
}

func (is Instances) GetFree() Tiles {
	return ^is.GetFull()
}

func (is Instances) GetFull() Tiles {
	result := Tiles(0)

	is.Each(func(m Mask) bool {
		if m.IsFull() {
			result = result.Set(m.Tile())
		}
		return true
	})
	return result
}

func (is Instances) Invert() Instances {
	return is.CopyFree(AllTiles)
}

func (is Instances) CopyFree(in Tiles) Instances {
	result := NewInstances()
	count := 0
	in.Each(func(t tile.Tile) bool {
		i := is.GetMask(t).InvertTiles()
		result.SetMask(i)
		count += i.Count()
		return true
	})
	result[counterIndex] = PackedMasks(count)
	return result
}

func (is Instances) CopyFrom(x Instances) {
	for k, v := range x {
		is[k] = v
	}
}

func (is Instances) extract(t tile.Tile, count int) Mask {
	original := is.GetMask(t) & 15
	if original.Count() < count {
		return 0
	}
	eraser := ^Mask(0)
	for (original & eraser).Count() > count {
		eraser <<= 1
	}

	is.setMaskImpl(uint(t), (^eraser)&original)
	return original & eraser & 15
}

func (is Instances) Merge(other Instances) Instances {
	for k, v := range is[:counterIndex] {
		is[k] = v | other[k]
	}
	is.invalidateCounter()
	return is
}

func (is Instances) packed() Instances {
	return is[:counterIndex]
}

func (is Instances) setMaskImpl(index uint, mask Mask) {
	blocknum := index / tilesPerPack
	shift := index % tilesPerPack
	is[blocknum] = is[blocknum].Set(mask, shift)
}

func (is Instances) SetMask(mask Mask) {
	is.setMaskImpl(shift(mask.Tile()), mask)
	is.invalidateCounter()
}

func (is Instances) Add(t tile.Instances) Instances {
	for _, v := range t {
		is.Set(v)
	}
	return is
}

func (is Instances) Clone() Instances {
	clone := make(Instances, len(is))
	clone.CopyFrom(is)
	return clone
}

func (is Instances) AddCount(t tile.Tile, x int) Instances {
	is.SetCount(t, is.GetMask(t).Count()+x)
	return is
}

func (is Instances) SetCount(t tile.Tile, x int) {
	m := NewMask(MaskByCount(x), t)
	is.SetMask(m)
}

func (is Instances) GetCount(t tile.Tile) int {
	return is.GetMask(t).Count()
}

func (is Instances) Set(t tile.Instance) {
	current := is.GetMask(t.Tile())
	is.SetMask(current.SetCopyBit(t.CopyID()))
}

func (is Instances) Check(t tile.Instance) bool {
	return is.GetMask(t.Tile()).Check(t.CopyID())
}

func (is Instances) Remove(t tile.Instance) bool {
	current := is.GetMask(t.Tile())
	next := current.UnsetCopyBit(t.CopyID())
	is.SetMask(next)
	return next != current
}

func (is Instances) UniqueTiles() Tiles {
	cts := Tiles(0)
	start := tile.TileBegin
	for _, v := range is.packed() {
		t := start
		for v != 0 {
			if (v & 15) != 0 {
				cts = cts.Set(t)
			}
			t++
			v >>= 4
		}
		start += tilesPerPack
	}
	return cts
}

func (is Instances) Instances() tile.Instances {
	ret := make(tile.Instances, is.Count())
	x := 0
	is.Each(func(mask Mask) bool {
		return mask.Each(func(inst tile.Instance) bool {
			ret[x] = inst
			x++
			return true
		})
	})
	return ret
}

func (is Instances) Count() int {
	val := is[counterIndex]
	if val == counterInvalid {
		return is.recountImpl()
	}
	return int(val)
}

func (is Instances) recountImpl() int {
	x := 0
	for _, v := range is.packed() {
		x += v.CountBits()
	}
	is[counterIndex] = PackedMasks(x)
	return x
}
