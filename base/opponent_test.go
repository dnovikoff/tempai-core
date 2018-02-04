package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpponent(t *testing.T) {
	assert.Equal(t, Right, Self.Next())
	assert.Equal(t, Front, Right.Next())
	assert.Equal(t, Left, Front.Next())
	assert.Equal(t, Self, Left.Next())

	assert.Equal(t, Left, Self.Prev())
	assert.Equal(t, Self, Right.Prev())
	assert.Equal(t, Right, Front.Prev())
	assert.Equal(t, Front, Left.Prev())
}

func TestOpponentAdvance(t *testing.T) {
	assert.Equal(t, Left, Self.Advance(-1))
	assert.Equal(t, Front, Self.Advance(-2))
	assert.Equal(t, Right, Self.Advance(-3))
	assert.Equal(t, Self, Self.Advance(-4))
	assert.Equal(t, Left, Self.Advance(-5))

	assert.Equal(t, Self, Right.Advance(-1))
	assert.Equal(t, Left, Right.Advance(-2))
	assert.Equal(t, Front, Right.Advance(-3))
	assert.Equal(t, Right, Right.Advance(-4))
	assert.Equal(t, Self, Right.Advance(-5))

	assert.Equal(t, Right, Front.Advance(-1))
	assert.Equal(t, Self, Front.Advance(-2))
	assert.Equal(t, Left, Front.Advance(-3))
	assert.Equal(t, Front, Front.Advance(-4))
	assert.Equal(t, Right, Front.Advance(-5))
}
