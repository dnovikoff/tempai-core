package tempai

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/meld"
	"github.com/dnovikoff/tempai-core/tile"
)

func tstMeldsString(in meld.Melds) string {
	ret := ""
	waits := compact.Tiles(0)
	for _, v := range in {
		waits |= v.Waits()
		tiles := v.Instances()
		if len(tiles) == 0 {
			continue
		}
		ret += tiles.String() + " "
	}
	ret += "(" + waits.Tiles().String() + ")"
	return ret
}

func tstTempai(t *testing.T, str string) string {
	obj := testTempai(t, str)
	if obj == nil {
		return "Invalid"
	}
	strs := make([]string, 0, len(obj))
	for _, v := range obj {
		strs = append(strs, tstMeldsString(v))
	}
	sort.Strings(strs)

	return strings.Join(strs, ", ")
}

func testTempai(t *testing.T, str string) TempaiMelds {
	tg := compact.NewTileGenerator()
	inst, err := tg.CompactFromString(str)
	require.NoError(t, err, str)
	return Calculate(inst, nil)
}

func testAway(t *testing.T, str string) compact.Tiles {
	tg := compact.NewTileGenerator()
	inst, err := tg.CompactFromString(str)
	require.NoError(t, err)
	return GetTempaiTiles(inst, nil)
}

func testAwayString(t *testing.T, str string) string {
	return testAway(t, str).Tiles().String()
}

func TestTempai7(t *testing.T) {
	tst := func(str string) string {
		return tstTempai(t, str)
	}
	assert.Equal(t, "11m 11p 11z 55z 66z 77z 2z (2z)", tst("11m11p112556677z"))
	// 14 tiles
	assert.Equal(t, "Invalid", tst("11m11p1122556677z"))
	assert.Equal(t, "1m1p12567z", testAwayString(t, "11m11p1122556677z"))
}

func TestTempai13(t *testing.T) {
	tst := func(str string) string {
		return tstTempai(t, str)
	}
	assert.Equal(t, "Invalid", tst("12m19s19p12234567z"))
	assert.Equal(t, "2m", testAwayString(t, "12m19s19p12234567z"))
	assert.Equal(t, "1m 9m 1p 9p 1s 9s 1z 3z 4z 5z 7z 22z (6z)", tst("19m19s19p1223457z"))
}

func TestTempaiRegular(t *testing.T) {
	tst := func(str string) string {
		return tstTempai(t, str)
	}
	assert.Equal(t, "55p 333m 888s 555z 66z (6z), 66z 333m 888s 555z 55p (5p)", tst("555z333m888s66z55p"))

	// This is not possible
	// 111m 44m(+4m) 456m 456m 55m
	// 111m 44m 456m 456m 55m(+5m)
	assert.Equal(t, "456m 111m 444m 555m 6m (6m), 55m 456m 111m 444m 56m (47m), 66m 111m 444m 555m 45m (36m)", tst("1114444555566m"))

	// This is not possible
	// m111 m4(+m4) m444 m777 m888
	assert.Equal(t, "Invalid", tst("1114444777888m"))

	// Hand is ready
	assert.Equal(t, "Invalid", tst("555z333m888s66z55p1z"))
}

func TestTempai9Gates(t *testing.T) {
	temp := testTempai(t, "1112345678999m")
	assert.Equal(t, 11, len(temp))
	assert.Equal(t, "123456789m", temp.Waits().Tiles().String())
}

func TestTempai9Gates14(t *testing.T) {
	assert.Nil(t, testTempai(t, "11123455678999m"))
}

func TestTempaiReady(t *testing.T) {
	tst := func(str string) string {
		return tstTempai(t, str)
	}

	tiles := func(str string) string {
		return testAwayString(t, str)
	}

	assert.Equal(t, "Invalid", tst("55p555z888s333m666z"))
	assert.Equal(t, "Invalid", tst("19m19p19s1234567z55s"))

	assert.Equal(t, "3m5p8s56z", tiles("55p555z888s333m666z"))
	assert.Equal(t, "5p", tiles("19m19p19s1234567z5p"))
}

func TestTempaiWaits(t *testing.T) {
	tst := func(str string) string {
		melds := testTempai(t, str)
		if !assert.NotNil(t, melds) {
			require.FailNow(t, "Bad", str)
		}
		return melds.Waits().Tiles().String()
	}

	require.Equal(t, "346m", tst("4555m123789s333z"))

	// Nobetan
	assert.Equal(t, "14m", tst("1234m123789s333z"))

	// Sanmenten
	assert.Equal(t, "147m", tst("23456m11p123456s"))
	// Sanmentan
	assert.Equal(t, "147m", tst("1234567m123456s"))

	// 8 sides
	assert.Equal(t, "23456789m", tst("1113334567888m"))
	assert.Equal(t, "1347s", tst("1222234456678s"))
	assert.Equal(t, "1247p", tst("1123333445566p"))
	assert.Equal(t, "1345m", tst("2333445556777m"))

	// 3 side-pons
	assert.Equal(t, "346m", tst("4555m123456789s"))
	// 2 side-pons
	assert.Equal(t, "34m", tst("3555m123456789s"))
	// 5-side
	assert.Equal(t, "35689m", tst("4445678m456789s"))
	// 4-side
	assert.Equal(t, "69p47s", tst("44456s66678p111z"))

	assert.Equal(t, "6789s", tst("7788999s123456p"))

	//ipeko
	assert.Equal(t, "25s2z", tst("22334455s11122z"))
	assert.Equal(t, "2356s", tst("2233445566s123p"))

	// kantankan 3335777
	assert.Equal(t, "456s", tst("3335777s111222z"))
	// tatsumaki
	assert.Equal(t, "23456s", tst("3334555s111222z"))
	// hapubidjin 2223456777
	assert.Equal(t, "12345678s", tst("2223456777s111z"))
	// kokushi 13
	assert.Equal(t, "12345678s", tst("2223456777s111z"))

	// Difiicult
	// 7(6)
	assert.Equal(t, "124569s", tst("1233334555678s"))
	assert.Equal(t, "14s5z", tst("22223333444s55z"))

	// Interpret
	assert.Equal(t, "34p", tst("222p123s111555z4p"))

	// 6 is used
	assert.Equal(t, "12345789m", tst("1112345666678m"))
	// 4 is used
	assert.Equal(t, "12356789m", tst("2344445678999m"))
}

func TestTempaiCorrectNumbers(t *testing.T) {
	// 555678m56788p678s
	ints := []int{16, 17, 19, 24, 30, 53, 58, 60, 65, 66, 94, 99, 101}
	require.Equal(t, 13, len(ints))

	validateInstances := func(winInstances tile.Instances) bool {
		winInts := make([]int, len(winInstances))
		for k, v := range winInstances {
			winInts[k] = int(v)
		}
		return assert.Equal(t, ints, winInts)
	}

	validate := func(melds meld.Melds) bool {
		i := getMeldsInstances(melds).Instances()
		return validateInstances(i)
	}

	instances := make(tile.Instances, len(ints))
	for k, v := range ints {
		instances[k] = tile.Instance(v)
	}
	compact := compact.NewInstances().Add(instances)
	require.True(t, validateInstances(compact.Instances()))
	temp := Calculate(compact, nil)
	require.NotNil(t, temp)

	for k, variant := range temp {
		require.True(t, validate(variant), fmt.Sprintf("[%v] %v", k, variant))
	}
}

func TestTempaiT1(t *testing.T) {
	temp := testTempai(t, "234m13999p567s11z")
	assert.Equal(t, 1, len(temp))
	assert.Equal(t, "2p", temp.Waits().Tiles().String())
}

func TestNilIndex(t *testing.T) {
	var x TempaiMelds
	assert.Nil(t, x.Index())
}
