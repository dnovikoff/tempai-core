package calc

import (
	"testing"

	"github.com/dnovikoff/tempai-core/compact"

	"github.com/stretchr/testify/assert"

	"github.com/dnovikoff/tempai-core/tile"
)

func TestCountersBlocks(t *testing.T) {
	b := counterBlock(0)
	for i := uint(0); i < 9; i++ {
		for c := 0; c < 5; c++ {
			b.set(i, c)
			assert.EqualValues(t, c, b.get(i))
		}
	}
}

func TestCounters(t *testing.T) {
	b := NewCounters()
	for tl := tile.TileBegin; tl < tile.TileEnd; tl++ {
		for c := 0; c < 5; c++ {
			b.Set(tl, c)
			assert.EqualValues(t, c, b.Get(tl))
		}
	}
}

func TestCountersInvert(t *testing.T) {
	b := NewCounters()
	assert.True(t, b.Empty())
	b.Invert()
	assert.False(t, b.Empty())
	for tl := tile.TileBegin; tl < tile.TileEnd; tl++ {
		t.Run("Get"+tl.String(), func(t *testing.T) {
			assert.Equal(t, 4, b.Get(tl))
		})
	}
	assert.Equal(t, compact.AllTiles, b.Tiles())
	for tl := tile.TileBegin; tl < tile.TileEnd; tl++ {
		b.Set(tl, 0)
	}
	t.Log(b)
	assert.True(t, b.Empty())
}
