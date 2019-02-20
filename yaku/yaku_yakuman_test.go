package yaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/tile"
)

func TestYakumanDaisangen(t *testing.T) {
	tester := newYakuTester(t, "555666777z123p5s")
	win := tester.ron(tile.Sou5)
	assert.Equal(t, "YakumanDaisangen", win.String())
}

func TestYakumanTenhouu(t *testing.T) {
	tester := newYakuTester(t, "225588s11m14466z")
	tester.ctx.IsFirstTake = true
	win := tester.tsumo(tile.East)
	assert.Equal(t, "YakumanTenhou", win.String())
}

func TestYakumanChuren9(t *testing.T) {
	for tl := tile.Pin1; tl <= tile.Pin9; tl++ {
		t.Run(tl.String(), func(t *testing.T) {
			tester := newYakuTester(t, "1112345678999p")
			assert.Equal(t, "YakumanChuurenpooto9", tester.ron(tl).String())
		})
	}
}

func TestYakumanChuren1(t *testing.T) {
	tester := newYakuTester(t, "1112334678999p")
	assert.Equal(t, "235p", tester.tempaiTiles())

	win := tester.ron(tile.Pin5)
	assert.Equal(t, "YakumanChuurenpooto", win.String())

	assert.Nil(t, tester.ron(tile.Pin1))
	assert.Nil(t, tester.ron(tile.Pin4))
	assert.Nil(t, tester.ron(tile.Pin6))
	assert.Nil(t, tester.ron(tile.Pin7))
	assert.Nil(t, tester.ron(tile.Pin8))
	assert.Nil(t, tester.ron(tile.Pin9))

	assert.Equal(t, "6 = YakuChinitsu: 6", tester.ron(tile.Pin2).String())
	assert.Equal(t, "6 = YakuChinitsu: 6", tester.ron(tile.Pin3).String())
}

func TestYakumanKokushi13(t *testing.T) {
	tester := newYakuTester(t, "19s19p19m1234567z")
	for i := tile.TileBegin; i < tile.TileEnd; i++ {
		t.Run(i.String(), func(t *testing.T) {
			win := tester.ron(i)
			if compact.TerminalOrHonor.Check(i) {
				require.NotNil(t, win)
				assert.Equal(t, "YakumanKokushi13", win.String())
			} else {
				assert.Nil(t, win, i.String())
			}
		})
	}
}

func TestYakumanKokushiOne(t *testing.T) {
	tester := newYakuTester(t, "19s19p1m12345677z")
	winTile := tile.Man9
	for i := tile.TileBegin; i < tile.TileEnd; i++ {
		t.Run(i.String(), func(t *testing.T) {
			win := tester.ron(i)
			if i == winTile {
				require.NotNil(t, win)
				assert.Equal(t, "YakumanKokushi", win.String())
			} else {
				assert.Nil(t, win)
			}
		})
	}
}

func TestYakumanSuuankou(t *testing.T) {
	tester := newYakuTester(t, "55566z333m888s55p")
	win := tester.tsumo(tile.Green)
	assert.Equal(t, "YakumanSuuankou", win.String())

	win = tester.tsumo(tile.Pin5)
	assert.Equal(t, "YakumanSuuankou", win.String())

	win = tester.ron(tile.Pin5)
	assert.Equal(t, "5 = YakuHaku: 1, YakuSanankou: 2, YakuToitoi: 2", win.String())
}

func TestYakumanSuuankouTanki(t *testing.T) {
	tester := newYakuTester(t, "555666z333m888s5p")
	win := tester.tsumo(tile.Pin5)
	assert.Equal(t, "YakumanSuuankouTanki", win.String())
	win = tester.ron(tile.Pin5)
	assert.Equal(t, "YakumanSuuankouTanki", win.String())
}

func TestYakumanSosushi(t *testing.T) {
	tester := newYakuTester(t, "1112223334z123m")
	win := tester.ron(tile.North)
	assert.Equal(t, "YakumanShousuushi", win.String())
}

func TestYakumanDaisuushi(t *testing.T) {
	tester := newYakuTester(t, "11122233344z11m")
	win := tester.ron(tile.North)
	assert.Equal(t, "YakumanDaisuushi", win.String())
}

func TestYakumanCombo(t *testing.T) {
	tester := newYakuTester(t, "1112223334445z")
	tester.ctx.IsFirstTake = true
	win := tester.tsumo(tile.White)
	assert.Equal(t, "YakumanDaisuushi, YakumanSuuankouTanki, YakumanTenhou, YakumanTsuiisou", win.String())
}

func TestYakumanTTError(t *testing.T) {
	tester := newYakuTester(t, "22m333p555777z88s")
	win := tester.tsumo(tile.Sou8)
	assert.Equal(t, "YakumanSuuankou", win.String())
	// assert.Equal(t, "22m 333p 555z 777z 88s(OpponentSelf:8s:WIN) (8s)", win.Melds.String())
}

func TestYakumanGreenRules(t *testing.T) {
	tester := newYakuTester(t, "2223344466688s")
	t.Run("no green not required", func(t *testing.T) {
		tester.rules.IsGreenRequired = false
		win := tester.ron(tile.Sou3)
		assert.Equal(t, "YakumanRyuuiisou", win.String())
	})
	t.Run("no green not required", func(t *testing.T) {
		tester := newYakuTester(t, "2223344466688s")
		tester.rules.IsGreenRequired = true
		win := tester.ron(tile.Sou3)
		assert.Equal(t, "11 = YakuChinitsu: 6, YakuSanankou: 2, YakuTanyao: 1, YakuToitoi: 2", win.String())
	})

	tester = newYakuTester(t, "22233444666s66z")
	t.Run("green not required", func(t *testing.T) {
		tester.rules.IsGreenRequired = false
		win := tester.ron(tile.Sou3)
		assert.Equal(t, "YakumanRyuuiisou", win.String())
	})
	t.Run("green required", func(t *testing.T) {
		tester.rules.IsGreenRequired = true
		win := tester.ron(tile.Sou3)
		assert.Equal(t, "YakumanRyuuiisou", win.String())
	})
}
