package base

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dnovikoff/tempai-core/tile"
)

func TestWindAdvance(t *testing.T) {
	assert.Equal(t, WindEast, WindEast.Advance(0))
	assert.Equal(t, WindSouth, WindEast.Advance(1))
	assert.Equal(t, WindWest, WindEast.Advance(2))
	assert.Equal(t, WindNorth, WindEast.Advance(3))
	assert.Equal(t, WindEast, WindEast.Advance(4))
	assert.Equal(t, WindSouth, WindEast.Advance(5))

	assert.Equal(t, WindEast, WindEast.Advance(0))
	assert.Equal(t, WindNorth, WindEast.Advance(-1))
	assert.Equal(t, WindWest, WindEast.Advance(-2))
	assert.Equal(t, WindSouth, WindEast.Advance(-3))
	assert.Equal(t, WindEast, WindEast.Advance(-4))
	assert.Equal(t, WindNorth, WindEast.Advance(-5))

	assert.Equal(t, WindSouth, WindSouth.Advance(0))
	assert.Equal(t, WindWest, WindSouth.Advance(1))
	assert.Equal(t, WindNorth, WindSouth.Advance(2))
	assert.Equal(t, WindEast, WindSouth.Advance(3))
	assert.Equal(t, WindSouth, WindSouth.Advance(4))
	assert.Equal(t, WindWest, WindSouth.Advance(5))

	assert.Equal(t, WindSouth, WindSouth.Advance(0))
	assert.Equal(t, WindEast, WindSouth.Advance(-1))
	assert.Equal(t, WindNorth, WindSouth.Advance(-2))
	assert.Equal(t, WindWest, WindSouth.Advance(-3))
	assert.Equal(t, WindSouth, WindSouth.Advance(-4))
	assert.Equal(t, WindEast, WindSouth.Advance(-5))
}

func TestWindCheck(t *testing.T) {
	assert.True(t, WindEast.CheckTile(tile.East))
	assert.True(t, WindSouth.CheckTile(tile.South))
	assert.True(t, WindWest.CheckTile(tile.West))
	assert.True(t, WindNorth.CheckTile(tile.North))

	assert.False(t, WindEast.CheckTile(tile.West))
}
