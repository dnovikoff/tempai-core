package shanten

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/tile"
)

func testEffectiveBest(t *testing.T, in string) string {
	tg := compact.NewTileGenerator()
	tiles, err := tg.CompactFromString(in)
	require.NoError(t, err, in)
	require.Equal(t, 14, tiles.Count())
	results := CalculateEffectivity(tiles, 0, nil)
	require.NotNil(t, results)
	result := results.Best()
	shanten := result.Shanten
	uke := shanten.UkeIre
	return fmt.Sprintf("%v -> %v/%v/%v = %v", result.Tile, shanten.Value, uke.UniqueTiles().Count(), uke.Count(), uke.UniqueTiles().Tiles())
}

func TestEffectiveSimple(t *testing.T) {
	tst := func(in string) string {
		return testEffectiveBest(t, in)
	}

	assert.Equal(t, "2m -> 1/11/34 = 34567m123456p", tst("2335m1122334p111s"))
}

func TestEffectiveTileCompare(t *testing.T) {
	assert.False(t, tileLess(tile.Pin5, tile.Pin4))
	assert.True(t, tileLess(tile.Pin1, tile.Pin4))
	assert.False(t, tileLess(tile.Pin3, tile.Man2))
	assert.True(t, tileLess(tile.Man1, tile.Pin1))
}

func TestEffectiveBug1(t *testing.T) {
	test := func(in string, opened int) string {
		tg := compact.NewTileGenerator()
		tiles, err := tg.CompactFromString(in)
		require.NoError(t, err, in)
		require.Equal(t, 0, (tiles.Count()-2)%3)
		results := CalculateEffectivity(tiles, opened, nil)
		require.NotNil(t, results)
		result := results.Best()
		shanten := result.Shanten
		uke := shanten.UkeIre
		return fmt.Sprintf("%v -> %v/%v/%v = %v", result.Tile, shanten.Value, uke.UniqueTiles().Count(), uke.Count(), uke.UniqueTiles().Tiles())
	}
	assert.Equal(t, "6p -> 0/2/8 = 69s", test("88m668p23478s7p111z", 0))
	assert.Equal(t, "6p -> 0/2/8 = 69s", test("88m668p23478s7p", 1))

	assert.Equal(t, "8s -> 1/5/16 = 8m67p69s", test("88m668p234788s222z", 0))
	assert.Equal(t, "8s -> 1/5/16 = 8m67p69s", test("88m668p234788s", 1))
}

// TODO: test for special and Best
