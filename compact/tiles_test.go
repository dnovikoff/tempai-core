package compact

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"bitbucket.org/dnovikoff/tempai-core/tile"
)

func TestTilesNew(t *testing.T) {
	assert.Equal(t, "47p", NewFromTiles(tile.Pin4, tile.Pin7).Tiles().String())
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
