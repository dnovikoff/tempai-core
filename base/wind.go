package base

import (
	"github.com/dnovikoff/tempai-core/tile"
)

//go:generate stringer -type=Wind
type Wind int

const (
	WindEast Wind = iota
	WindSouth
	WindWest
	WindNorth
	WindEnd
)

func (w Wind) tile() tile.Tile {
	return tile.Tile(w) + tile.East
}

func (w Wind) CheckTile(t tile.Tile) bool {
	return w.tile() == t
}

func (w Wind) fix() Wind {
	if w < 0 {
		return (w + (4 * (w / -4)) + 4) % 4
	}
	return w % 4
}

func (w Wind) Advance(num int) Wind {
	return (w + Wind(num)).fix()
}
