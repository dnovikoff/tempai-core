package calc

import (
	"testing"

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
	b.Invert()
	for tl := tile.TileBegin; tl < tile.TileEnd; tl++ {
		assert.Equal(t, 4, b.Get(tl))
	}
}
