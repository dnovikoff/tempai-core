package yaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/tile"
)

func TestYakuRinshanWin(t *testing.T) {
	tester := newYakuTester(t, "123p345s2255z")
	tester.ctx.SelfWind = base.WindSouth
	tester.kan(tile.Sou9)
	tester.ctx.IsRinshan = true
	win := tester.tsumo(tile.South)
	assert.Equal(t, "48 = 20(FuBase) + 2(FuTsumo) + 2(FuPair)[55z] + 16(FuMeld)[9999s+] + 8(FuMeld)[222z]", win.Fus.String())
	assert.Equal(t, "2 = YakuNanSelf: 1, YakuRinshan: 1", win.String())
	assert.Equal(t, 0, len(win.Yakuman))

	tester.ctx.IsRinshan = false
	win = tester.tsumo(tile.South)
	assert.Equal(t, "1 = YakuNanSelf: 1", win.String())
}

func TestYakuHoiteiWin(t *testing.T) {
	tester := newYakuTester(t, "123p345s2255z")
	tester.kan(tile.Sou9)
	tester.ctx.IsLastTile = true
	assert.Equal(t, "1 = YakuHoutei: 1", tester.ron(tile.South).String())
	assert.Equal(t, "1 = YakuHaitei: 1", tester.tsumo(tile.South).String())
}

func TestYakuRinshanIsNotHoitei(t *testing.T) {
	tester := newYakuTester(t, "123p555s2255z")
	tester.ctx.IsLastTile = true
	tester.ctx.IsRinshan = true
	tester.rules.IsHaiteiFromLiveOnly = true
	tester.kan(tile.Sou9)

	win := tester.tsumo(tile.South)
	assert.Equal(t, "1 = YakuRinshan: 1", win.String())

	tester.rules.IsHaiteiFromLiveOnly = false
	win = tester.tsumo(tile.South)
	assert.Equal(t, "2 = YakuHaitei: 1, YakuRinshan: 1", win.String())
}

func TestYakuKuitan(t *testing.T) {
	tester := newYakuTester(t, "234567p3388m")
	tester.rules.IsOpenTanyao = true
	tester.chi(tile.Sou2)
	win := tester.tsumo(tile.Man3)
	require.NotNil(t, win)
	assert.Equal(t, "1 = YakuTanyao: 1", win.String())
	assert.Equal(t, "26 = 20(FuBase) + 2(FuTsumo) + 4(FuMeld)[333m]", win.Fus.String())
}

func TestToitoiSananko(t *testing.T) {
	tester := newYakuTester(t, "55566z333m888s55p")
	win := tester.ron(tile.Pin5)
	assert.Equal(t, "5 = YakuHaku: 1, YakuSanankou: 2, YakuToitoi: 2", win.String())
}

func TestIppeiko(t *testing.T) {
	tester := newYakuTester(t, "22456m55566777p")
	win := tester.ron(tile.Pin6)
	assert.Equal(t, "2 = YakuIppeiko: 1, YakuTanyao: 1", win.String())
}

func TestRyanpeiko(t *testing.T) {
	tester := newYakuTester(t, "1122335566799s")
	win := tester.ron(tile.Sou7)
	assert.Equal(t, "10 = YakuChinitsu: 6, YakuPinfu: 1, YakuRyanpeikou: 3", win.String())
	// assert.Equal(t, "99s 123s 123s 567s 56s(OpponentRight:7s:WIN) (47s)", win.Melds.String())
}

func TestPinfuOtherWind(t *testing.T) {
	tester := newYakuTester(t, "33z123m456p66778s")
	win := tester.tsumo(tile.Sou5)
	assert.Equal(t, "2 = YakuPinfu: 1, YakuTsumo: 1", win.String())
	assert.Equal(t, "20 = 20(FuBase)", win.Fus.String())
	// assert.Equal(t, "33z 123m 456p 678s 67s(OpponentSelf:5s:WIN) (58s)", win.Melds.String())
}

func TestPinfuDragonPair(t *testing.T) {
	tester := newYakuTester(t, "77z123m456p66778s")
	win := tester.tsumo(tile.Sou5)
	assert.Equal(t, "1 = YakuTsumo: 1", win.String())
	assert.Equal(t, "24 = 20(FuBase) + 2(FuTsumo) + 2(FuPair)[77z]", win.Fus.String())
}

func TestPinfuBadWindSelfPair(t *testing.T) {
	tester := newYakuTester(t, "33z123m456p66778s")
	tester.ctx.SelfWind = base.WindWest
	win := tester.tsumo(tile.Sou5)
	assert.Equal(t, "1 = YakuTsumo: 1", win.String())
	assert.Equal(t, "24 = 20(FuBase) + 2(FuTsumo) + 2(FuPair)[33z]", win.Fus.String())
}

func TestPinfuBadWindRoundPair(t *testing.T) {
	tester := newYakuTester(t, "33z123m456p66778s")
	tester.ctx.SelfWind = base.WindWest
	win := tester.tsumo(tile.Sou5)
	assert.Equal(t, "1 = YakuTsumo: 1", win.String())
	assert.Equal(t, "24 = 20(FuBase) + 2(FuTsumo) + 2(FuPair)[33z]", win.Fus.String())
}

func TestPinfuBadWindBoth(t *testing.T) {
	tester := newYakuTester(t, "33z123m456p66778s")
	tester.ctx.SelfWind = base.WindWest
	tester.ctx.RoundWind = base.WindWest
	win := tester.tsumo(tile.Sou5)
	assert.Equal(t, "1 = YakuTsumo: 1", win.String())
	assert.Equal(t, "26 = 20(FuBase) + 2(FuTsumo) + 4(FuPair)[33z]", win.Fus.String())
}

func TestPinfuPenchan(t *testing.T) {
	tester := newYakuTester(t, "33z12m456p566778s")
	win := tester.tsumo(tile.Man3)
	assert.Equal(t, "1 = YakuTsumo: 1", win.String())
	assert.Equal(t, "24 = 20(FuBase) + 2(FuTsumo) + 2(FuBadWait)[12m (3m)]", win.Fus.String())
}

func TestPinfuKanchan(t *testing.T) {
	tester := newYakuTester(t, "33z13m456p566778s")
	win := tester.tsumo(tile.Man2)
	assert.Equal(t, "1 = YakuTsumo: 1", win.String())
	assert.Equal(t, "24 = 20(FuBase) + 2(FuTsumo) + 2(FuBadWait)[13m (2m)]", win.Fus.String())
}

func TestNoFuOpen(t *testing.T) {
	tester := newYakuTester(t, "33z123p12367s")
	tester.chi(tile.Man1)
	win := tester.ron(tile.Sou8)
	assert.Equal(t, "1 = YakuSanshoku: 1", win.String())
	assert.Equal(t, "22 = 20(FuBase) + 2(FuNoOpenFu)", win.Fus.String())
}

func TestTanyaoClosed(t *testing.T) {
	tester := newYakuTester(t, "33s45p222666m777s")
	win := tester.tsumo(tile.Pin6)
	require.NotNil(t, win)
	assert.Equal(t, "4 = YakuSanankou: 2, YakuTanyao: 1, YakuTsumo: 1", win.String())
	assert.Equal(t, "34 = 20(FuBase) + 2(FuTsumo) + 4(FuMeld)[222m] + 4(FuMeld)[666m] + 4(FuMeld)[777s]", win.Fus.String())
}

func TestTanyaoOpen(t *testing.T) {
	tester := newYakuTester(t, "33s45p222666m")
	tester.pon(tile.Sou7)
	win := tester.tsumo(tile.Pin6)
	require.NotNil(t, win)
	assert.Equal(t, "1 = YakuTanyao: 1", win.String())
	assert.Equal(t, "32 = 20(FuBase) + 2(FuTsumo) + 4(FuMeld)[222m] + 4(FuMeld)[666m] + 2(FuMeld)[777s+]", win.Fus.String())
}

func TestYakuChiitoi(t *testing.T) {
	tester := newYakuTester(t, "225588s11m14455z")
	win := tester.tsumo(tile.East)
	assert.Equal(t, "3 = YakuChiitoi: 2, YakuTsumo: 1", win.String())
	assert.Equal(t, "25 = 25(FuBase7)", win.Fus.String())
}

func TestRуnho(t *testing.T) {
	tester := newYakuTester(t, "123456s234p55z99m")
	tester.ctx.IsFirstTake = true
	win := tester.ron(tile.Man9)
	assert.Equal(t, "5 = YakuRenhou: 5", win.String())
	tester.ctx.IsFirstTake = false

	win = tester.ron(tile.Man9)
	require.Nil(t, win)
}

func TestRуnhoCoolHandPriority(t *testing.T) {
	tester := newYakuTester(t, "1122335566799p")
	win := tester.ron(tile.Pin7)
	assert.Equal(t, "10 = YakuChinitsu: 6, YakuPinfu: 1, YakuRyanpeikou: 3", win.String())
}

func TestYakuCase1(t *testing.T) {
	tester := newYakuTester(t, "77p33777z")
	tester.pon(tile.Sou2)
	tester.pon(tile.Sou4)

	win := tester.ron(tile.Sou2)
	require.Nil(t, win)
	win = tester.ron(tile.Pin7)
	assert.Equal(t, "3 = YakuChun: 1, YakuToitoi: 2", win.String())
}

func TestTankiCase(t *testing.T) {
	tester := newYakuTester(t, "1233455m666z")
	tester.kanClosed(tile.Pin1)
	win := tester.tsumo(tile.Man5)
	assert.Equal(t, "2 = YakuHatsu: 1, YakuTsumo: 1", win.String())
	assert.Equal(t, "64 = 20(FuBase) + 2(FuTsumo) + 2(FuBadWait)[5m (5m)] + 8(FuMeld)[666z] + 32(FuMeld)[1111p]", win.Fus.String())
}

func TestTankiCase2(t *testing.T) {
	tester := newYakuTester(t, "1233455m234789s")
	win := tester.tsumo(tile.Man5)
	assert.Equal(t, "2 = YakuPinfu: 1, YakuTsumo: 1", win.String())
	assert.Equal(t, "20 = 20(FuBase)", win.Fus.String())
}

func TestOpenTanyao(t *testing.T) {
	tester := newYakuTester(t, "4588p444s")
	tester.chi(tile.Man4)
	tester.chi(tile.Man5)
	win := tester.ron(tile.Pin3)
	assert.Equal(t, "1 = YakuTanyao: 1", win.String())
	assert.Equal(t, "24 = 20(FuBase) + 4(FuMeld)[444s]", win.Fus.String())
}

func TestKanchanTest(t *testing.T) {
	tester := newYakuTester(t, "45556m456s456p11z")
	win := tester.tsumo(tile.Man5)
	assert.Equal(t, "3 = YakuSanshoku: 2, YakuTsumo: 1", win.String())
	assert.Equal(t, "32 = 20(FuBase) + 2(FuTsumo) + 2(FuBadWait)[46m (5m)] + 4(FuPair)[11z] + 4(FuMeld)[555m]", win.Fus.String())
}

// TODO:
// nagashi test
// steal tests
// openings tests
// 13 draw test
