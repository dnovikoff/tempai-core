package compact

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tempai-core/tile"
)

func TestPackedMasks(t *testing.T) {
	packed := PackedMasks(0)
	assert.Equal(t, 0, packed.CountBits())
	packed.Set(NewMask(15, tile.Pin6), 1)
	// Unchanged
	assert.Equal(t, 0, packed.CountBits())
	packed = packed.Set(NewMask(15, tile.Pin6), 1)
	assert.Equal(t, 4, packed.CountBits())

	packed = packed.Set(NewMask(15, tile.Pin6), 0)
	assert.Equal(t, 8, packed.CountBits())

	packed = packed.Set(NewMask(2, tile.Pin6), 0)
	assert.Equal(t, 5, packed.CountBits())
}

func TestInstanceCountingOne(t *testing.T) {
	compact := NewInstances()
	assert.Equal(t, 0, compact.Count())
	mask := NewMask(15, tile.Man1)
	assert.Equal(t, 4, mask.Count())
	compact.SetMask(mask)
	assert.Equal(t, 4, compact.Count())

	compact.SetMask(NewMask(15, tile.Pin6))
	assert.Equal(t, 8, compact.Count())
}

func TestInstancePrinting(t *testing.T) {
	compact := NewInstances()
	compact.SetMask(NewMask(15, tile.Pin6))
	assert.Equal(t, "6666p", compact.Instances().String())
}

func TestInstanceCounting(t *testing.T) {
	compact := NewInstances()
	assert.Equal(t, 0, compact.Count())
	compact.SetMask(NewMask(15, tile.Pin6))
	assert.Equal(t, 4, compact.Count())
	compact.SetMask(NewMask(2, tile.Pin6))
	assert.Equal(t, 1, compact.Count())

	compact.SetMask(NewMask(3, tile.Red))
	assert.Equal(t, 3, compact.Count())

	assert.Equal(t, "6p77z", compact.Instances().String())
	assert.Equal(t, "6p7z", compact.UniqueTiles().Tiles().String())
}

func TestInstanceMerge(t *testing.T) {
	first := NewInstances().AddCount(tile.Man1, 2)
	second := NewInstances().AddCount(tile.Sou8, 3)

	third := first.Merge(second)

	assert.Equal(t, "11m888s", third.Instances().String())
	assert.Equal(t, 5, third.Count())
}

func TestInstanceMaskError1(t *testing.T) {
	st := NewInstances()
	require.Equal(t, "", st.Instances().String())
	st.Set(tile.Sou4.Instance(0))
	assert.Equal(t, 1, st.Count())
	assert.Equal(t, "4s", st.Instances().String())
}

func TestInstanceMaskErrors(t *testing.T) {
	tg := NewTileGenerator()
	str := "22223333444s55z"
	inst, err := tg.InstancesFromString(str)
	require.NoError(t, err)
	require.Equal(t, 13, len(inst))
	st := NewInstances()
	st.Add(inst)

	assert.Equal(t, len(inst), st.Count())
	assert.Equal(t, str, st.Instances().String())

	assert.Equal(t, len(inst), st.Count())
}

func TestInstanceCounters(t *testing.T) {
	st := NewInstances()
	assert.Equal(t, 0, st.GetCount(tile.Man3))
	assert.False(t, st.GetMask(tile.Man3).IsFull())
	assert.True(t, st.GetMask(tile.Man3).IsEmpty())

	st.SetCount(tile.Man3, 3)
	assert.Equal(t, 3, st.GetCount(tile.Man3))
	assert.False(t, st.GetMask(tile.Man3).IsFull())
	assert.False(t, st.GetMask(tile.Man3).IsEmpty())

	st.SetCount(tile.Man3, 4)
	assert.Equal(t, 4, st.GetCount(tile.Man3))
	assert.True(t, st.GetMask(tile.Man3).IsFull())
	assert.False(t, st.GetMask(tile.Man3).IsEmpty())
}

func TestInstanceClone(t *testing.T) {
	tg := NewTestGenerator(t)
	x1 := tg.CompactFromString("126m")
	x2 := x1.Clone()
	assert.Equal(t, x1.Instances(), x2.Instances())
	x2.Set(tile.Red.Instance(1))
	assert.NotEqual(t, x1.Instances(), x2.Instances())
}

func TestInstanceCheck(t *testing.T) {
	tg := NewTestGenerator(t)
	x1 := tg.CompactFromString("1126m")
	assert.True(t, x1.Check(tile.Man1.Instance(0)))
	assert.True(t, x1.Check(tile.Man1.Instance(1)))
	assert.True(t, x1.Check(tile.Man2.Instance(0)))
	assert.True(t, x1.Check(tile.Man6.Instance(0)))

	assert.False(t, x1.Check(tile.Man1.Instance(2)))
	assert.False(t, x1.Check(tile.Man1.Instance(3)))
	assert.False(t, x1.Check(tile.Red.Instance(0)))
}

func TestInstanceEach(t *testing.T) {
	tg := NewTestGenerator(t)
	x1 := tg.CompactFromString("1m7z")
	t.Run("each", func(t *testing.T) {
		var res tile.Instances
		x1.Each(func(m Mask) bool {
			res = append(res, m.Instances()...)
			return true
		})
		assert.Equal(t, tile.Instances{
			tile.Man1.Instance(0),
			tile.Red.Instance(0),
		}, res)
	})
	t.Run("stop", func(t *testing.T) {
		var res tile.Instances
		x1.Each(func(m Mask) bool {
			res = append(res, m.Instances()...)
			return false
		})
		assert.Equal(t, tile.Instances{
			tile.Man1.Instance(0),
		}, res)
	})
}
