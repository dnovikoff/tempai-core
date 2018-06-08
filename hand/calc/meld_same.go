package calc

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/tile"
)

type same struct {
	meld
	count int
}

type tanki struct {
	same
}

type ponPart struct {
	same
}

func Kan(t tile.Tile) Meld {
	return &same{meld{
		t:     t,
		tags:  TagPon | TagKan | TagComplete | addTags(t),
		tiles: tile.Tiles{t, t, t, t},
	}, 4}
}

func Pon(t tile.Tile) Meld {
	return &same{meld{
		tags:  TagComplete | TagPon | addTags(t),
		t:     t,
		tiles: tile.Tiles{t, t, t},
	}, 3}
}

func PonPart(t tile.Tile) Meld {
	return &ponPart{same{meld{
		tags:   TagPon,
		waits:  tile.Tiles{t},
		cwaits: compact.FromTile(t),
		t:      t,
		tiles:  tile.Tiles{t, t},
	}, 2}}
}

func Pair(t tile.Tile) Meld {
	return &same{meld{
		tags:  TagPair | TagComplete | addTags(t),
		t:     t,
		tiles: tile.Tiles{t, t},
	}, 2}
}

func Tanki(t tile.Tile) Meld {
	return &tanki{same{meld{
		tags:   TagPair | TagTanki,
		t:      t,
		waits:  tile.Tiles{t},
		cwaits: compact.FromTile(t),
		tiles:  tile.Tiles{t},
	}, 1}}
}

func (m *tanki) Complete(t tile.Tile) Meld {
	if m.t != t {
		return nil
	}
	return Pair(t)
}

func (m *ponPart) Complete(t tile.Tile) Meld {
	if m.t != t {
		return nil
	}
	return Pon(t)
}

func (m *same) Extract(x Counters) bool {
	return x.Dec(m.t, m.count)
}

func addTags(t tile.Tile) Tags {
	if compact.Honor.Check(t) {
		return TagHonor
	}
	if compact.Terminal.Check(t) {
		return TagTerminal
	}
	return TagMiddle
}
