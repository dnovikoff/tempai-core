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

func NewAllInstancesFromTo(from, to tile.Tile) Instances {
	x := NewInstances()
	for t := from; t < to; t++ {
		x.SetCount(t, 4)
	}
	return x
}

func NewAllInstances() Instances {
	return NewAllInstancesFromTo(tile.Begin, tile.End)
}

func NewInstances() Instances {
	return make(Instances, endIndex)
}

func (this Instances) Each(f func(mask Mask) bool) bool {
	start := tile.Man1
	for _, v := range this.packed() {
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

func (this Instances) EachTile(f func(t tile.Tile) bool) bool {
	start := tile.Man1
	for _, v := range this.packed() {
		cur := start
		for v != 0 {
			if uint(v&15) != 0 {
				if !f(cur) {
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

func (this Instances) invalidateCounter() {
	this[counterIndex] = counterInvalid
}

func (this Instances) GetMask(t tile.Tile) Mask {
	block := this[int(shift(t)/tilesPerPack)]
	return block.Get(shift(t)%tilesPerPack, t)
}

func (this Instances) CountFree(in Tiles) int {
	result := 0
	in.Each(func(t tile.Tile) bool {
		result += 4 - this.GetCount(t)
		return true
	})
	return result
}

func (this Instances) GetFree() Tiles {
	return ^this.GetFull()
}

func (this Instances) GetFull() Tiles {
	result := Tiles(0)

	this.Each(func(m Mask) bool {
		if m.IsFull() {
			result = result.Set(m.Tile())
		}
		return true
	})
	return result
}

func (this Instances) Invert() Instances {
	return this.CopyFree(AllTiles)
}

func (this Instances) CopyFree(in Tiles) Instances {
	result := NewInstances()
	count := 0
	in.Each(func(t tile.Tile) bool {
		i := this.GetMask(t).InvertTiles()
		result.SetMask(i)
		count += i.Count()
		return true
	})
	result[counterIndex] = PackedMasks(count)
	return result
}

func (this Instances) CopyFrom(x Instances) {
	for k, v := range x.all() {
		this[k] = v
	}
}

func (this Instances) extract(t tile.Tile, count int) Mask {
	original := this.GetMask(t) & 15
	if original.Count() < count {
		return 0
	}
	eraser := ^Mask(0)
	for (original & eraser).Count() > count {
		eraser <<= 1
	}

	this.setMaskImpl(uint(t), (^eraser)&original)
	return original & eraser & 15
}

func (this Instances) Merge(other Instances) Instances {
	for k, v := range this.all() {
		this[k] = v | other[k]
	}
	this.invalidateCounter()
	return this
}

func (this Instances) packed() Instances {
	return this[:counterIndex]
}

func (this Instances) all() Instances {
	return this[:endIndex]
}

func (this Instances) setMaskImpl(index uint, mask Mask) {
	blocknum := index / tilesPerPack
	shift := index % tilesPerPack
	this[blocknum] = this[blocknum].Set(mask, shift)
}

func (this Instances) SetMask(mask Mask) {
	this.setMaskImpl(shift(mask.Tile()), mask)
	this.invalidateCounter()
}

func (this Instances) Add(t tile.Instances) Instances {
	for _, v := range t {
		this.Set(v)
	}
	return this
}

func (this Instances) Clone() Instances {
	clone := make(Instances, len(this))
	for k, v := range this.all() {
		clone[k] = v
	}
	return clone
}

func (this Instances) AddCount(t tile.Tile, x int) Instances {
	m := this.GetMask(t)
	val := m.Count() + x
	if val < 0 {
		val = 0
	} else if val > 4 {
		val = 4
	}
	this.SetCount(t, val)
	return this
}

func (this Instances) AddCounts(t tile.Tiles) Instances {
	for _, v := range t {
		this.AddCount(v, 1)
	}
	return this
}

func (this Instances) CheckEmpty(t tile.Tile) bool {
	return this.GetMask(t).IsEmpty()
}

func (this Instances) CheckFull(t tile.Tile) bool {
	return this.GetMask(t).IsFull()
}

func (this Instances) SetCount(t tile.Tile, x int) {
	this.SetMask(NewMaskByCount(x, t))
}

func (this Instances) GetCount(t tile.Tile) int {
	return this.GetMask(t).Count()
}

func (this Instances) Set(t tile.Instance) {
	current := this.GetMask(t.Tile())
	this.SetMask(current.SetCopyBit(t.CopyId()))
}

func (this Instances) RemoveAll(t tile.Tile) {
	this.SetMask(this.GetMask(t).SetCount(0))
}

func (this Instances) RemoveTile(t tile.Tile) tile.Instance {
	mask := this.GetMask(t)
	first := mask.First()
	if first != tile.InstanceNull {
		this.SetMask(mask.UnsetInstance(first))
	}
	return first
}

func (this Instances) Check(t tile.Instance) bool {
	return this.GetMask(t.Tile()).Check(t.CopyId())
}

func (this Instances) Remove(t tile.Instance) bool {
	current := this.GetMask(t.Tile())
	next := current.UnsetCopyBit(t.CopyId())
	this.SetMask(next)
	return next != current
}

func (this Instances) UniqueTiles() Tiles {
	cts := Tiles(0)
	start := tile.Begin
	for _, v := range this.packed() {
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

func (this Instances) UniqueCount() int {
	cnt := 0
	for _, val := range this.packed() {
		for val != 0 {
			if val&15 != 0 {
				cnt++
			}
			val >>= 4
		}
	}
	return cnt
}

func (this Instances) Instances() tile.Instances {
	ret := make(tile.Instances, this.Count())
	x := 0
	this.Each(func(mask Mask) bool {
		return mask.Each(func(inst tile.Instance) bool {
			ret[x] = inst
			x++
			return true
		})
	})
	return ret
}

func (this Instances) Count() int {
	val := this[counterIndex]
	if val == counterInvalid {
		return this.recountImpl()
	}
	return int(val)
}

func (this Instances) recountImpl() int {
	x := 0
	for _, v := range this.packed() {
		x += v.CountBits()
	}
	this[counterIndex] = PackedMasks(x)
	return x
}
