package calc

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/tile"
)

type chi struct{ meld }

// Ryanman and penchan
type chiPart1 struct{ meld }

// Kanchan
type chiPart2 struct{ meld }

func Chi(t tile.Tile) Meld {
	if !compact.Sequence.Check(t) || t.Number() > 7 {
		return nil
	}
	tags := TagComplete | TagChi
	switch t.Number() {
	case 1, 7:
		tags |= TagTerminal
	default:
		tags |= TagMiddle
	}
	return &chi{meld{
		tags:  tags,
		t:     t,
		tiles: tile.Tiles{t, t + 1, t + 2},
	}}
}

func (m *chi) Extract(x Counters) bool {
	return x.Dec3(m.t)
}

// Penchan or ryanman
func ChiPart1(t tile.Tile) Meld {
	if !compact.Sequence.Check(t) || t.Number() > 8 {
		return nil
	}
	tags := TagChi
	var waits tile.Tiles
	switch t.Number() {
	case 1:
		waits = tile.Tiles{t + 2}
		tags |= TagPenchan
	case 8:
		waits = tile.Tiles{t - 1}
		tags |= TagPenchan
	default:
		tags |= TagRyanman
		waits = tile.Tiles{t - 1, t + 2}
	}
	return &chiPart1{meld{
		tags:   tags,
		waits:  waits,
		cwaits: compact.FromTiles(waits...),
		t:      t,
		tiles:  tile.Tiles{t, t + 1},
	}}
}

// Kanchan
func ChiPart2(t tile.Tile) Meld {
	if !compact.Sequence.Check(t) || t.Number() > 7 {
		return nil
	}
	return &chiPart2{meld{
		tags:   TagChi | TagKanchan,
		waits:  tile.Tiles{t + 1},
		cwaits: compact.FromTile(t + 1),
		t:      t,
		tiles:  tile.Tiles{t, t + 2},
	}}
}

func (m *chiPart2) Extract(x Counters) bool {
	return x.Dec2(m.t, 2)
}

func (m *chiPart1) Extract(x Counters) bool {
	return x.Dec2(m.t, 1)
}

func (m *chiPart1) Complete(t tile.Tile) Meld {
	if !m.cwaits.Check(t) {
		return nil
	}
	if t < m.t {
		return Chi(t)
	}
	return Chi(m.t)
}

func (m *chiPart2) Complete(t tile.Tile) Meld {
	if m.t != t-1 {
		return nil
	}
	return Chi(m.t)
}
