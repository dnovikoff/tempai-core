package yaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tempai-core/tile"
)

func TestYakuNoOtherFuHand(t *testing.T) {
	tester := newYakuTester(t, "11m56p456789s")
	tester.chi(tile.Sou1)
	win := tester.ron(tile.Pin7)
	assert.Equal(t, "1 = YakuItsuu: 1", win.String())
	assert.Equal(t, "22 = 20(FuBase) + 2(FuNoOpenFu)", win.Fus.String())
}

func TestYakuFuOpenKanTerminal16(t *testing.T) {
	tester := newYakuTester(t, "23p456789m22z")
	tester.kan(tile.Sou1)
	assert.Equal(t, "14p", tester.tempaiTiles())
	// At least on han
	tester.ctx.IsRinshan = true
	win := tester.tsumo(tile.Pin1)
	require.NotNil(t, win)
	assert.Equal(t, "38 = 20(FuBase) + 2(FuTsumo) + 16(FuMeld)[1111s+]", win.Fus.String())
}

func TestYakuFuOpenKanDragon16(t *testing.T) {
	tester := newYakuTester(t, "22z23p456789m")
	tester.kan(tile.White)
	win := tester.tsumo(tile.Pin1)
	assert.Equal(t, "38 = 20(FuBase) + 2(FuTsumo) + 16(FuMeld)[5555z+]", win.Fus.String())
}

func TestYakuFuOpenKanWind16(t *testing.T) {
	tester := newYakuTester(t, "22z23p456789m")
	tester.kan(tile.East)
	win := tester.tsumo(tile.Pin1)
	assert.Equal(t, "38 = 20(FuBase) + 2(FuTsumo) + 16(FuMeld)[1111z+]", win.Fus.String())
}

func TestYakuFuClosedKanMiddle16(t *testing.T) {
	tester := newYakuTester(t, "23p456789m22z")
	tester.kanClosed(tile.Sou2)
	win := tester.tsumo(tile.Pin1)
	assert.Equal(t, "38 = 20(FuBase) + 2(FuTsumo) + 16(FuMeld)[2222s]", win.Fus.String())
}

func TestYakuFuOpenKanMiddle8(t *testing.T) {
	tester := newYakuTester(t, "23p456789m22z")
	tester.kan(tile.Sou2)
	tester.ctx.IsRinshan = true
	win := tester.tsumo(tile.Pin1)
	assert.Equal(t, "30 = 20(FuBase) + 2(FuTsumo) + 8(FuMeld)[2222s+]", win.Fus.String())
}

func TestYakuFuClosedKanTerminal32(t *testing.T) {
	tester := newYakuTester(t, "23p456789m22z")
	tester.kanClosed(tile.Sou1)
	win := tester.tsumo(tile.Pin1)
	assert.Equal(t, "54 = 20(FuBase) + 2(FuTsumo) + 32(FuMeld)[1111s]", win.Fus.String())
}

func TestYakuFuClosedKanDragon32(t *testing.T) {
	tester := newYakuTester(t, "22z23p456789m")
	tester.kanClosed(tile.White)
	win := tester.tsumo(tile.Pin1)
	assert.Equal(t, "54 = 20(FuBase) + 2(FuTsumo) + 32(FuMeld)[5555z]", win.Fus.String())
}

func TestYakuFuClosedKanWind32(t *testing.T) {
	tester := newYakuTester(t, "22z23p456789m")
	tester.kanClosed(tile.East)
	win := tester.tsumo(tile.Pin1)
	assert.Equal(t, "54 = 20(FuBase) + 2(FuTsumo) + 32(FuMeld)[1111z]", win.Fus.String())
}

func TestYakuRinshanRules(t *testing.T) {
	tester := newYakuTester(t, "22z23p456789m")
	tester.kanClosed(tile.East)
	tester.ctx.IsRinshan = true

	tester.rules.IsRinshanFu = true
	win := tester.tsumo(tile.Pin1)
	assert.Equal(t, "54 = 20(FuBase) + 2(FuTsumo) + 32(FuMeld)[1111z]", win.Fus.String())

	tester.rules.IsRinshanFu = false
	win = tester.tsumo(tile.Pin1)
	assert.Equal(t, "52 = 20(FuBase) + 32(FuMeld)[1111z]", win.Fus.String())
}
