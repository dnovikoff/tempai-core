package base

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dnovikoff/tempai-core/tile"
)

func TestWindToOpponent(t *testing.T) {
	assert.Equal(t, Self, WindEast.Opponent(WindEast))
	assert.Equal(t, Right, WindEast.Opponent(WindSouth))
	assert.Equal(t, Front, WindEast.Opponent(WindWest))
	assert.Equal(t, Left, WindEast.Opponent(WindNorth))

	assert.Equal(t, Left, WindSouth.Opponent(WindEast))
	assert.Equal(t, Self, WindSouth.Opponent(WindSouth))
	assert.Equal(t, Right, WindSouth.Opponent(WindWest))
	assert.Equal(t, Front, WindSouth.Opponent(WindNorth))

	assert.Equal(t, Front, WindWest.Opponent(WindEast))
	assert.Equal(t, Left, WindWest.Opponent(WindSouth))
	assert.Equal(t, Self, WindWest.Opponent(WindWest))
	assert.Equal(t, Right, WindWest.Opponent(WindNorth))

	assert.Equal(t, Right, WindNorth.Opponent(WindEast))
	assert.Equal(t, Front, WindNorth.Opponent(WindSouth))
	assert.Equal(t, Left, WindNorth.Opponent(WindWest))
	assert.Equal(t, Self, WindNorth.Opponent(WindNorth))
}

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
