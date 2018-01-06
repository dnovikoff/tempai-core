package meld

import (
	"bitbucket.org/dnovikoff/tempai-core/base"
	"bitbucket.org/dnovikoff/tempai-core/compact"
	"bitbucket.org/dnovikoff/tempai-core/tile"
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
	return newPair(t1.Tile(), pairSubtypeTanki, t1.CopyId(), 0)
}

func NewOne(t1 tile.Instance) Pair {
	return newPair(t1.Tile(), pairSubtypeOne, t1.CopyId(), 0)
}

func NewPair(t tile.Tile, c1, c2 tile.CopyId) Pair {
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

func newPair(base tile.Tile, sub pairSubtype, c1, c2 tile.CopyId) Pair {
	x := int(c2) & 3
	x = (x << 2) | (int(c1) & 3)
	x = (x << 2) | (int(sub) & 3)
	x = (x << 6) | (int(base) & 63)
	x = (x << 2) | int(TypePair)
	return Pair(x)
}

func (this Pair) IsBadWait() bool {
	return this.subType() == pairSubtypeTanki
}

func (this Pair) Base() tile.Tile {
	return tile.Tile((this >> 2) & 63)
}

func (this Pair) subType() pairSubtype {
	return pairSubtype((this >> (2 + 6)) & 3)
}

func (this Pair) c1() tile.CopyId {
	return tile.CopyId((this >> (2 + 6 + 2)) & 3)
}

func (this Pair) c2() tile.CopyId {
	return tile.CopyId((this >> (2 + 6 + 4)) & 3)
}

func (this Pair) IsComplete() bool {
	return this.subType() == pairSubtypePair || this.subType() == pairSubtypeOne
}

func (this Pair) IsOpened() bool {
	return false
}

func (this Pair) IsTanki() bool {
	return this.subType() == pairSubtypeTanki
}

func (this Pair) OriginalWaits() compact.Tiles {
	if this.IsComplete() {
		return 0
	}
	return compact.NewFromTile(this.Base())
}

func (this Pair) Opponent() base.Opponent {
	return base.Self
}

func (this Pair) Waits() compact.Tiles {
	return this.OriginalWaits()
}

func (this Pair) Open(tile.Instance, base.Opponent) Meld {
	return 0
}

func (this Pair) OpenedBy() compact.Tiles {
	return 0
}

func (this Pair) Meld() Meld {
	return Meld(this)
}

func (this Pair) Rebase(in compact.Instances) Meld {
	meld := this
	mask := in.GetMask(this.Base())
	switch this.subType() {
	case pairSubtypeTanki:
		return NewTanki(mask.First()).Meld()
	case pairSubtypeOne:
		return NewOne(mask.First()).Meld()
	case pairSubtypePair:
		first := mask.FirstCopy()
		mask = mask.UnsetCopyBit(first)
		second := mask.FirstCopy()
		return NewPair(this.Base(), first, second).Meld()
	}
	return meld.Meld()
}

func (this Pair) AddTo(in compact.Instances) {
	switch this.subType() {
	case pairSubtypeTanki, pairSubtypeOne:
		in.Set(this.Base().Instance(this.c1()))
	case pairSubtypePair:
		mask := in.GetMask(this.Base())
		mask = mask.SetCopyBit(this.c1()).SetCopyBit(this.c2())
		in.SetMask(mask)
	}
}

func (this Pair) ExtractFrom(in compact.Instances) bool {
	switch this.subType() {
	case pairSubtypeTanki, pairSubtypeOne:
		return in.Remove(this.Base().Instance(this.c1()))
	case pairSubtypePair:
		mask := in.GetMask(this.Base())
		next := mask.UnsetCopyBit(this.c1()).UnsetCopyBit(this.c2())
		in.SetMask(next)
		return next != mask
	}
	return false
}
