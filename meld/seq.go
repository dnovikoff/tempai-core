package meld

import (
	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/tile"
)

// Seq is used to represent Chi
// Opponent for Seq is always Left
// Type - 2 | Chibase -5 | TileCopies - 2-2-2 | Complete -1 | Opened - 2 |
// TilesCopies - values 0-3
// ChiBase - 1-27
// Complete. 0 - means Opened will be treated as 'Hole' Tile
// Opened - opened tile or hole
// 16 In Total
type Seq Meld

var _ Interface = Seq(0)

const (
	HoleCopy      tile.CopyID = -5
	seqCopiesMask             = 63 << 7
)

type seqCopies int

func OpenCopy(c tile.CopyID) tile.CopyID {
	return -(c + 1)
}

func getCopy(c tile.CopyID) tile.CopyID {
	if c < 0 {
		return -(c + 1)
	}
	return c
}

func newSeqCopies(c1, c2, c3 tile.CopyID) seqCopies {
	x := seqCopies(c3)
	x = (x << 2) | (seqCopies(c2) & 3)
	x = (x << 2) | (seqCopies(c1) & 3)
	return x
}

func (c seqCopies) at(idx int) tile.CopyID {
	return tile.CopyID((c >> (uint(idx) * 2)) & 3)
}

func newSeq(base tile.Tile, copies seqCopies, complete int, opened int) Seq {
	x := opened & 3
	x = (x << 1) | (complete & 1)
	x = (x << 6) | int(copies&63)
	x = (x << 5) | (int(base) & 31)
	x = (x << 2) | int(TypeSeq)
	return Seq(x)
}

func newSeq2(base tile.Tile, c1, c2, c3 tile.CopyID, complete int, opened int) Seq {
	return newSeq(base, newSeqCopies(c1, c2, c3), complete, opened)
}

func NewSeq(base tile.Tile, t1, t2, t3 tile.CopyID) Seq {
	if !base.IsSequence() {
		return 0
	}
	switch base.Number() {
	case 8:
		// special case
		if t3 != HoleCopy {
			return 0
		}
		return NewSeq(base-1, HoleCopy, t1, t2)
	case 9:
		return 0
	}

	c := 0
	isOpened := false
	check := func(t tile.CopyID) {
		if t < 0 {
			if t != HoleCopy {
				t = getCopy(t)
				isOpened = true
			}
			c++
		}
		if t > 3 {
			c = 100
		}
	}
	check(t1)
	check(t2)
	check(t3)
	switch c {
	case 1:
		if !isOpened {
			if t1 == HoleCopy {
				if base.Number() == 7 {
					return newSeq2(base, 0, t2, t3, 0, 1)
				}
				return NewSeq(base+1, t2, t3, HoleCopy)
			} else if t2 == HoleCopy {
				return newSeq2(base, t1, 0, t3, 0, 2)
			}
			return newSeq2(base, t1, t2, 0, 0, 3)
		}
		if t1 < 0 {
			return newSeq2(base, getCopy(t1), t2, t3, 1, 1)
		} else if t2 < 0 {
			return newSeq2(base, t1, getCopy(t2), t3, 1, 2)
		}
		return newSeq2(base, t1, t2, getCopy(t3), 1, 3)
	case 0:
		return newSeq2(base, t1, t2, t3, 1, 0)
	}
	return 0
}

func (s Seq) each(f func(i tile.Instance) bool) bool {
	idx := s.OpenedIndex()
	base := s.Base()
	isComplete := s.IsComplete()
	copies := s.copies()
	openedTile := tile.Tile(idx-1) + base
	end := base + 3
	for base < end {
		c := tile.CopyID(copies & 3)
		copies >>= 2
		if isComplete || openedTile != base {
			if !f(base.Instance(c)) {
				return false
			}
		}
		base++
	}
	return true
}

func (s Seq) IsBadWait() bool {
	return s.Waits().Count() == 1
}

func (s Seq) Open(t tile.Instance, opponent base.Opponent) Meld {
	if opponent != base.Left {
		return 0
	}
	base := s.Base()
	ot := t.Tile()
	if !s.OpenedBy().Check(ot) {
		return 0
	}
	mask := s.copies()
	switch ot {
	case base - 1:
		mask = (mask << 2) | seqCopies(t.CopyID())
		return newSeq(base-1, mask, 1, 1).Meld()
	case base:
		mask = (mask & (^3)) | seqCopies(t.CopyID())
		return newSeq(base, mask, 1, 1).Meld()
	case base + 1:
		return newSeq2(base, mask.at(0), t.CopyID(), mask.at(2), 1, 2).Meld()
	case base + 2:
		mask = (mask & ^(3 << 4)) | (seqCopies(t.CopyID()) << 4)
		return newSeq(base, mask, 1, 3).Meld()
	}
	return 0
}

func (s Seq) setCopies(mask seqCopies) Seq {
	x := int(mask) & 63
	x <<= 7

	return Seq((int(s) & (^seqCopiesMask)) | x)
}

func (s Seq) copies() seqCopies {
	return seqCopies((s >> 7) & 63)
}

func (s Seq) Base() tile.Tile {
	return tile.Tile((s >> 2) & 31)
}

func (s Seq) Meld() Meld {
	return Meld(s)
}

func (s Seq) Opponent() base.Opponent {
	return base.Left
}

func (s Seq) IsComplete() bool {
	return ((s >> (2 + 5 + 3*2)) & 1) == 1
}

func (s Seq) OpenedIndex() int {
	return int((s >> (2 + 5 + 3*2 + 1)) & 3)
}

func (s Seq) IsOpened() bool {
	return s.IsComplete() && s.OpenedIndex() != 0
}

func (s Seq) OriginalWaits() compact.Tiles {
	idx := s.OpenedIndex()
	base := s.Base()
	switch idx {
	case 2:
		// hole waits
		return compact.FromTile(base + 1)
	case 1:
		if base.Number() == 7 {
			return compact.FromTile(base)
		}
		return compact.FromTiles(base, base+3)
	case 3:
		if base.Number() == 1 {
			return compact.FromTile(base + 2)
		}
		return compact.FromTiles(base-1, base+2)
	}
	// Complete chi closed in hand = 0
	return 0
}

func (s Seq) OpenedBy() compact.Tiles {
	return s.Waits()
}

func (s Seq) Waits() compact.Tiles {
	if s.IsComplete() {
		return 0
	}
	return s.OriginalWaits()
}

func (s Seq) Rebase(in compact.Instances) Meld {
	meld := s
	mask := seqCopies(0)
	base := s.Base()
	end := base + 3
	excluded := tile.Tile(-1)
	if !s.IsComplete() {
		excluded = tile.Tile(s.OpenedIndex()-1) + base
	}
	shift := uint(0)
	for base < end {
		if excluded != base {
			first := in.GetMask(base).FirstCopy()
			if first == tile.NullCopy {
				return 0
			}
			positionMask := seqCopies(first) << shift
			mask |= positionMask
		}
		base++
		shift += 2
	}
	return meld.setCopies(mask).Meld()
}

func (s Seq) Instances() (ret tile.Instances) {
	s.each(func(i tile.Instance) bool {
		ret = append(ret, i)
		return true
	})
	return
}

func (s Seq) AddTo(in compact.Instances) {
	s.each(func(i tile.Instance) bool {
		in.Set(i)
		return true
	})
}

func (s Seq) ExtractFrom(in compact.Instances) bool {
	return s.each(func(i tile.Instance) bool {
		return in.Remove(i)
	})
}
