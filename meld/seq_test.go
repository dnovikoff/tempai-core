package meld

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/tile"
)

func TestSeqBad(t *testing.T) {
	assert.EqualValues(t, 0, NewSeq(tile.Red, 1, 2, 3))
	assert.EqualValues(t, 0, NewSeq(tile.West, 1, 2, 3))
	assert.EqualValues(t, 0, NewSeq(tile.Pin8, 1, 2, 3))
	assert.True(t, 0 != NewSeq(tile.Pin7, 1, 2, 3))
}

func TestSeqComplete(t *testing.T) {
	m := NewSeq(tile.Man1, 1, 2, 3)
	require.Equal(t, TypeSeq, m.Meld().Type())
	assert.True(t, 0 != m)
	assert.True(t, m.IsComplete())
	assert.False(t, m.IsOpened())
	assert.False(t, m.IsBadWait())
	assert.EqualValues(t, 0, m.OriginalWaits())
	assert.EqualValues(t, 0, m.Waits())
	assert.EqualValues(t, 0, m.OpenedBy())
	i := m.Meld().Instances()
	assert.Equal(t, "123m", i.Tiles().String())
	assert.Equal(t, tile.Instances{
		tile.Man1.Instance(1),
		tile.Man2.Instance(2),
		tile.Man3.Instance(3),
	}, i)
}

func TestSeqPenchan1(t *testing.T) {
	m := NewSeq(tile.Man1, 1, 2, HoleCopy)
	assert.False(t, m.IsComplete())
	assert.False(t, m.IsOpened())
	assert.True(t, m.IsBadWait())
	w := compact.NewFromTile(tile.Man3)
	assert.Equal(t, w, m.OriginalWaits())
	assert.Equal(t, w, m.Waits())
	assert.Equal(t, w, m.OpenedBy())
	i := m.Meld().Instances()
	assert.Equal(t, "12m", i.Tiles().String())
	assert.Equal(t, tile.Instances{
		tile.Man1.Instance(1),
		tile.Man2.Instance(2),
	}, i)

	assert.EqualValues(t, 0, m.Open(tile.Man3.Instance(0), base.Front))

	m = Seq(m.Open(tile.Man3.Instance(0), base.Left))
	require.Equal(t, TypeSeq, m.Meld().Type())
	assert.True(t, 0 != m)
	assert.True(t, m.IsComplete())
	assert.True(t, m.IsOpened())
	assert.False(t, m.IsBadWait())
	assert.Equal(t, w, m.OriginalWaits())
	assert.EqualValues(t, 0, m.Waits())
	assert.EqualValues(t, 0, m.OpenedBy())
	i = m.Meld().Instances()
	assert.Equal(t, "123m", i.Tiles().String())
	assert.Equal(t, tile.Instances{
		tile.Man1.Instance(1),
		tile.Man2.Instance(2),
		tile.Man3.Instance(0),
	}, i)
}

func TestSeqPenchan2(t *testing.T) {
	m := NewSeq(tile.Man7, HoleCopy, 2, 3)
	assert.False(t, m.IsComplete())
	assert.False(t, m.IsOpened())
	assert.True(t, m.IsBadWait())
	w := compact.NewFromTile(tile.Man7)
	assert.Equal(t, w, m.OriginalWaits())
	assert.Equal(t, w, m.Waits())
	i := m.Meld().Instances()
	assert.Equal(t, "89m", i.Tiles().String())
	assert.Equal(t, tile.Instances{
		tile.Man8.Instance(2),
		tile.Man9.Instance(3),
	}, i)

	m = Seq(m.Open(tile.Man7.Instance(0), base.Left))
	assert.True(t, 0 != m)
	assert.True(t, m.IsComplete())
	assert.True(t, m.IsOpened())
	assert.False(t, m.IsBadWait())
	assert.Equal(t, w, m.OriginalWaits())
	assert.EqualValues(t, 0, m.Waits())
	assert.EqualValues(t, 0, m.OpenedBy())
	i = m.Meld().Instances()
	assert.Equal(t, "789m", i.Tiles().String())
	assert.Equal(t, tile.Instances{
		tile.Man7.Instance(0),
		tile.Man8.Instance(2),
		tile.Man9.Instance(3),
	}, i)
}

func TestSeqPenchanSpecialCase(t *testing.T) {
	m1 := NewSeq(tile.Man8, 2, 3, HoleCopy)
	m2 := NewSeq(tile.Man7, HoleCopy, 2, 3)
	assert.Equal(t, m1, m2)
}

func TestSeqKanchan(t *testing.T) {
	m := NewSeq(tile.Pin5, 0, HoleCopy, 3)
	assert.False(t, m.IsComplete())
	assert.False(t, m.IsOpened())
	assert.True(t, m.IsBadWait())
	w := compact.NewFromTile(tile.Pin6)
	assert.Equal(t, w, m.OriginalWaits())
	assert.Equal(t, w, m.Waits())
	i := m.Meld().Instances()
	assert.Equal(t, "57p", i.Tiles().String())
	assert.Equal(t, tile.Instances{
		tile.Pin5.Instance(0),
		tile.Pin7.Instance(3),
	}, i)

	m = Seq(m.Open(tile.Pin6.Instance(0), base.Left))
	assert.True(t, 0 != m)
	assert.True(t, m.IsComplete())
	assert.True(t, m.IsOpened())
	assert.False(t, m.IsBadWait())
	assert.Equal(t, w, m.OriginalWaits())
	assert.EqualValues(t, 0, m.Waits())
	assert.EqualValues(t, 0, m.OpenedBy())
	i = m.Meld().Instances()
	assert.Equal(t, "567p", i.Tiles().String())
	assert.Equal(t, tile.Instances{
		tile.Pin5.Instance(0),
		tile.Pin6.Instance(0),
		tile.Pin7.Instance(3),
	}, i)
}

func TestSeqRynman(t *testing.T) {
	m := NewSeq(tile.Pin5, 1, 2, HoleCopy)
	assert.False(t, m.IsComplete())
	assert.False(t, m.IsOpened())
	assert.False(t, m.IsBadWait())
	w := compact.NewFromTiles(tile.Pin4, tile.Pin7)
	assert.Equal(t, w, m.OriginalWaits())
	assert.Equal(t, w, m.Waits())
	i := m.Meld().Instances()
	assert.Equal(t, "56p", i.Tiles().String())
	assert.Equal(t, tile.Instances{
		tile.Pin5.Instance(1),
		tile.Pin6.Instance(2),
	}, i)

	original := m

	m = Seq(original.Open(tile.Pin4.Instance(0), base.Left))
	assert.True(t, 0 != m)
	assert.True(t, m.IsComplete())
	assert.True(t, m.IsOpened())
	assert.False(t, m.IsBadWait())
	assert.Equal(t, "47p", m.OriginalWaits().Tiles().String())
	assert.EqualValues(t, 0, m.Waits())
	assert.EqualValues(t, 0, m.OpenedBy())
	i = m.Meld().Instances()
	assert.Equal(t, "456p", i.Tiles().String())
	assert.Equal(t, tile.Instances{
		tile.Pin4.Instance(0),
		tile.Pin5.Instance(1),
		tile.Pin6.Instance(2),
	}, i)

	m = Seq(original.Open(tile.Pin7.Instance(0), base.Left))
	assert.True(t, 0 != m)
	assert.True(t, m.IsComplete())
	assert.True(t, m.IsOpened())
	assert.False(t, m.IsBadWait())
	assert.Equal(t, "47p", m.OriginalWaits().Tiles().String())
	assert.EqualValues(t, 0, m.Waits())
	assert.EqualValues(t, 0, m.OpenedBy())
	i = m.Meld().Instances()
	assert.Equal(t, "567p", i.Tiles().String())
	assert.Equal(t, tile.Instances{
		tile.Pin5.Instance(1),
		tile.Pin6.Instance(2),
		tile.Pin7.Instance(0),
	}, i)
}

func TestSeqCompact(t *testing.T) {
	gen := compact.NewTileGenerator()
	tiles, err := gen.CompactFromString("123456778p")
	require.NoError(t, err)
	m := NewSeq(tile.Pin5, 1, 2, 3)
	str := func() string {
		return tiles.Instances().String()
	}
	// Unchanged - no such instances
	m.ExtractFrom(tiles)
	assert.Equal(t, "123456778p", str())
	m = Seq(m.Rebase(tiles))
	assert.NotEmpty(t, m)
	m.ExtractFrom(tiles)
	assert.Equal(t, "123478p", str())
	m.AddTo(tiles)
	assert.Equal(t, "123456778p", str())
}

func TestSeqOpen(t *testing.T) {
	m := NewSeq(tile.Pin4, OpenCopy(0), 0, 0)
	c := compact.NewInstances()
	m.AddTo(c)
	assert.True(t, m.IsOpened())
	assert.True(t, c.Check(tile.Pin4.Instance(0)))
	assert.True(t, c.Check(tile.Pin5.Instance(0)))
	assert.True(t, c.Check(tile.Pin6.Instance(0)))
}

func TestSeqMultiRebase(t *testing.T) {
	gen := compact.NewTileGenerator()
	tiles, err := gen.CompactFromString("111122223333s")
	require.NoError(t, err)
	m := NewSeq(tile.Sou1, 1, 2, 3)
	str := func() string {
		return tiles.Instances().String()
	}

	m.ExtractFrom(tiles)
	assert.Equal(t, "111222333s", str())
	m.Rebase(tiles).ExtractFrom(tiles)
	assert.Equal(t, "112233s", str())
	m.Rebase(tiles).ExtractFrom(tiles)
	assert.Equal(t, "123s", str())
	m.Rebase(tiles).ExtractFrom(tiles)
	assert.Equal(t, "", str())
}

func TestSeqRebaseNotEnoughtTiles(t *testing.T) {
	gen := compact.NewTileGenerator()
	tiles, err := gen.CompactFromString("111122223333s")
	require.NoError(t, err)
	assert.EqualValues(t, 0, NewSeq(tile.Pin1, 1, 2, 3).Rebase(tiles))
}

func TestNewSeqBug(t *testing.T) {
	i := compact.NewInstances()
	s := NewSeq(tile.Sou1, HoleCopy, 1, 1)
	s.AddTo(i)
	assert.Equal(t, "23s", i.Instances().String())
}
