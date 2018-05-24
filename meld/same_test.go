package meld

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/tile"
)

// samePonPart
// sameKanUpgraded

// samePonOpened1
// samePonOpened2
// samePonOpened3

func TestSameKan(t *testing.T) {
	m := NewKan(tile.Sou9, 1)
	require.Equal(t, TypeSame, m.Meld().Type())
	assert.True(t, 0 != m)
	assert.True(t, m.IsComplete())
	assert.False(t, m.IsOpened())
	assert.False(t, m.IsBadWait())
	assert.EqualValues(t, 0, m.OriginalWaits())
	assert.EqualValues(t, 0, m.Waits())
	assert.EqualValues(t, 0, m.OpenedBy())
	i := m.Meld().Instances()
	assert.Equal(t, "9999s", i.Tiles().String())
	assert.Equal(t, tile.Instances{
		tile.Sou9.Instance(0),
		tile.Sou9.Instance(1),
		tile.Sou9.Instance(2),
		tile.Sou9.Instance(3),
	}, i)
	assert.EqualValues(t, 0, m.Open(tile.Sou9.Instance(0), base.Left))
	assert.EqualValues(t, 0, m.Upgrade())
}

func TestSameKanO(t *testing.T) {
	kans := []Same{
		NewKan(tile.Green, 0),
		NewKan(tile.Green, 1),
		NewKan(tile.Green, 2),
		NewKan(tile.Green, 3),
	}
	mp := map[Same]bool{}
	for k, v := range kans {
		mp[v] = true
		assert.False(t, v.IsOpened())
		assert.EqualValues(t, k, v.OpenedCopy())
	}
	assert.Equal(t, 4, len(mp))
}

func TestSamePon(t *testing.T) {
	m := NewPon(tile.Sou8, 2)
	require.Equal(t, TypeSame, m.Meld().Type())
	assert.True(t, 0 != m)
	assert.True(t, m.IsComplete())
	assert.False(t, m.IsOpened())
	assert.False(t, m.IsBadWait())
	assert.EqualValues(t, 0, m.OriginalWaits())
	assert.EqualValues(t, 0, m.Waits())
	assert.Equal(t, compact.NewFromTile(tile.Sou8), m.OpenedBy())
	i := m.Meld().Instances()
	assert.Equal(t, "888s", i.Tiles().String())
	assert.Equal(t, tile.Instances{
		tile.Sou8.Instance(0),
		tile.Sou8.Instance(1),
		tile.Sou8.Instance(3),
	}, i)

	assert.EqualValues(t, 0, m.Upgrade())

	assert.EqualValues(t, 0, m.Open(tile.Sou8.Instance(0), base.Right))
	m = Same(m.Open(tile.Sou8.Instance(2), base.Right))
	require.Equal(t, NewKanOpened(tile.Sou8, 2, base.Right), m)
	require.Equal(t, TypeSame, m.Meld().Type())
	require.NotEqual(t, 0, m)
	assert.True(t, m.IsComplete())
	assert.True(t, m.IsOpened())
	assert.False(t, m.IsBadWait())
	assert.EqualValues(t, 0, m.OriginalWaits())
	assert.EqualValues(t, 0, m.Waits())
	assert.EqualValues(t, 0, m.OpenedBy())
	i = m.Meld().Instances()
	assert.Equal(t, "8888s", i.Tiles().String())
	assert.Equal(t, tile.Instances{
		tile.Sou8.Instance(0),
		tile.Sou8.Instance(1),
		tile.Sou8.Instance(2),
		tile.Sou8.Instance(3),
	}, i)
}

func TestSamePart(t *testing.T) {
	m := NewPonPart(tile.West, 2, 3)
	require.Equal(t, TypeSame, m.Meld().Type())
	require.True(t, 0 != m)
	assert.False(t, m.IsComplete())
	assert.False(t, m.IsOpened())
	assert.False(t, m.IsBadWait())
	w := compact.NewFromTile(tile.West)
	assert.Equal(t, w, m.OriginalWaits())
	assert.Equal(t, w, m.Waits())
	assert.Equal(t, w, m.OpenedBy())
	i := m.Meld().Instances()
	assert.Equal(t, "33z", i.Tiles().String())
	assert.Equal(t, tile.Instances{
		tile.West.Instance(0),
		tile.West.Instance(1),
	}, i)

	assert.EqualValues(t, 0, m.Upgrade())

	assert.EqualValues(t, 0, m.Open(tile.West.Instance(0), base.Front))
	m = Same(m.Open(tile.West.Instance(2), base.Front))
	require.Equal(t, NewPonOpened(tile.West, 2, 3, base.Front), m)
	require.Equal(t, TypeSame, m.Meld().Type())
	require.True(t, 0 != m)
	assert.True(t, m.IsComplete())
	assert.True(t, m.IsOpened())
	assert.False(t, m.IsBadWait())
	assert.Equal(t, w, m.OriginalWaits())
	assert.EqualValues(t, 0, m.Waits())
	assert.EqualValues(t, 0, m.OpenedBy())
	i = m.Meld().Instances()
	assert.Equal(t, "333z", i.Tiles().String())
	assert.Equal(t, tile.Instances{
		tile.West.Instance(0),
		tile.West.Instance(1),
		tile.West.Instance(2),
	}, i)

	m = m.Upgrade()
	require.Equal(t, NewKanUpgraded(tile.West, 2, 3, base.Front), m)
	require.Equal(t, TypeSame, m.Meld().Type())
	require.True(t, 0 != m)
	assert.True(t, m.IsComplete())
	assert.True(t, m.IsOpened())
	assert.False(t, m.IsBadWait())
	assert.Equal(t, w, m.OriginalWaits())
	assert.EqualValues(t, 0, m.Waits())
	assert.EqualValues(t, 0, m.OpenedBy())
	i = m.Meld().Instances()
	assert.Equal(t, "3333z", i.Tiles().String())
	assert.Equal(t, tile.Instances{
		tile.West.Instance(0),
		tile.West.Instance(1),
		tile.West.Instance(2),
		tile.West.Instance(3),
	}, i)
}

func TestSamePartFromExisting(t *testing.T) {
	assert.Equal(t, NewPonPart(tile.Pin2, 1, 3), NewPonPartFromExisting(tile.Pin2, 0, 2))
}

func TestSameCompact(t *testing.T) {
	gen := compact.NewTileGenerator()
	tiles, err := gen.CompactFromString("3m22p111s")
	require.NoError(t, err)
	m := NewPonPart(tile.Pin2, 0, 1)
	str := func() string {
		return tiles.Instances().String()
	}
	// Unchanged - no such instances
	m.ExtractFrom(tiles)
	assert.Equal(t, "3m22p111s", str())
	m = Same(m.Rebase(tiles))

	m.ExtractFrom(tiles)
	assert.Equal(t, "3m111s", str())
	m.AddTo(tiles)
	assert.Equal(t, "3m22p111s", str())
}

func TestSameBugInstances(t *testing.T) {
	m := Meld(3462)
	require.Equal(t, TypeSame, m.Type())
	s := Same(m)
	require.Equal(t, tile.Green, s.Base())
	assert.Equal(t, samePon, s.subType())
	for k, v := range s.Meld().Instances() {
		t.Run("Green"+strconv.Itoa(k), func(t *testing.T) {
			assert.Equal(t, tile.Green, v.Tile())
		})
	}
	assert.Equal(t, tile.Instances{
		tile.Green.Instance(0),
		tile.Green.Instance(1),
		tile.Green.Instance(2),
	}, s.Meld().Instances())
}
