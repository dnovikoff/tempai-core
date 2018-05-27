package yaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/meld"
	"github.com/dnovikoff/tempai-core/tile"
)

func TestYakuNoOtherFuHand(t *testing.T) {
	tester := NewYakuTester(t, "12456789s56p11m")
	tester.Declare(meld.NewSeq(tile.Sou1, 0, 0, meld.HoleCopy), tile.Sou3, base.Left)
	win := tester.Ron(tile.Pin7)
	assert.Equal(t, "1 = YakuItsuu: 1", win.String())
	assert.Equal(t, "22 = 20(FuBase) + 2(FuNoOpenFu)", win.Fus.String())
}

func TestYakuFuOpenKanTerminal16(t *testing.T) {
	tester := NewYakuTester(t, "111s23p456789m22z")
	tester.Declare(meld.NewPon(tile.Sou1.Instance(0)), tile.Sou1, base.Left)
	assert.Equal(t, "14p", tester.TempaiTiles())
	// At least on han
	tester.ctx.IsRinshan = true
	win := tester.Tsumo(tile.Pin1)
	require.NotNil(t, win)
	assert.Equal(t, "38 = 20(FuBase) + 2(FuTsumo) + 16(FuOther)[1111s+]", win.Fus.String())
}

func TestYakuFuOpenKanDragon16(t *testing.T) {
	tester := NewYakuTester(t, "22555z23p456789m")
	tester.Declare(meld.NewPon(tile.White.Instance(0)), tile.White, base.Left)
	win := tester.Tsumo(tile.Pin1)
	assert.Equal(t, "38 = 20(FuBase) + 2(FuTsumo) + 16(FuOther)[5555z+]", win.Fus.String())
}

func TestYakuFuOpenKanWind16(t *testing.T) {
	tester := NewYakuTester(t, "11122z23p456789m")
	tester.Declare(meld.NewPon(tile.East.Instance(0)), tile.East, base.Left)
	win := tester.Tsumo(tile.Pin1)
	assert.Equal(t, "38 = 20(FuBase) + 2(FuTsumo) + 16(FuOther)[1111z+]", win.Fus.String())
}

func TestYakuFuClosedKanMiddle16(t *testing.T) {
	tester := NewYakuTester(t, "2222s23p456789m22z")
	tester.Kan(tile.Sou2)
	win := tester.Tsumo(tile.Pin1)
	assert.Equal(t, "38 = 20(FuBase) + 2(FuTsumo) + 16(FuOther)[2222s]", win.Fus.String())
}

func TestYakuFuOpenKanMiddle8(t *testing.T) {
	tester := NewYakuTester(t, "222s23p456789m22z")
	tester.Declare(meld.NewPon(tile.Sou2.Instance(0)), tile.Sou2, base.Left)
	tester.ctx.IsRinshan = true
	win := tester.Tsumo(tile.Pin1)
	assert.Equal(t, "30 = 20(FuBase) + 2(FuTsumo) + 8(FuOther)[2222s+]", win.Fus.String())
}

func TestYakuFuClosedKanTerminal32(t *testing.T) {
	tester := NewYakuTester(t, "1111s23p456789m22z")
	tester.Kan(tile.Sou1)
	win := tester.Tsumo(tile.Pin1)
	assert.Equal(t, "54 = 20(FuBase) + 2(FuTsumo) + 32(FuOther)[1111s]", win.Fus.String())
}

func TestYakuFuClosedKanDragon32(t *testing.T) {
	tester := NewYakuTester(t, "225555z23p456789m")
	tester.Kan(tile.White)
	win := tester.Tsumo(tile.Pin1)
	assert.Equal(t, "54 = 20(FuBase) + 2(FuTsumo) + 32(FuOther)[5555z]", win.Fus.String())
}

func TestYakuFuClosedKanWind32(t *testing.T) {
	tester := NewYakuTester(t, "111122z23p456789m")
	tester.Kan(tile.East)
	win := tester.Tsumo(tile.Pin1)
	assert.Equal(t, "54 = 20(FuBase) + 2(FuTsumo) + 32(FuOther)[1111z]", win.Fus.String())
}
