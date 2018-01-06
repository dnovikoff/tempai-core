package compact

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"bitbucket.org/dnovikoff/tempai-core/tile"
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

	mask.SetCopyBit(tile.CopyId(0))
	// original does not change
	assert.True(t, mask.IsEmpty())

	mask = mask.SetCopyBit(tile.CopyId(0))
	assert.False(t, mask.IsEmpty())
	assert.False(t, mask.IsFull())
	assert.Equal(t, 1, mask.Count())

	mask = mask.SetCopyBit(tile.CopyId(0))
	assert.False(t, mask.IsEmpty())
	assert.False(t, mask.IsFull())
	assert.Equal(t, 1, mask.Count())

	mask = mask.SetCopyBit(tile.CopyId(1)).SetCopyBit(2).SetCopyBit(3)
	assert.False(t, mask.IsEmpty())
	assert.True(t, mask.IsFull())
	assert.Equal(t, 4, mask.Count())

	assert.Equal(t, "6p", mask.Tile().String())

	// Still 4
	mask = mask.SetCopyBit(tile.CopyId(5))
	assert.False(t, mask.IsEmpty())
	assert.True(t, mask.IsFull())
	assert.Equal(t, 4, mask.Count())

	assert.Equal(t, "6p", mask.Tile().String())

	mask = mask.UnsetCopyBit(tile.CopyId(2))

	assert.False(t, mask.IsEmpty())
	assert.False(t, mask.IsFull())
	assert.Equal(t, 3, mask.Count())
	assert.Equal(t, "6p", mask.Tile().String())
}

func TestMaskCounters(t *testing.T) {
	mask := NewMaskByCount(0, tile.Man3)
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
