package calc

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/tile"
)

type Meld interface {
	Tile() tile.Tile
	Tags() Tags
	Waits() tile.Tiles
	CompactWaits() compact.Tiles
	// For debug
	Tiles() tile.Tiles
	Complete(t tile.Tile) Meld
	Extract(Counters) bool
}

type meld struct {
	tags   Tags
	waits  tile.Tiles
	cwaits compact.Tiles
	t      tile.Tile
	tiles  tile.Tiles
}

func (m *meld) Tags() Tags {
	return m.tags
}

func (m *meld) Complete(t tile.Tile) Meld {
	return nil
}

func (m *meld) Tile() tile.Tile {
	return m.t
}

func (m *meld) Waits() tile.Tiles {
	return m.waits
}

func (m *meld) Tiles() tile.Tiles {
	return m.tiles
}

func (m *meld) CompactWaits() compact.Tiles {
	return m.cwaits
}

type openWrapper struct {
	Meld
}

func (w openWrapper) Tags() Tags {
	return w.Meld.Tags() | TagOpened
}

func Open(m Meld) Meld {
	return openWrapper{m}
}
