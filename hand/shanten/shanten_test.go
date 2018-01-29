package shanten

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tempai-core/compact"
)

func testCompact(t *testing.T, str string) compact.Instances {
	tg := compact.NewTileGenerator()
	tiles, err := tg.CompactFromString(str)
	require.NoError(t, err, str)
	return tiles
}

func testGetShantent(t *testing.T, str string) Results {
	tiles := testCompact(t, str)
	require.Equal(t, 13, tiles.Count())
	return Calculate(tiles, 0, nil)
}

func TestShantenSimple1(t *testing.T) {
	tst := func(str string) int {
		results := testGetShantent(t, str)
		return results.Total.Value
	}
	assert.Equal(t, 0, tst("44p456678s44777z"))
}

func TestShantenSimple(t *testing.T) {
	tst := func(str string) int {
		results := testGetShantent(t, str)
		return results.Total.Value
	}

	assert.Equal(t, 0, tst("11558899s11223z"))

	assert.Equal(t, 2, tst("8m1367p4566677s1z"))
	assert.Equal(t, 0, tst("123456789s1122z"))
	assert.Equal(t, 0, tst("44p456678s44777z"))

	// 13
	assert.Equal(t, 0, tst("19s19p19m1234456z"))
	assert.Equal(t, 0, tst("19s19p19m1234567z"))
	assert.Equal(t, 1, tst("19s19p18m1234567z"))
	assert.Equal(t, 2, tst("19s29p18m1234567z"))
	// This leads to 7 pairs
	assert.Equal(t, 4, tst("27s29p28m1134777z"))
}

func TestShantenBugs(t *testing.T) {
	tg := compact.NewTileGenerator()
	compact, err := tg.CompactFromString("29m3677p27s13457z")
	require.NoError(t, err)
	require.Equal(t, 13, compact.Count())

	res := Calculate(compact, 0, nil)
	assert.Equal(t, 5, res.Total.Value)
}

func TestShantenLockEasy(t *testing.T) {
	tiles := testCompact(t, "12m123456789s55z")
	require.Equal(t, 13, tiles.Count())

	used3 := testCompact(t, "333m")
	used4 := testCompact(t, "3333m")

	res := Calculate(tiles, 0, nil)
	assert.Equal(t, 0, res.Total.Value)
	res = Calculate(tiles, 0, used3)
	assert.Equal(t, 0, res.Total.Value)
	res = Calculate(tiles, 0, used4)
	assert.Equal(t, 1, res.Total.Value)
}

func TestShantenLock(t *testing.T) {
	tiles := testCompact(t, "12m123456s55z")
	require.Equal(t, 13-3, tiles.Count())

	used3 := testCompact(t, "333m789s")
	used4 := testCompact(t, "3333m789s")

	res := Calculate(tiles, 1, nil)
	assert.Equal(t, 0, res.Total.Value)
	res = Calculate(tiles, 1, used3)
	assert.Equal(t, 0, res.Total.Value)
	res = Calculate(tiles, 1, used4)
	assert.Equal(t, 1, res.Total.Value)
}

func TestShantenBug0(t *testing.T) {
	tiles := testCompact(t, "3678m3356p14s256z")
	res := Calculate(tiles, 0, nil)
	m := res.Total
	assert.Equal(t, 4, m.Value)
	uke := m.CalculateUkeIre(compact.NewTotals().Merge(tiles))
	assert.Equal(t, "12345m347p123456s256z", uke.UniqueTiles().Tiles().String())
}

func TestShantenBug1(t *testing.T) {
	tiles := testCompact(t, "369m7p1559s13567z")
	res := Calculate(tiles, 0, nil)
	m := res.Total
	assert.Equal(t, 5, m.Value)
	uke := m.CalculateUkeIre(compact.NewTotals().Merge(tiles))
	assert.Equal(t, "1369m179p19s1234567z", uke.UniqueTiles().Tiles().String())
}

func TestMonocolorBug(t *testing.T) {
	tiles := testCompact(t, "1111222235555m")
	res := Calculate(tiles, 0, nil)
	m := res.Total
	assert.Equal(t, 1, m.Value)
	uke := m.CalculateUkeIre(compact.NewTotals().Merge(tiles))
	assert.Equal(t, "3467m", uke.UniqueTiles().Tiles().String())
}

// TODO:
// used m1122333
// form is m1122 - should lead to on meld dropout
