package shanten

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bitbucket.org/dnovikoff/tempai-core/compact"
)

func testUkeIre(t *testing.T, in string) string {
	tg := compact.NewTileGenerator()
	tiles, err := tg.CompactFromString(in)
	require.NoError(t, err, in)
	require.Equal(t, 13, tiles.Count())
	results := CalculateShanten(tiles, 0, nil)
	uke := results.UkeIre

	return fmt.Sprintf("%v/%v/%v = %v", results.Value, uke.UniqueTiles().Count(), uke.Count(), uke.UniqueTiles().Tiles())
}

func TestUkeIreSimpleOthers(t *testing.T) {
	tst := func(in string) string {
		return testUkeIre(t, in)
	}

	// Kokushi hand
	assert.Equal(t, "2/3/12 = 567z", tst("119s19p19m1234z23m"))
	// Kokushi 13
	assert.Equal(t, "3/13/42 = 19m19p19s1234567z", tst("19s19p19m1234z234m"))

	// 7 Pairs waits
	assert.Equal(t, "1/3/9 = 129m", tst("1122556677z129m"))
}

func TestUkeIreSimple(t *testing.T) {
	tst := func(in string) string {
		return testUkeIre(t, in)
	}

	// From marujan
	assert.Equal(t, "1/10/30 = 1234m123456p", tst("233m1122334p111s"))
	assert.Equal(t, "1/11/34 = 34567m123456p", tst("335m1122334p111s"))

	assert.Equal(t, "1/8/22 = 12345678m", tst("1233446888m444p"))
	assert.Equal(t, "1/8/22 = 12345789m", tst("1233448889m444p"))

	assert.Equal(t, "1/9/29 = 34567p1234s", tst("3335p233778899s"))
	assert.Equal(t, "1/13/42 = 345678m3456789s", tst("444556m2225678s"))
	assert.Equal(t, "1/9/29 = 34567m1234p", tst("445999m1123p555z"))

	assert.Equal(t, "1/5/14 = 12378m", tst("1133889m777p789s"))
	assert.Equal(t, "1/11/34 = 12345m256789p", tst("3444m226778p567s"))

	assert.Equal(t, "1/4/10 = 34m8s6z", tst("335m56788s22266z"))

	assert.Equal(t, "1/7/24 = 3456789m", tst("4578m234567789p"))
	assert.Equal(t, "1/7/24 = 3456789m", tst("4578m345567789p"))

	assert.Equal(t, "1/7/24 = 3456789m", tst("5778m123345666p"))
	assert.Equal(t, "1/10/28 = 34567m12346p", tst("577m1123345666p"))

	assert.Equal(t, "1/5/17 = 1458m5p", tst("236788m5789p555z"))
	assert.Equal(t, "1/9/29 = 123456789m", tst("2366788m789p555z"))

	assert.Equal(t, "1/7/24 = 1234567m", tst("2356m111123777p"))
	assert.Equal(t, "1/7/24 = 1234678m", tst("2368m111123777p"))

	assert.Equal(t, "1/6/14 = 12345m3z", tst("12233355666m33z"))
	assert.Equal(t, "1/7/24 = 1234567m", tst("1445m222567p567s"))

	assert.Equal(t, "1/12/35 = 12345m1234567s", tst("112355m2344556s"))
}

func TestUkeIreCool(t *testing.T) {
	tst := func(in string) string {
		return testUkeIre(t, in)
	}
	// Good example of wrong calculation, caused by full set cuts first
	// m33 p35 66778 s567
	assert.Equal(t, "1/11/38 = 34567m123458p", tst("335m3566778p567s"))
}
