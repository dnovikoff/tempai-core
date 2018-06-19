package tempai

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tempai-core/compact"
)

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
	assert.Equal(t, "19m19p19s1223457z (6z)", tst("19m19s19p1223457z"))
	assert.Equal(t, "1m19p19s12345677z (9m)", tst("19s19p1m12345677z"))
}

func TestTempaiRegular(t *testing.T) {
	for _, v := range []struct {
		hand     string
		expected string
	}{
		{"555z333m888s66z55p", "55p 333m 888s 555z 66z (6z), 66z 333m 888s 555z 55p (5p)"},
		// This is not possible
		// 111m 44m(+4m) 456m 456m 55m
		// 111m 44m 456m 456m 55m(+5m)
		{"1114444555566m", "456m 111m 444m 555m 6m (6m), 55m 456m 111m 444m 56m (47m), 66m 111m 444m 555m 45m (36m)"},

		// This is not possible
		// m111 m4(+m4) m444 m777 m888
		{"1114444777888m", "Invalid"},
		// Hand is ready
		{"555z333m888s66z55p1z", "Invalid"},
	} {
		t.Run(v.hand, func(t *testing.T) {
			assert.Equal(t, v.expected, tstTempai(t, v.hand))
		})
	}

}

func TestTempai9Gates(t *testing.T) {
	temp := testTempai(t, "1112345678999m")
	require.NotNil(t, temp)
	assert.Equal(t, 11, len(temp.Results))
	assert.Equal(t, "123456789m", GetWaits(temp).Tiles().String())
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
	for _, v := range []struct {
		hand     string
		expected string
	}{
		{"4555m123789s333z", "346m"},
		// Nobetan
		{"1234m123789s333z", "14m"},
		// Sanmenten
		{"23456m11p123456s", "147m"},
		// Sanmentan
		{"1234567m123456s", "147m"},
		// 8 sides
		{"1113334567888m", "23456789m"},
		{"1222234456678s", "1347s"},
		{"1123333445566p", "1247p"},
		{"2333445556777m", "1345m"},
		// 3 side-pons
		{"4555m123456789s", "346m"},
		// 2 side-pons
		{"3555m123456789s", "34m"},
		// 5-side
		{"4445678m456789s", "35689m"},
		// 4-side
		{"44456s66678p111z", "69p47s"},

		{"7788999s123456p", "6789s"},

		//ipeko
		{"22334455s11122z", "25s2z"},
		{"2233445566s123p", "2356s"},

		// kantankan 3335777
		{"3335777s111222z", "456s"},
		// tatsumaki
		{"3334555s111222z", "23456s"},
		// hapubidjin 2223456777
		{"2223456777s111z", "12345678s"},

		// kokushi
		{"19s19p1m12345677z", "9m"},
		{"19s19p19m1234567z", "19m19p19s1234567z"},

		// Difiicult
		// 7(6)
		{"1233334555678s", "124569s"},
		{"22223333444s55z", "14s5z"},

		// Interpret
		{"222p123s111555z4p", "34p"},

		// 6 is used
		{"1112345666678m", "12345789m"},
		// 4 is used
		{"2344445678999m", "12356789m"},
	} {
		t.Run(v.hand, func(t *testing.T) {
			res := testTempai(t, v.hand)
			require.NotNil(t, res)
			assert.Equal(t, v.expected, GetWaits(res).Tiles().String())
		})
	}
}

// func TestTempaiCorrectNumbers(t *testing.T) {
// 	// 555678m56788p678s
// 	ints := []int{17, 18, 20, 25, 31, 54, 59, 61, 66, 67, 95, 100, 102}
// 	require.Equal(t, 13, len(ints))

// 	validateInstances := func(winInstances tile.Instances) bool {
// 		winInts := make([]int, len(winInstances))
// 		for k, v := range winInstances {
// 			winInts[k] = int(v)
// 		}
// 		return assert.Equal(t, ints, winInts)
// 	}

// 	validate := func(melds meld.Melds) bool {
// 		i := getMeldsInstances(melds).Instances()
// 		return validateInstances(i)
// 	}

// 	instances := make(tile.Instances, len(ints))
// 	for k, v := range ints {
// 		instances[k] = tile.Instance(v)
// 	}
// 	compact := compact.NewInstances().Add(instances)
// 	require.True(t, validateInstances(compact.Instances()))
// 	temp := Calculate(compact)
// 	require.NotNil(t, temp)

// 	for k, variant := range temp {
// 		require.True(t, validate(variant), fmt.Sprintf("[%v] %v", k, variant))
// 	}
// }

func TestTempaiT1(t *testing.T) {
	temp := testTempai(t, "234m13999p567s11z")
	require.NotNil(t, temp)
	assert.Equal(t, 1, len(temp.Results))
	waits := GetWaits(temp).Tiles()
	assert.Equal(t, "2p", waits.String())
}

func tstTempai(t *testing.T, str string) string {
	obj := testTempai(t, str)
	if obj == nil {
		return "Invalid"
	}
	strs := make([]string, 0, len(obj.Results))
	for _, v := range obj.Results {
		strs = append(strs, DebugTempai(v))
	}
	sort.Strings(strs)

	return strings.Join(strs, ", ")
}

func testTempai(t *testing.T, str string) *TempaiResults {
	tg := compact.NewTileGenerator()
	inst, err := tg.CompactFromString(str)
	require.NoError(t, err, str)
	return Calculate(inst)
}

func testAway(t *testing.T, str string) compact.Tiles {
	tg := compact.NewTileGenerator()
	inst, err := tg.CompactFromString(str)
	require.NoError(t, err)
	return GetTempaiTiles(inst)
}

func testAwayString(t *testing.T, str string) string {
	return testAway(t, str).Tiles().String()
}
