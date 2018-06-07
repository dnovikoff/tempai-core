package compact

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dnovikoff/tempai-core/tile"
)

func TestMaskCountingNaive(t *testing.T) {
	for i := Mask(0); i < 16; i++ {
		assert.Equal(t, i.NaiveCount(), i.Count())
	}
	assert.Equal(t, 4, Mask(15).Count())
}

func TestMask(t *testing.T) {
	mask := NewMask(0, tile.Pin6)
	assert.Equal(t, "6p", mask.Tile().String())

	assert.True(t, mask.IsEmpty())
	assert.False(t, mask.IsFull())
	assert.Equal(t, 0, mask.Count())

	mask.SetCopyBit(tile.CopyID(0))
	// original does not change
	assert.True(t, mask.IsEmpty())

	mask = mask.SetCopyBit(tile.CopyID(0))
	assert.False(t, mask.IsEmpty())
	assert.False(t, mask.IsFull())
	assert.Equal(t, 1, mask.Count())

	mask = mask.SetCopyBit(tile.CopyID(0))
	assert.False(t, mask.IsEmpty())
	assert.False(t, mask.IsFull())
	assert.Equal(t, 1, mask.Count())

	mask = mask.SetCopyBit(tile.CopyID(1)).SetCopyBit(2).SetCopyBit(3)
	assert.False(t, mask.IsEmpty())
	assert.True(t, mask.IsFull())
	assert.Equal(t, 4, mask.Count())

	assert.Equal(t, "6p", mask.Tile().String())

	// Still 4
	mask = mask.SetCopyBit(tile.CopyID(5))
	assert.False(t, mask.IsEmpty())
	assert.True(t, mask.IsFull())
	assert.Equal(t, 4, mask.Count())

	assert.Equal(t, "6p", mask.Tile().String())

	mask = mask.UnsetCopyBit(tile.CopyID(2))

	assert.False(t, mask.IsEmpty())
	assert.False(t, mask.IsFull())
	assert.Equal(t, 3, mask.Count())
	assert.Equal(t, "6p", mask.Tile().String())
}

func TestMaskCounters(t *testing.T) {
	mask := NewMask(0, tile.Man3)
	assert.Equal(t, 0, mask.Count())
	assert.False(t, mask.IsFull())
	assert.True(t, mask.IsEmpty())

	mask = mask.SetCount(3)
	assert.Equal(t, 3, mask.Count())
	assert.False(t, mask.IsFull())
	assert.False(t, mask.IsEmpty())

	mask = mask.SetCount(4)
	assert.Equal(t, 4, mask.Count())
	assert.True(t, mask.IsFull())
	assert.False(t, mask.IsEmpty())
}

func TestMaskRemove(t *testing.T) {
	tst := func(mask Mask) string {
		return fmt.Sprintf("%04b", mask.Mask())
	}
	mask := NewMask(0, tile.Green)
	assert.Equal(t, tile.Instances{}, mask.Instances())

	mask1 := mask.InvertTiles().UnsetCopyBit(3)
	mask2 := mask.SetCopyBit(1).SetCopyBit(3)
	assert.Equal(t, "0111", tst(mask1))
	assert.Equal(t, "1010", tst(mask2))
	assert.Equal(t, "0101", tst(mask1.Remove(mask2)))
	assert.Equal(t, "1111", tst(mask1.Merge(mask2)))
}

func TestMaskPrinting(t *testing.T) {
	assert.Equal(t, "6666p", NewMask(15, tile.Pin6).Instances().String())
}

func TestMaskCount(t *testing.T) {
	tst := func(x int) int {
		return NewMask(0, tile.Man1).SetCount(x).Count()
	}
	assert.Equal(t, 0, tst(0))
	assert.Equal(t, 1, tst(1))
	assert.Equal(t, 2, tst(2))
	assert.Equal(t, 3, tst(3))
	assert.Equal(t, 4, tst(4))

	assert.Equal(t, 0, tst(5))
	assert.Equal(t, 0, tst(400))
	assert.Equal(t, 0, tst(-1))
}

func TestMaskFirstCopy(t *testing.T) {
	m := NewMask(0, tile.Man1).SetCount(4)
	assert.EqualValues(t, 0, m.FirstCopy())
	assert.True(t, m.Check(0))
	assert.True(t, m.Check(1))
	assert.True(t, m.Check(2))
	assert.True(t, m.Check(3))

	m = m.UnsetCopyBit(0)
	assert.EqualValues(t, 1, m.FirstCopy())
	assert.False(t, m.Check(0))
	assert.True(t, m.Check(1))
	assert.True(t, m.Check(2))
	assert.True(t, m.Check(3))

	m = m.UnsetCopyBit(1)
	assert.EqualValues(t, 2, m.FirstCopy())
	assert.False(t, m.Check(0))
	assert.False(t, m.Check(1))
	assert.True(t, m.Check(2))
	assert.True(t, m.Check(3))

	m = m.UnsetCopyBit(2)
	assert.EqualValues(t, 3, m.FirstCopy())
	assert.False(t, m.Check(0))
	assert.False(t, m.Check(1))
	assert.False(t, m.Check(2))
	assert.True(t, m.Check(3))

	m = m.UnsetCopyBit(3)
	assert.EqualValues(t, tile.NullCopy, m.FirstCopy())
	assert.False(t, m.Check(0))
	assert.False(t, m.Check(1))
	assert.False(t, m.Check(2))
	assert.False(t, m.Check(3))
}

func TestMaskFirst(t *testing.T) {
	assert.Equal(t, tile.Pin4.Instance(0), NewMask(0, tile.Pin4).SetCount(4).First())
	assert.Equal(t, tile.InstanceNull, NewMask(0, tile.Pin4).SetCount(0).First())
}

func TestMaskEach(t *testing.T) {
	t.Run("each", func(t *testing.T) {
		var res tile.Instances
		NewMask(0, tile.Red).SetCount(2).Each(func(i tile.Instance) bool {
			res = append(res, i)
			return true
		})
		assert.Equal(t, tile.Instances{
			tile.Red.Instance(0),
			tile.Red.Instance(1),
		}, res)
	})
	t.Run("stop", func(t *testing.T) {
		var res tile.Instances
		NewMask(0, tile.Red).SetCount(2).Each(func(i tile.Instance) bool {
			res = append(res, i)
			return false
		})
		assert.Equal(t, tile.Instances{
			tile.Red.Instance(0),
		}, res)
	})
}
