package meld

import (
	"strings"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/tile"
)

// 2 - Type

type Meld int
type Melds []Meld

type ReducedInterface interface {
	Meld() Meld
	Opponent() base.Opponent
	IsOpened() bool
	IsBadWait() bool
	OriginalWaits() compact.Tiles
	OpenedBy() compact.Tiles
	Open(tile.Instance, base.Opponent) Meld
}

type Interface interface {
	ReducedInterface

	Base() tile.Tile
	Rebase(compact.Instances) Meld
	Waits() compact.Tiles
	IsComplete() bool
	AddTo(compact.Instances)
	ExtractFrom(compact.Instances) bool
}

func (this Meld) IsPon() bool {
	return this.Type() == TypeSame
}

func (this Meld) IsSeq() bool {
	return this.Type() == TypeSeq
}

func (this Meld) IsTanki() bool {
	if this.Type() != TypePair {
		return false
	}
	return Pair(this).IsTanki()
}

func (this Meld) IsKan() bool {
	if this.Type() != TypeSame {
		return false
	}
	return Same(this).IsKan()
}

func (this Meld) Type() Type {
	return Type(this & 3)
}

func (this Meld) IsNull() bool {
	return this == 0
}

func (this Meld) Interface() Interface {
	switch this.Type() {
	case TypeSeq:
		return Seq(this)
	case TypeSame:
		return Same(this)
	case TypePair:
		return Pair(this)
	}
	return nil
}

func (this Meld) Waits() compact.Tiles {
	switch this.Type() {
	case TypeSeq:
		return Seq(this).Waits()
	case TypeSame:
		return Same(this).Waits()
	case TypePair:
		return Pair(this).Waits()
	}
	return 0
}

func (this Meld) IsComplete() bool {
	switch this.Type() {
	case TypeSeq:
		return Seq(this).IsComplete()
	case TypeSame:
		return Same(this).IsComplete()
	case TypePair:
		return Pair(this).IsComplete()
	}
	return false
}

func (this Meld) AddTo(in compact.Instances) {
	switch this.Type() {
	case TypeSeq:
		Seq(this).AddTo(in)
	case TypeSame:
		Same(this).AddTo(in)
	case TypePair:
		Pair(this).AddTo(in)
	}
}

func (this Meld) RebaseAndExtractFrom(in compact.Instances) Meld {
	fixed := this.Rebase(in)
	if fixed == 0 {
		return 0
	}
	fixed.ExtractFrom(in)
	return fixed
}

func (this Meld) ExtractFrom(in compact.Instances) bool {
	switch this.Type() {
	case TypeSeq:
		return Seq(this).ExtractFrom(in)
	case TypeSame:
		return Same(this).ExtractFrom(in)
	case TypePair:
		return Pair(this).ExtractFrom(in)
	}
	return false
}

func (this Meld) Rebase(in compact.Instances) Meld {
	switch this.Type() {
	case TypeSeq:
		return Seq(this).Rebase(in)
	case TypeSame:
		return Same(this).Rebase(in)
	case TypePair:
		return Pair(this).Rebase(in)
	}
	return 0
}

func (this Meld) Instances() tile.Instances {
	i := compact.NewInstances()
	this.AddTo(i)
	return i.Instances()
}

func (this Meld) Base() tile.Tile {
	switch this.Type() {
	case TypeSeq:
		return Seq(this).Base()
	case TypeSame:
		return Same(this).Base()
	case TypePair:
		return Pair(this).Base()
	}
	return 0
}

// Only for 1,2 tiles left
func ExtractLastMeld(t compact.Instances) Meld {
	i := t.Instances()
	switch len(i) {
	case 1:
		return NewTanki(i[0]).Meld()
	case 2:
		i1 := i[0]
		i2 := i[1]
		if i1 > i2 {
			i1, i2 = i2, i1
		}
		t1 := i1.Tile()
		t2 := i2.Tile()
		switch t2 - t1 {
		case 0:
			return NewPonPartFromExisting(t1, i1.CopyId(), i2.CopyId()).Meld()
		case 1:
			return NewSeq(t1, i1.CopyId(), i2.CopyId(), HoleCopy).Meld()
		case 2:
			return NewSeq(t1, i1.CopyId(), HoleCopy, i2.CopyId()).Meld()
		}

	}
	return 0
}

func (this Melds) Win(t tile.Tile) Meld {
	for _, m := range this {
		w := m.Waits()
		if w == 0 {
			continue
		}
		if w.Check(t) {
			return m
		}
	}
	return 0
}

func MeldsToString(melds Melds) string {
	x := make([]string, len(melds))
	for k, v := range melds {
		c := compact.NewInstances()
		v.Interface().AddTo(c)
		x[k] = c.Instances().String()
	}
	return strings.Join(x, " ")
}

func (this Melds) AddTo(x compact.Instances) {
	for _, v := range this {
		v.AddTo(x)
	}
}
