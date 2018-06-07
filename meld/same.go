package meld

import (
	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/tile"
)

// Opponent for Seq is always Left
// Type - 2 | Base - 6 | Subtype - 2 | c1 -2 | c2 - 2 | op -2
// Base - tile
// Subtype
// 1. part c1 - 1 | c2 - 2
// 2.
// 3.
// 4.
// samePonOpenedX - specialEncoded to fit 16

// Part
// c1 - not in pon/opened
// c2 - opened

// Pon c1 - tile not in base Opponent
// c1 = not in base
// c2 = unused

// Kan c1 - opened
// Upgraded
// c1 = opened
// c2 = upgraded

type Same Meld

var _ Interface = Same(0)

type sameSubtype int

const (
	samePart sameSubtype = iota
	samePon
	sameKan
	sameUpgraded
)

func newSame(base tile.Tile, subType sameSubtype, c1, c2 tile.CopyID, op base.Opponent) Same {
	x := int(op) & 3
	x = (x << 2) | (int(c2) & 3)
	x = (x << 2) | (int(c1) & 3)
	x = (x << 2) | (int(subType) & 3)
	x = (x << 6) | (int(base) & 63)
	x = (x << 2) | int(TypeSame)
	return Same(x)
}

func (s Same) Base() tile.Tile {
	return tile.Tile((s >> 2) & 63)
}

func (s Same) subType() sameSubtype {
	return sameSubtype((s >> (2 + 6)) & 3)
}

func (s Same) c1() tile.CopyID {
	return tile.CopyID((s >> (2 + 6 + 2)) & 3)
}

func (s Same) c2() tile.CopyID {
	return tile.CopyID((s >> (2 + 6 + 2 + 2)) & 3)
}

func (s Same) Opponent() base.Opponent {
	return base.Opponent((s >> (2 + 6 + 2 + 2 + 2)) & 3)
}

func (s Same) baseCompact() compact.Tiles {
	return compact.FromTile(s.Base())
}

func (s Same) IsComplete() bool {
	return s.subType() != samePart || s.Opponent() != base.Self
}

func (s Same) Meld() Meld {
	return Meld(s)
}

func (s Same) IsOpened() bool {
	return s.Opponent() != base.Self
}

func (s Same) IsBadWait() bool {
	return false
}

func (s Same) OpenedBy() compact.Tiles {
	switch s.subType() {
	case samePart:
		if s.Opponent() == base.Self {
			return s.baseCompact()
		}
	case samePon:
		return s.baseCompact()
	}
	return 0
}

func (s Same) IsUpgraded() bool {
	return s.subType() == sameUpgraded
}

func (s Same) UpgradeOrAnyInstance() tile.Instance {
	if !s.IsUpgraded() {
		return s.Base().Instance(0)
	}
	return s.Base().Instance(s.c2())
}

func (s Same) UpgradeInstance() tile.Instance {
	if !s.IsUpgraded() {
		return tile.InstanceNull
	}
	return s.Base().Instance(s.c2())
}

func (s Same) IsKan() bool {
	switch s.subType() {
	case sameKan, sameUpgraded:
		return true
	}
	return false
}

func (s Same) Upgrade() Same {
	if s.subType() != samePart {
		return 0
	}
	op := s.Opponent()
	if op == base.Self {
		return 0
	}
	return newSame(s.Base(), sameUpgraded, s.c1(), s.c2(), op)
}

func (s Same) OpenedCopy() tile.CopyID {
	return s.c1()
}

func (s Same) NotInPonCopy() tile.CopyID {
	return s.c2()
}

func (s Same) Open(t tile.Instance, opponent base.Opponent) Meld {
	if !s.OpenedBy().Check(t.Tile()) {
		return 0
	}
	c := t.CopyID()
	base := s.Base()
	c1 := s.c1()
	c2 := s.c2()
	switch s.subType() {
	case samePart:
		if c == c1 {
			return newSame(base, samePart, c1, c2, opponent).Meld()
		} else if c == c2 {
			return newSame(base, samePart, c2, c1, opponent).Meld()
		}
	case samePon:
		if c1 == c {
			return newSame(base, sameKan, c1, 0, opponent).Meld()
		}
	}
	return 0
}

func (s Same) InstancesMask() compact.Mask {
	b := s.Base()
	kan := compact.NewMask(compact.MaskByCount(4), b)
	switch s.subType() {
	case samePart:
		kan = kan.UnsetCopyBit(s.c2())
		if s.Opponent() == base.Self {
			kan = kan.UnsetCopyBit(s.c1())
		}
	case samePon:
		kan = kan.UnsetCopyBit(s.c1())
	case sameKan, sameUpgraded:
	default:
		return 0
	}
	return kan
}

func (s Same) Waits() compact.Tiles {
	if s.subType() == samePart && s.Opponent() == base.Self {
		return s.baseCompact()
	}
	return 0
}

func (s Same) OriginalWaits() compact.Tiles {
	switch s.subType() {
	case samePart, sameUpgraded:
		return s.baseCompact()
	}
	return 0
}

// Closed kan does not need to have opened param
// You can use 0 in any case
// Still it is required for full compat with tenhou kans
func NewKan(t tile.Instance) Same {
	return NewKanOpened(t, base.Self)
}

func NewKanOpened(opened tile.Instance, opponent base.Opponent) Same {
	return newSame(opened.Tile(), sameKan, opened.CopyID(), 0, opponent)
}

func NewKanUpgraded(opened tile.Instance, upgraded tile.CopyID, opponent base.Opponent) Same {
	return newSame(opened.Tile(), sameUpgraded, opened.CopyID(), upgraded, opponent)
}

func NewPon(notInPon tile.Instance) Same {
	return newSame(notInPon.Tile(), samePon, notInPon.CopyID(), 0, base.Self)
}

func NewPonOpened(opened tile.Instance, notInPon tile.CopyID, opponent base.Opponent) Same {
	return newSame(opened.Tile(), samePart, opened.CopyID(), notInPon, opponent)
}

func NewPonPart(t tile.Tile, c1, c2 tile.CopyID) Same {
	if c1 > c2 {
		c1, c2 = c2, c1
	}
	return newSame(t, samePart, c1, c2, base.Self)
}

func NewPonPartFromExisting(t tile.Tile, c1, c2 tile.CopyID) Same {
	mask := compact.NewMask(0, t).SetCopyBit(c1).SetCopyBit(c2).InvertTiles()
	c1 = mask.FirstCopy()
	c2 = mask.UnsetCopyBit(c1).FirstCopy()
	return NewPonPart(t, c1, c2)
}

func (s Same) Rebase(in compact.Instances) Meld {
	original := in.GetMask(s.Base())
	switch s.subType() {
	case samePart:
		if original.Count() < 2 {
			return 0
		}
		inverted := original.InvertTiles()
		first := inverted.FirstCopy()
		second := inverted.UnsetCopyBit(first).FirstCopy()
		return NewPonPart(s.Base(), first, second).Meld()
	case sameKan, sameUpgraded:
		if original.Count() != 4 {
			return 0
		}
		return s.Meld()
	case samePon:
		if original.Count() < 3 {
			return 0
		}
		// In case of all for in a hand - any will fit
		notInPon := tile.AnyCopy
		if !original.IsFull() {
			notInPon = original.InvertTiles().FirstCopy()
		}
		op := s.Opponent()
		if op == base.Self {
			return NewPon(s.Base().Instance(notInPon)).Meld()
		}
		// Not enough information to choose one
		opened := original.FirstCopy()
		return NewPonOpened(s.Base().Instance(opened), notInPon, op).Meld()
	}
	return 0
}

func (s Same) AddTo(in compact.Instances) {
	original := in.GetMask(s.Base())
	mask := s.InstancesMask()
	in.SetMask(original.Merge(mask))
}

func (s Same) ExtractFrom(in compact.Instances) bool {
	original := in.GetMask(s.Base())
	mask := s.InstancesMask()
	next := original.Remove(mask)
	in.SetMask(next)
	return next != original
}
