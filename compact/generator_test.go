package compact

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tempai-core/tile"
)

func TestGeneratorTiles(t *testing.T) {
	tg := NewTestGenerator(t)
	assert.Equal(t, "111122223333444455556666777788889999m111122223333444455556666777788889999p111122223333444455556666777788889999s1111222233334444555566667777z", tg.TilesLeft().String())
	tg.InstancesFromString("1223334444s")
	assert.Equal(t, "111122223333444455556666777788889999m111122223333444455556666777788889999p11122355556666777788889999s1111222233334444555566667777z", tg.TilesLeft().String())
}

func TestGeneratorNull(t *testing.T) {
	tg := NewTileGenerator()
	assert.Equal(t, tile.Man1.Instance(0), tg.Instance(tile.Man1))
	assert.Equal(t, tile.Man1.Instance(1), tg.Instance(tile.Man1))
	assert.Equal(t, tile.Man1.Instance(2), tg.Instance(tile.Man1))
	assert.Equal(t, tile.Man1.Instance(3), tg.Instance(tile.Man1))
	assert.Equal(t, tile.InstanceNull, tg.Instance(tile.Man1))
	assert.Equal(t, tile.InstanceNull, tg.Instance(tile.Man1))
}

func TestGeneratorCompact(t *testing.T) {
	tg := NewTileGenerator()
	tiles, err := tg.CompactFromString("12334z")
	require.NoError(t, err)
	assert.Equal(t, "12334z", tiles.Instances().String())
}

func TestGeneratorError1(t *testing.T) {
	tg := NewTileGenerator()
	_, err := tg.InstancesFromString("1111z")
	require.NoError(t, err)
	_, err = tg.InstancesFromString("11111z")
	require.Error(t, err)
	_, err = tg.InstancesFromString("incorrect")
	require.Error(t, err)
}

func TestGeneratorError(t *testing.T) {
	tg := NewTileGenerator()
	_, err := tg.CompactFromString("1111z")
	require.NoError(t, err)
	_, err = tg.CompactFromString("11111z")
	require.Error(t, err)
}
