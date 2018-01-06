package compact

import "bitbucket.org/dnovikoff/tempai-core/tile"

type Mask uint

func NewMask(mask uint, t tile.Tile) Mask {
	this := Mask(t) << 4
	return this | Mask(mask&15)
}

func NewMaskByCount(count int, t tile.Tile) Mask {
	return NewMask(0, t).SetCount(count)
}

func (this Mask) Tile() tile.Tile {
	return tile.Tile(this >> 4)
}

func (this Mask) Mask() uint {
	return (uint(this) & 15)
}

func (this Mask) Count() int {
	switch this.Mask() {
	case 0:
		return 0
	case 1, 2, 4, 8:
		return 1
	case 3, 5, 6, 9, 10, 12:
		return 2
	case 7, 11, 13, 14:
		return 3
	case 15:
		return 4
	}
	return 0
}

func (this Mask) NaiveCount() int {
	cnt := 0
	for i := 0; i < 4; i++ {
		cnt += int(this & 1)
		this >>= 1
	}
	return cnt
}

func (this Mask) SetCount(in int) Mask {
	x := uint(0)
	switch in {
	case 0:
	case 1:
		x = 1
	case 2:
		x = 1 + 2
	case 3:
		x = 1 + 2 + 4
	case 4:
		x = 1 + 2 + 4 + 8
	}
	return NewMask(x, this.Tile())
}

func (this Mask) Instances() tile.Instances {
	ret := make(tile.Instances, 0, 4)
	this.Each(func(t tile.Instance) bool {
		ret = append(ret, t)
		return true
	})
	return ret
}

func (this Mask) InvertTiles() Mask {
	result := NewMask(^this.Mask(), this.Tile())
	return result
}

func (this Mask) FirstCopy() tile.CopyId {
	switch {
	case this&1 == 1:
		return 0
	case this&2 == 2:
		return 1
	case this&4 == 4:
		return 2
	case this&8 == 8:
		return 3
	}
	return tile.NullCopy
}

func (this Mask) First() tile.Instance {
	c := this.FirstCopy()
	if c == tile.NullCopy {
		return tile.InstanceNull
	}
	return this.Tile().Instance(c)
}

func (this Mask) Each(f func(tile.Instance) bool) bool {
	t := this.Tile()
	for i := tile.CopyId(0); i < 4; i++ {
		if this&1 == 1 {
			if !f(t.Instance(i)) {
				return false
			}
		}
		this >>= 1
	}
	return true
}

func (this Mask) Check(index tile.CopyId) bool {
	return ((1 << uint(index)) & this) != 0
}

func (this Mask) SetIntBit(index uint) Mask {
	return this | (1<<index)&15
}

func (this Mask) SetCopyBit(cid tile.CopyId) Mask {
	return this.SetIntBit(uint(cid))
}

func (this Mask) UnsetIntBit(index uint) Mask {
	mask := Mask(1<<index) & 15
	return this &^ mask
}

func (this Mask) UnsetInstance(i tile.Instance) Mask {
	return this.UnsetCopyBit(i.CopyId())
}

func (this Mask) UnsetInstances(i tile.Instances) Mask {
	for _, v := range i {
		this = this.UnsetCopyBit(v.CopyId())
	}
	return this
}

func (this Mask) SetInstances(i tile.Instances) Mask {
	for _, v := range i {
		this = this.SetCopyBit(v.CopyId())
	}
	return this
}

func (this Mask) Merge(m Mask) Mask {
	return NewMask(m.Mask()|this.Mask(), this.Tile())
}

func (this Mask) Remove(m Mask) Mask {
	return NewMask(this.Mask()&(^m.Mask()), this.Tile())
}

func (this Mask) UnsetCopyBit(cid tile.CopyId) Mask {
	return this.UnsetIntBit(uint(cid))
}

func (this Mask) IsFull() bool {
	return this.Mask() == 15
}

func (this Mask) IsEmpty() bool {
	return this.Mask() == 0
}
