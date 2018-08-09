package effective

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/calc"
	"github.com/dnovikoff/tempai-core/tile"
)

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
	for _, v := range []struct {
		Hand     string
		Opened   int
		Expected string
	}{
		{"88m668p23478s7p111z", 0, "6p -> 0/2/8 = 69s"},
		{"88m668p23478s7p", 1, "6p -> 0/2/8 = 69s"},
		{"88m668p234788s222z", 0, "8s -> 1/5/16 = 8m67p69s"},
		{"88m668p234788s", 1, "8s -> 1/5/16 = 8m67p69s"},
	} {
		t.Run(v.Hand, func(t *testing.T) {
			tg := compact.NewTileGenerator()
			tiles, err := tg.CompactFromString(v.Hand)
			require.NoError(t, err)
			require.Equal(t, 0, (tiles.CountBits()-2)%3)
			results := Calculate(tiles, calc.Opened(v.Opened))
			require.NotNil(t, results)
			sorted := results.Sorted(tiles)
			used := compact.NewTotals().Merge(tiles)
			result := sorted.Best()
			s := result.Shanten
			uke := s.Total.CalculateUkeIre(used)
			assert.Equal(t, v.Expected,
				fmt.Sprintf("%v -> %v/%v/%v = %v",
					result.Tile,
					s.Total.Value,
					uke.UniqueTiles().Count(),
					uke.Count(),
					uke.UniqueTiles().Tiles()))
		})
	}
}

func TestEffectivePairsBug2(t *testing.T) {
	tg := compact.NewTileGenerator()
	tiles, err := tg.CompactFromString("333667799p22444z")
	require.NoError(t, err)
	results := Calculate(tiles)
	require.NotNil(t, results)
	assert.Equal(t, 2, results[tile.North].Regular.Value)
	assert.Equal(t, 1, results[tile.North].Pairs.Value)
	assert.Equal(t, 1, results[tile.North].Total.Value)
	assert.Equal(t, 1, results[tile.Pin3].Total.Value)
	assert.Equal(t, 1, results[tile.South].Total.Value)
}

func TestEffectivePairsBug3(t *testing.T) {
	tg := compact.NewTileGenerator()
	tiles, err := tg.CompactFromString("11112222333344z")
	require.NoError(t, err)
	results := Calculate(tiles)
	require.NotNil(t, results)
	assert.Equal(t, 5, results[tile.East].Pairs.Value)
	assert.Equal(t, 5, results[tile.South].Pairs.Value)
	assert.Equal(t, 6, results[tile.North].Pairs.Value)
}

func testEffectiveBest(t *testing.T, in string) string {
	tg := compact.NewTileGenerator()
	tiles, err := tg.CompactFromString(in)
	require.NoError(t, err, in)
	require.Equal(t, 14, tiles.CountBits())
	results := Calculate(tiles)
	require.NotNil(t, results)
	used := compact.NewTotals().Merge(tiles)
	result := results.Sorted(tiles).Best()
	s := result.Shanten
	uke := s.Total.CalculateUkeIre(used)
	return fmt.Sprintf("%v -> %v/%v/%v = %v",
		result.Tile,
		s.Total.Value,
		uke.UniqueTiles().Count(),
		uke.Count(),
		uke.UniqueTiles().Tiles())
}
