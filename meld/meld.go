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

func (m Meld) IsPon() bool {
	return m.Type() == TypeSame
}

func (m Meld) IsSeq() bool {
	return m.Type() == TypeSeq
}

func (m Meld) IsTanki() bool {
	if m.Type() != TypePair {
		return false
	}
	return Pair(m).IsTanki()
}

func (m Meld) IsKan() bool {
	if m.Type() != TypeSame {
		return false
	}
	return Same(m).IsKan()
}

func (m Meld) Type() Type {
	return Type(m & 3)
}

func (m Meld) IsNull() bool {
	return m == 0
}

func (m Meld) Interface() Interface {
	switch m.Type() {
	case TypeSeq:
		return Seq(m)
	case TypeSame:
		return Same(m)
	case TypePair:
		return Pair(m)
	}
	return nil
}

func (m Meld) Waits() compact.Tiles {
	switch m.Type() {
	case TypeSeq:
		return Seq(m).Waits()
	case TypeSame:
		return Same(m).Waits()
	case TypePair:
		return Pair(m).Waits()
	}
	return 0
}

func (m Meld) IsComplete() bool {
	switch m.Type() {
	case TypeSeq:
		return Seq(m).IsComplete()
	case TypeSame:
		return Same(m).IsComplete()
	case TypePair:
		return Pair(m).IsComplete()
	}
	return false
}

func (m Meld) AddTo(in compact.Instances) {
	switch m.Type() {
	case TypeSeq:
		Seq(m).AddTo(in)
	case TypeSame:
		Same(m).AddTo(in)
	case TypePair:
		Pair(m).AddTo(in)
	}
}

func (m Meld) RebaseAndExtractFrom(in compact.Instances) Meld {
	fixed := m.Rebase(in)
	if fixed == 0 {
		return 0
	}
	fixed.ExtractFrom(in)
	return fixed
}

func (m Meld) ExtractFrom(in compact.Instances) bool {
	switch m.Type() {
	case TypeSeq:
		return Seq(m).ExtractFrom(in)
	case TypeSame:
		return Same(m).ExtractFrom(in)
	case TypePair:
		return Pair(m).ExtractFrom(in)
	}
	return false
}

func (m Meld) Rebase(in compact.Instances) Meld {
	switch m.Type() {
	case TypeSeq:
		return Seq(m).Rebase(in)
	case TypeSame:
		return Same(m).Rebase(in)
	case TypePair:
		return Pair(m).Rebase(in)
	}
	return 0
}

func (m Meld) Instances() tile.Instances {
	i := compact.NewInstances()
	m.AddTo(i)
	return i.Instances()
}

func (m Meld) Base() tile.Tile {
	switch m.Type() {
	case TypeSeq:
		return Seq(m).Base()
	case TypeSame:
		return Same(m).Base()
	case TypePair:
		return Pair(m).Base()
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
			return NewPonPartFromExisting(t1, i1.CopyID(), i2.CopyID()).Meld()
		case 1:
			return NewSeq(t1, i1.CopyID(), i2.CopyID(), HoleCopy).Meld()
		case 2:
			return NewSeq(t1, i1.CopyID(), HoleCopy, i2.CopyID()).Meld()
		}

	}
	return 0
}

func (m Melds) Win(t tile.Tile) Meld {
	for _, m := range m {
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

func (m Melds) AddTo(x compact.Instances) {
	for _, v := range m {
		v.AddTo(x)
	}
}
