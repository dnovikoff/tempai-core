package meld

import (
	"bitbucket.org/dnovikoff/tempai-core/base"
	"bitbucket.org/dnovikoff/tempai-core/compact"
	"bitbucket.org/dnovikoff/tempai-core/tile"
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

func newSame(base tile.Tile, subType sameSubtype, c1, c2 tile.CopyId, op base.Opponent) Same {
	x := int(op) & 3
	x = (x << 2) | (int(c2) & 3)
	x = (x << 2) | (int(c1) & 3)
	x = (x << 2) | (int(subType) & 3)
	x = (x << 6) | (int(base) & 63)
	x = (x << 2) | int(TypeSame)
	return Same(x)
}

func (this Same) Base() tile.Tile {
	return tile.Tile((this >> 2) & 63)
}

func (this Same) subType() sameSubtype {
	return sameSubtype((this >> (2 + 6)) & 3)
}

func (this Same) c1() tile.CopyId {
	return tile.CopyId((this >> (2 + 6 + 2)) & 3)
}

func (this Same) c2() tile.CopyId {
	return tile.CopyId((this >> (2 + 6 + 2 + 2)) & 3)
}

func (this Same) Opponent() base.Opponent {
	return base.Opponent((this >> (2 + 6 + 2 + 2 + 2)) & 3)
}

func (this Same) baseCompact() compact.Tiles {
	return compact.NewFromTile(this.Base())
}

func (this Same) IsComplete() bool {
	return this.subType() != samePart || this.Opponent() != base.Self
}

func (this Same) Meld() Meld {
	return Meld(this)
}

func (this Same) IsOpened() bool {
	return this.Opponent() != base.Self
}

func (this Same) IsBadWait() bool {
	return false
}

func (this Same) OpenedBy() compact.Tiles {
	switch this.subType() {
	case samePart:
		if this.Opponent() == base.Self {
			return this.baseCompact()
		}
	case samePon:
		return this.baseCompact()
	}
	return 0
}

func (this Same) IsUpgraded() bool {
	return this.subType() == sameUpgraded
}

func (this Same) UpgradeOrAnyInstance() tile.Instance {
	if !this.IsUpgraded() {
		return this.Base().Instance(0)
	}
	return this.Base().Instance(this.c2())
}

func (this Same) UpgradeInstance() tile.Instance {
	if !this.IsUpgraded() {
		return tile.InstanceNull
	}
	return this.Base().Instance(this.c2())
}

func (this Same) IsKan() bool {
	switch this.subType() {
	case sameKan, sameUpgraded:
		return true
	}
	return false
}

func (this Same) Upgrade() Same {
	if this.subType() != samePart {
		return 0
	}
	op := this.Opponent()
	if op == base.Self {
		return 0
	}
	return newSame(this.Base(), sameUpgraded, this.c1(), this.c2(), op)
}

func (this Same) OpenedCopy() tile.CopyId {
	return this.c1()
}

func (this Same) NotInPonCopy() tile.CopyId {
	return this.c2()
}

func (this Same) Open(t tile.Instance, opponent base.Opponent) Meld {
	if !this.OpenedBy().Check(t.Tile()) {
		return 0
	}
	c := t.CopyId()
	base := this.Base()
	c1 := this.c1()
	c2 := this.c2()
	switch this.subType() {
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

func (this Same) InstancesMask() compact.Mask {
	b := this.Base()
	kan := compact.NewMaskByCount(4, b)
	switch this.subType() {
	case samePart:
		kan = kan.UnsetCopyBit(this.c2())
		if this.Opponent() == base.Self {
			kan = kan.UnsetCopyBit(this.c1())
		}
	case samePon:
		kan = kan.UnsetCopyBit(this.c1())
	case sameKan, sameUpgraded:
	default:
		return 0
	}
	return kan
}

func (this Same) Waits() compact.Tiles {
	if this.subType() == samePart && this.Opponent() == base.Self {
		return this.baseCompact()
	}
	return 0
}

func (this Same) OriginalWaits() compact.Tiles {
	switch this.subType() {
	case samePart, sameUpgraded:
		return this.baseCompact()
	}
	return 0
}

// Closed kan does not need to have opened param
// You can use 0 in any case
// Still it is required for full compat with tenhou kans
func NewKan(t tile.Tile, opened tile.CopyId) Same {
	return newSame(t, sameKan, opened, 0, base.Self)
}

func NewKanOpened(t tile.Tile, opened tile.CopyId, opponent base.Opponent) Same {
	return newSame(t, sameKan, opened, 0, opponent)
}

func NewKanUpgraded(t tile.Tile, opened, upgraded tile.CopyId, opponent base.Opponent) Same {
	return newSame(t, sameUpgraded, opened, upgraded, opponent)
}

func NewPon(t tile.Tile, notInPon tile.CopyId) Same {
	return newSame(t, samePon, notInPon, 0, base.Self)
}

func NewPonOpened(t tile.Tile, opened, notInPon tile.CopyId, opponent base.Opponent) Same {
	return newSame(t, samePart, opened, notInPon, opponent)
}

func NewPonPart(t tile.Tile, c1, c2 tile.CopyId) Same {
	if c1 > c2 {
		c1, c2 = c2, c1
	}
	return newSame(t, samePart, c1, c2, base.Self)
}

func NewPonPartFromExisting(t tile.Tile, c1, c2 tile.CopyId) Same {
	mask := compact.NewMask(0, t).SetCopyBit(c1).SetCopyBit(c2).InvertTiles()
	c1 = mask.FirstCopy()
	c2 = mask.UnsetCopyBit(c1).FirstCopy()
	return NewPonPart(t, c1, c2)
}

func (this Same) Rebase(in compact.Instances) Meld {
	original := in.GetMask(this.Base())
	switch this.subType() {
	case samePart:
		if original.Count() < 2 {
			return 0
		}
		inverted := original.InvertTiles()
		first := inverted.FirstCopy()
		second := inverted.UnsetCopyBit(first).FirstCopy()
		return NewPonPart(this.Base(), first, second).Meld()
	case sameKan, sameUpgraded:
		if original.Count() != 4 {
			return 0
		}
		return this.Meld()
	case samePon:
		if original.Count() < 3 {
			return 0
		}
		notInPon := original.InvertTiles().FirstCopy()
		op := this.Opponent()
		if op == base.Self {
			return NewPon(this.Base(), notInPon).Meld()
		}
		// Not enough information to choose one
		opened := original.FirstCopy()
		return NewPonOpened(this.Base(), opened, notInPon, op).Meld()
	}
	return 0
}

func (this Same) AddTo(in compact.Instances) {
	original := in.GetMask(this.Base())
	mask := this.InstancesMask()
	in.SetMask(original.Merge(mask))
}

func (this Same) ExtractFrom(in compact.Instances) bool {
	original := in.GetMask(this.Base())
	mask := this.InstancesMask()
	next := original.Remove(mask)
	in.SetMask(next)
	return next != original
}
