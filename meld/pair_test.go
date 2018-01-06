package meld

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bitbucket.org/dnovikoff/tempai-core/compact"
	"bitbucket.org/dnovikoff/tempai-core/tile"
)

func TestPair(t *testing.T) {
	assert.EqualValues(t, 0, NewPair(tile.Green, 1, 1))

	meld := NewPair(tile.Green, 0, 1)
	require.Equal(t, TypePair, meld.Meld().Type())
	assert.Equal(t, tile.Green, meld.Base())
	assert.EqualValues(t, 0, meld.Waits())
	assert.True(t, meld.IsComplete())
	assert.EqualValues(t, 0, meld.c1())
	assert.EqualValues(t, 1, meld.c2())
	assert.Equal(t, tile.Instances{128, 129}, meld.Meld().Instances())
}

func TestPairTanki(t *testing.T) {
	meld := NewTanki(tile.Pin1.Instance(2))
	require.Equal(t, TypePair, meld.Meld().Type())
	assert.Equal(t, tile.Pin1, meld.Base())
	assert.Equal(t, compact.NewFromTile(tile.Pin1), meld.Waits())
	assert.False(t, meld.IsComplete())
	assert.EqualValues(t, 2, meld.c1())
	assert.EqualValues(t, 0, meld.c2())
	assert.Equal(t, tile.Instances{38}, meld.Meld().Instances())
}

func TestPairOne(t *testing.T) {
	meld := NewOne(tile.Pin1.Instance(2))
	require.Equal(t, TypePair, meld.Meld().Type())
	assert.Equal(t, tile.Pin1, meld.Base())
	assert.EqualValues(t, 0, meld.Waits())
	assert.True(t, meld.IsComplete())
	assert.EqualValues(t, 2, meld.c1())
	assert.EqualValues(t, 0, meld.c2())
	assert.Equal(t, tile.Instances{38}, meld.Meld().Instances())
}

func TestPairHole(t *testing.T) {
	meld := NewHole(tile.East)
	require.Equal(t, TypePair, meld.Meld().Type())
	assert.Equal(t, tile.East, meld.Base())
	assert.Equal(t, compact.NewFromTile(tile.East), meld.Waits())
	assert.False(t, meld.IsComplete())
	assert.EqualValues(t, 0, meld.c1())
	assert.EqualValues(t, 0, meld.c2())
	assert.Equal(t, tile.Instances{}, meld.Meld().Instances())
}

func TestPairRebase(t *testing.T) {
	tg := compact.NewTileGenerator()
	i, err := tg.CompactFromString("44p")
	require.NoError(t, err)
	pair := Pair(NewPair(tile.Pin4, 0, 1).Rebase(i))
	require.NotEqual(t, 0, pair)
	require.EqualValues(t, 0, pair.c1())
	require.EqualValues(t, 1, pair.c2())
}
