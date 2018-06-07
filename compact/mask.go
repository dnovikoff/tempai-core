package compact

import (
	"github.com/dnovikoff/tempai-core/tile"
)

type Mask uint

const (
	FullMask = 15
)

func MaskByCount(c int) uint {
	return FullMask >> uint(4-c)
}

func NewMask(mask uint, t tile.Tile) Mask {
	m := Mask(shift(t)) << 4
	return m | Mask(mask&15)
}

func (m Mask) Tile() tile.Tile {
	return tile.Tile(m>>4) + tile.TileBegin
}

func (m Mask) Mask() uint {
	return (uint(m) & 15)
}

func (m Mask) Count() int {
	switch m.Mask() {
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

func (m Mask) NaiveCount() int {
	cnt := 0
	for i := 0; i < 4; i++ {
		cnt += int(m & 1)
		m >>= 1
	}
	return cnt
}

func (m Mask) SetCount(in int) Mask {
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
	return NewMask(x, m.Tile())
}

func (m Mask) Instances() tile.Instances {
	ret := make(tile.Instances, 0, 4)
	m.Each(func(t tile.Instance) bool {
		ret = append(ret, t)
		return true
	})
	return ret
}

func (m Mask) InvertTiles() Mask {
	return NewMask(^m.Mask(), m.Tile())
}

func (m Mask) FirstCopy() tile.CopyID {
	switch {
	case m&1 == 1:
		return 0
	case m&2 == 2:
		return 1
	case m&4 == 4:
		return 2
	case m&8 == 8:
		return 3
	}
	return tile.NullCopy
}

func (m Mask) First() tile.Instance {
	c := m.FirstCopy()
	if c == tile.NullCopy {
		return tile.InstanceNull
	}
	return m.Tile().Instance(c)
}

func (m Mask) Each(f func(tile.Instance) bool) bool {
	t := m.Tile()
	for i := tile.CopyID(0); i < 4; i++ {
		if m&1 == 1 {
			if !f(t.Instance(i)) {
				return false
			}
		}
		m >>= 1
	}
	return true
}

func (m Mask) Check(index tile.CopyID) bool {
	return ((1 << uint(index)) & m) != 0
}

func (m Mask) SetIntBit(index uint) Mask {
	return m | (1<<index)&15
}

func (m Mask) SetCopyBit(cid tile.CopyID) Mask {
	return m.SetIntBit(uint(cid))
}

func (m Mask) UnsetIntBit(index uint) Mask {
	mask := Mask(1<<index) & 15
	return m &^ mask
}

func (m Mask) Merge(x Mask) Mask {
	return NewMask(x.Mask()|m.Mask(), m.Tile())
}

func (m Mask) Remove(x Mask) Mask {
	return NewMask(m.Mask()&(^x.Mask()), m.Tile())
}

func (m Mask) UnsetCopyBit(cid tile.CopyID) Mask {
	return m.UnsetIntBit(uint(cid))
}

func (m Mask) IsFull() bool {
	return m.Mask() == FullMask
}

func (m Mask) IsEmpty() bool {
	return m.Mask() == 0
}
