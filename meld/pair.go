package meld

import (
	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/tile"
)

// Type - 2 | Base - 6 | subtype -2 | c1 - 2 | c2 - 2
// =14

type Pair Meld

var _ Interface = Pair(0)

type pairSubtype int

const (
	pairSubtypePair pairSubtype = iota
	pairSubtypeTanki
	pairSubtypeHole
	// One tile for kokushi
	pairSubtypeOne
)

func NewTanki(t1 tile.Instance) Pair {
	return newPair(t1.Tile(), pairSubtypeTanki, t1.CopyID(), 0)
}

func NewOne(t1 tile.Instance) Pair {
	return newPair(t1.Tile(), pairSubtypeOne, t1.CopyID(), 0)
}

func NewPair(t tile.Tile, c1, c2 tile.CopyID) Pair {
	if c1 > c2 {
		c1, c2 = c2, c1
	} else if c1 == c2 {
		return 0
	}
	return newPair(t, pairSubtypePair, c1, c2)
}

func NewPairFromMask(m compact.Mask) Pair {
	if m.Count() < 2 {
		return 0
	}
	c1 := m.FirstCopy()
	c2 := m.UnsetCopyBit(c1).FirstCopy()
	return NewPair(m.Tile(), c1, c2)
}

func NewHole(base tile.Tile) Pair {
	return newPair(base, pairSubtypeHole, 0, 0)
}

func newPair(base tile.Tile, sub pairSubtype, c1, c2 tile.CopyID) Pair {
	x := int(c2) & 3
	x = (x << 2) | (int(c1) & 3)
	x = (x << 2) | (int(sub) & 3)
	x = (x << 6) | (int(base) & 63)
	x = (x << 2) | int(TypePair)
	return Pair(x)
}

func (p Pair) IsBadWait() bool {
	return p.subType() == pairSubtypeTanki
}

func (p Pair) Base() tile.Tile {
	return tile.Tile((p >> 2) & 63)
}

func (p Pair) subType() pairSubtype {
	return pairSubtype((p >> (2 + 6)) & 3)
}

func (p Pair) c1() tile.CopyID {
	return tile.CopyID((p >> (2 + 6 + 2)) & 3)
}

func (p Pair) c2() tile.CopyID {
	return tile.CopyID((p >> (2 + 6 + 4)) & 3)
}

func (p Pair) IsComplete() bool {
	return p.subType() == pairSubtypePair || p.subType() == pairSubtypeOne
}

func (p Pair) IsOpened() bool {
	return false
}

func (p Pair) IsTanki() bool {
	return p.subType() == pairSubtypeTanki
}

func (p Pair) OriginalWaits() compact.Tiles {
	if p.IsComplete() {
		return 0
	}
	return compact.FromTile(p.Base())
}

func (p Pair) Opponent() base.Opponent {
	return base.Self
}

func (p Pair) Waits() compact.Tiles {
	return p.OriginalWaits()
}

func (p Pair) Open(tile.Instance, base.Opponent) Meld {
	return 0
}

func (p Pair) OpenedBy() compact.Tiles {
	return 0
}

func (p Pair) Meld() Meld {
	return Meld(p)
}

func (p Pair) Rebase(in compact.Instances) Meld {
	meld := p
	mask := in.GetMask(p.Base())
	switch p.subType() {
	case pairSubtypeTanki:
		return NewTanki(mask.First()).Meld()
	case pairSubtypeOne:
		return NewOne(mask.First()).Meld()
	case pairSubtypePair:
		first := mask.FirstCopy()
		mask = mask.UnsetCopyBit(first)
		second := mask.FirstCopy()
		return NewPair(p.Base(), first, second).Meld()
	}
	return meld.Meld()
}

func (p Pair) AddTo(in compact.Instances) {
	switch p.subType() {
	case pairSubtypeTanki, pairSubtypeOne:
		in.Set(p.Base().Instance(p.c1()))
	case pairSubtypePair:
		mask := in.GetMask(p.Base())
		mask = mask.SetCopyBit(p.c1()).SetCopyBit(p.c2())
		in.SetMask(mask)
	}
}

func (p Pair) ExtractFrom(in compact.Instances) bool {
	switch p.subType() {
	case pairSubtypeTanki, pairSubtypeOne:
		return in.Remove(p.Base().Instance(p.c1()))
	case pairSubtypePair:
		mask := in.GetMask(p.Base())
		next := mask.UnsetCopyBit(p.c1()).UnsetCopyBit(p.c2())
		in.SetMask(next)
		return next != mask
	}
	return false
}
