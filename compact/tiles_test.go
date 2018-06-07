package compact

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dnovikoff/tempai-core/tile"
)

func TestTilesNew(t *testing.T) {
	assert.Equal(t, "47p", FromTiles(tile.Pin4, tile.Pin7).Tiles().String())
}

func TestCTS(t *testing.T) {
	x := Tiles(0)
	assert.True(t, x.IsEmpty())
	assert.Equal(t, 0, x.Count())

	x.Set(tile.Man1)
	// original not changed
	assert.True(t, x.IsEmpty())

	assert.False(t, x.Check(tile.Man1))
	x = x.Set(tile.Man1)
	assert.False(t, x.IsEmpty())
	assert.Equal(t, 1, x.Count())
	assert.True(t, x.Check(tile.Man1))
	assert.False(t, x.Check(tile.Man2))

	x = x.Set(tile.Man2).Set(tile.East)
	assert.Equal(t, "12m1z", x.Tiles().String())

	assert.Equal(t, 3, x.Count())
	x = x.Unset(tile.Man6)
	assert.Equal(t, "12m1z", x.Tiles().String())
	assert.Equal(t, 3, x.Count())
	x = x.Unset(tile.East)
	assert.Equal(t, "12m", x.Tiles().String())
	assert.Equal(t, 2, x.Count())
}

func TestTilesNormalize(t *testing.T) {
	assert.Equal(t, FromTile(tile.Man1), FromTile(tile.Man1).Normalize())
	assert.Equal(t, FromTile(tile.Red), FromTile(tile.Red).Normalize())
	assert.EqualValues(t, 0, FromTile(tile.TileEnd).Normalize())

	assert.Equal(t, AllTiles, AllTiles.Normalize())
	assert.Equal(t, KokushiTiles, KokushiTiles.Normalize())
}

func TestTilesInvert(t *testing.T) {
	assert.EqualValues(t, 0, AllTiles.Invert())
	assert.Equal(t, FromTile(tile.Pin4), AllTiles.Unset(tile.Pin4).Invert())
}

func TestTilesMerge(t *testing.T) {
	x1 := FromTiles(tile.Man1, tile.Red, tile.Pin1)
	x2 := FromTiles(tile.Man2, tile.Red, tile.Pin2)
	assert.Equal(t, FromTiles(tile.Man1, tile.Man2, tile.Red, tile.Pin1, tile.Pin2), x1.Merge(x2))

	assert.Equal(t, FromTiles(tile.Man1, tile.Pin1), x1.Sub(x2))
	assert.Equal(t, FromTiles(tile.Man2, tile.Pin2), x2.Sub(x1))
}

func TestTilesEach(t *testing.T) {
	x1 := FromTiles(tile.Man1, tile.Red, tile.Pin1)
	t.Run("all", func(t *testing.T) {
		var res tile.Tiles
		x1.Each(func(t tile.Tile) bool {
			res = append(res, t)
			return true
		})
		assert.Equal(t, tile.Tiles{tile.Man1, tile.Pin1, tile.Red}, res)
	})
	t.Run("sequence range", func(t *testing.T) {
		var res tile.Tiles
		x1.EachRange(tile.SequenceBegin, tile.SequenceEnd, func(t tile.Tile) bool {
			res = append(res, t)
			return true
		})
		assert.Equal(t, tile.Tiles{tile.Man1, tile.Pin1}, res)
	})
	t.Run("not sequence range", func(t *testing.T) {
		var res tile.Tiles
		x1.EachRange(tile.SequenceEnd, tile.TileEnd, func(t tile.Tile) bool {
			res = append(res, t)
			return true
		})
		assert.Equal(t, tile.Tiles{tile.Red}, res)
	})
	t.Run("stop", func(t *testing.T) {
		var res tile.Tiles
		x1.Each(func(t tile.Tile) bool {
			res = append(res, t)
			return false
		})
		assert.Equal(t, tile.Tiles{tile.Man1}, res)
	})
}

func TestTilesValues(t *testing.T) {
	assert.Equal(t, 34, AllTiles.Count())
	assert.Equal(t, 0, AllTiles.Invert().Count())
	assert.Equal(t, 13, KokushiTiles.Count())
	assert.Equal(t, 21, KokushiTiles.Invert().Count())
	assert.EqualValues(t, 1, FromTile(tile.Man1))

	assert.Equal(t, "19m19p19s1234567z", KokushiTiles.Tiles().String())
}
