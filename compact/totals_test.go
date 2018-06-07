package compact

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dnovikoff/tempai-core/tile"
)

func TestTotalsAdd(t *testing.T) {
	tl := tile.Man1

	x := NewTotals()
	assert.Equal(t, 0, x.Get(tl))
	assert.False(t, x.IsFull(tl))
	assert.Equal(t, AllTiles, x.FreeTiles())
	assert.Equal(t, Tiles(0), x.FullTiles())
	assert.Equal(t, 0, x.Count())

	x.Add(tl, 1)
	assert.Equal(t, 1, x.Get(tl))
	assert.False(t, x.IsFull(tl))
	assert.Equal(t, AllTiles, x.FreeTiles())
	assert.Equal(t, Tiles(0), x.FullTiles())
	assert.Equal(t, 1, x.Count())

	x.Add(tl, 2)
	assert.Equal(t, 3, x.Get(tl))
	assert.False(t, x.IsFull(tl))
	assert.Equal(t, AllTiles, x.FreeTiles())
	assert.Equal(t, Tiles(0), x.FullTiles())
	assert.Equal(t, 3, x.Count())

	x.Add(tl, 1)
	assert.Equal(t, 4, x.Get(tl))
	assert.True(t, x.IsFull(tl))
	assert.Equal(t, AllTiles.Unset(tl), x.FreeTiles())
	assert.Equal(t, FromTile(tl), x.FullTiles())
	assert.Equal(t, 4, x.Count())

	x.Add(tl, 1)
	assert.Equal(t, 5, x.Get(tl))
	assert.True(t, x.IsFull(tl))
	assert.Equal(t, AllTiles.Unset(tl), x.FreeTiles())
	assert.Equal(t, FromTile(tl), x.FullTiles())
	assert.Equal(t, 5, x.Count())
}

func TestTotalsTiles(t *testing.T) {
	x := NewTotals()
	x.Set(tile.Pin1, 1)
	x.Set(tile.Red, 3)
	x.Set(tile.West, 2)
	assert.Equal(t, NewTiles(
		tile.Pin1,
		tile.Red,
		tile.West,
	), x.UniqueTiles())
	assert.Equal(t, 3, x.UniqueCount())
	assert.Equal(t, 6, x.Count())
}

func TestTotalsMerge(t *testing.T) {
	x1 := NewTotals()
	x1.Add(tile.Pin1, 1)
	x1.Add(tile.Red, 3)
	x1.Add(tile.West, 2)

	x2 := NewInstances()
	x2.SetCount(tile.Man1, 1)
	x2.SetCount(tile.Pin1, 1)
	x2.SetCount(tile.White, 2)
	x2.SetCount(tile.West, 1)

	x1.Merge(x2).Merge(nil)
	assert.Equal(t, NewTiles(
		tile.Man1,
		tile.Pin1,
		tile.Red,
		tile.West,
		tile.White,
	), x1.UniqueTiles())

	assert.Equal(t, 1, x1.Get(tile.Man1))
	assert.Equal(t, 2, x1.Get(tile.Pin1))
	assert.Equal(t, 3, x1.Get(tile.Red))
	assert.Equal(t, 3, x1.Get(tile.West))
	assert.Equal(t, 2, x1.Get(tile.White))
}

func TestTotalsTilesClone(t *testing.T) {
	x1 := NewTotals()
	x1.Set(tile.Man1, 1)
	x2 := x1.Clone()
	assert.Equal(t, x1.UniqueTiles(), x2.UniqueTiles())
	x2.Set(tile.Man2, 1)
	assert.NotEqual(t, x1.UniqueTiles(), x2.UniqueTiles())
}
