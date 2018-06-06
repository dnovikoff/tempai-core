package yaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tempai-core/tile"
)

func TestYakumanDaisangen(t *testing.T) {
	tester := NewYakuTester(t, "555666777z123p5s")
	win := tester.Ron(tile.Sou5)
	assert.Equal(t, "YakumanDaisangen:1", win.String())
}

func TestYakumanTenhouu(t *testing.T) {
	tester := NewYakuTester(t, "225588s11m14466z")
	tester.ctx.IsFirstTake = true
	win := tester.Tsumo(tile.East)
	assert.Equal(t, "YakumanTenhou:1", win.String())
}

func TestYakumanChuren9(t *testing.T) {
	tester := NewYakuTester(t, "1112345678999p")
	str := "YakumanChuurenpooto9:2"

	assert.Equal(t, str, tester.Ron(tile.Pin1).String(), tile.Pin1.String())
	assert.Equal(t, str, tester.Ron(tile.Pin2).String(), tile.Pin2.String())
	assert.Equal(t, str, tester.Ron(tile.Pin3).String(), tile.Pin3.String())
	assert.Equal(t, str, tester.Ron(tile.Pin4).String(), tile.Pin4.String())
	assert.Equal(t, str, tester.Ron(tile.Pin5).String(), tile.Pin5.String())
	assert.Equal(t, str, tester.Ron(tile.Pin6).String(), tile.Pin6.String())
	assert.Equal(t, str, tester.Ron(tile.Pin7).String(), tile.Pin7.String())
	assert.Equal(t, str, tester.Ron(tile.Pin8).String(), tile.Pin8.String())
	assert.Equal(t, str, tester.Ron(tile.Pin9).String(), tile.Pin9.String())
}

func TestYakumanChuren1(t *testing.T) {
	tester := NewYakuTester(t, "1112334678999p")
	assert.Equal(t, "235p", tester.TempaiTiles())

	win := tester.Ron(tile.Pin5)
	assert.Equal(t, "YakumanChuurenpooto:1", win.String())

	assert.Nil(t, tester.Ron(tile.Pin1))
	assert.Nil(t, tester.Ron(tile.Pin4))
	assert.Nil(t, tester.Ron(tile.Pin6))
	assert.Nil(t, tester.Ron(tile.Pin7))
	assert.Nil(t, tester.Ron(tile.Pin8))
	assert.Nil(t, tester.Ron(tile.Pin9))

	assert.Equal(t, "6 = YakuChinitsu: 6", tester.Ron(tile.Pin2).String())
	assert.Equal(t, "6 = YakuChinitsu: 6", tester.Ron(tile.Pin3).String())
}

func TestYakumanKokushi13(t *testing.T) {
	tester := NewYakuTester(t, "19s19p19m1234567z")
	for i := tile.TileBegin; i < tile.TileEnd; i++ {
		win := tester.Ron(i)
		if i.IsTerminalOrHonor() {
			require.NotNil(t, win, i.String())
			assert.Equal(t, "YakumanKokushi13:2", win.String())
		} else {
			assert.Nil(t, win, i.String(), i.String())
		}
	}
}

func TestYakumanKokushi1(t *testing.T) {
	tester := NewYakuTester(t, "19s19p1m12345677z")
	win := tester.Ron(tile.Man9)
	assert.Equal(t, "YakumanKokushi:1", win.String())

	winTile := tile.Man9
	for i := tile.TileBegin; i < tile.TileEnd; i++ {
		if i == winTile {
			continue
		}
		assert.Nil(t, tester.Ron(i), i.String())
	}
}

func TestYakumanSuuankou(t *testing.T) {
	tester := NewYakuTester(t, "55566z333m888s55p")
	win := tester.Tsumo(tile.Green)
	assert.Equal(t, "YakumanSuuankou:1", win.String())

	win = tester.Tsumo(tile.Pin5)
	assert.Equal(t, "YakumanSuuankou:1", win.String())

	win = tester.Ron(tile.Pin5)
	assert.Equal(t, "5 = YakuHaku: 1, YakuSanankou: 2, YakuToitoi: 2", win.String())
}

func TestYakumanSuuankouTanki(t *testing.T) {
	tester := NewYakuTester(t, "555666z333m888s5p")
	win := tester.Tsumo(tile.Pin5)
	assert.Equal(t, "YakumanSuuankouTanki:2", win.String())
	win = tester.Ron(tile.Pin5)
	assert.Equal(t, "YakumanSuuankouTanki:2", win.String())
}

func TestYakumanSosushi(t *testing.T) {
	tester := NewYakuTester(t, "1112223334z123m")
	win := tester.Ron(tile.North)
	assert.Equal(t, "YakumanShousuushi:1", win.String())
}

func TestYakumanDaisuushi(t *testing.T) {
	tester := NewYakuTester(t, "11122233344z11m")
	win := tester.Ron(tile.North)
	assert.Equal(t, "YakumanDaisuushi:2", win.String())
}

func TestYakumanCombo(t *testing.T) {
	tester := NewYakuTester(t, "1112223334445z")
	tester.ctx.IsFirstTake = true
	win := tester.Tsumo(tile.White)
	assert.Equal(t, "YakumanDaisuushi:2, YakumanSuuankouTanki:2, YakumanTenhou:1, YakumanTsuiisou:1", win.String())
}

func TestYakumanTTError(t *testing.T) {
	tester := NewYakuTester(t, "22m333p555777z88s")
	win := tester.Tsumo(tile.Sou8)
	assert.Equal(t, "YakumanSuuankou:1", win.String())
	// assert.Equal(t, "22m 333p 555z 777z 88s(OpponentSelf:8s:WIN) (8s)", win.Melds.String())
}

func TestYakumanGreenRules(t *testing.T) {
	tester := NewYakuTester(t, "2223344466688s")
	t.Run("no green not required", func(t *testing.T) {
		tester.rules.IsGreenRequired = false
		win := tester.Ron(tile.Sou3)
		assert.Equal(t, "YakumanRyuuiisou:1", win.String())
	})
	t.Run("no green not required", func(t *testing.T) {
		tester := NewYakuTester(t, "2223344466688s")
		tester.rules.IsGreenRequired = true
		win := tester.Ron(tile.Sou3)
		assert.Equal(t, "11 = YakuChinitsu: 6, YakuSanankou: 2, YakuTanyao: 1, YakuToitoi: 2", win.String())
	})

	tester = NewYakuTester(t, "22233444666s66z")
	t.Run("green not required", func(t *testing.T) {
		tester.rules.IsGreenRequired = false
		win := tester.Ron(tile.Sou3)
		assert.Equal(t, "YakumanRyuuiisou:1", win.String())
	})
	t.Run("green required", func(t *testing.T) {
		tester.rules.IsGreenRequired = true
		win := tester.Ron(tile.Sou3)
		assert.Equal(t, "YakumanRyuuiisou:1", win.String())
	})
}
