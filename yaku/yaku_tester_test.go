package yaku

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/calc"
	"github.com/dnovikoff/tempai-core/hand/tempai"
	"github.com/dnovikoff/tempai-core/tile"
)

type yakuTester struct {
	hand     compact.Instances
	declared calc.Melds
	ctx      *Context
	tg       *compact.Generator
	t        *testing.T
	rules    *RulesStruct
}

func newYakuTester(t *testing.T, in string) *yakuTester {
	yt := &yakuTester{t: t}
	tg := compact.NewTileGenerator()
	hand, err := tg.CompactFromString(in)
	require.NoError(t, err)
	yt.hand = hand
	yt.rules = RulesEMA()
	yt.ctx = &Context{
		Rules: yt.rules,
	}
	yt.tg = tg
	return yt
}

func (yt *yakuTester) riichi() *yakuTester {
	yt.ctx.IsRiichi = true
	return yt
}

func (yt *yakuTester) iDora(str string) *yakuTester {
	tiles, err := yt.tg.InstancesFromString(str)
	require.NoError(yt.t, err, str)
	yt.ctx.DoraTiles = IndicatorsToDoraTiles(tiles)
	return yt
}

func (yt *yakuTester) iUra(str string) *yakuTester {
	tiles, err := yt.tg.InstancesFromString(str)
	require.NoError(yt.t, err, str)
	yt.ctx.UraTiles = IndicatorsToDoraTiles(tiles)
	return yt
}

func (yt *yakuTester) tempai() *tempai.TempaiResults {
	cnt := yt.hand.CountBits() + len(yt.declared)*3
	require.Equal(yt.t, 13, cnt, yt.hand.Instances().String())
	t := tempai.Calculate(yt.hand, calc.Declared(yt.declared))
	require.NotNil(yt.t, t)
	return t
}

func (yt *yakuTester) win(t tile.Tile) *Result {
	i := yt.tempai()
	yt.ctx.Tile = yt.tg.Instance(t)
	return Win(i, yt.ctx, nil)
}

func (yt *yakuTester) ron(t tile.Tile) *Result {
	yt.ctx.IsTsumo = false
	return yt.win(t)
}

func (yt *yakuTester) tempaiTiles() string {
	return tempai.GetWaits(yt.tempai()).Tiles().String()
}

func (yt *yakuTester) tsumo(t tile.Tile) *Result {
	yt.ctx.IsTsumo = true
	return yt.win(t)
}

func (yt *yakuTester) kan(t tile.Tile) {
	yt.open(calc.Kan(t))
}

func (yt *yakuTester) kanClosed(t tile.Tile) {
	yt.declared = append(yt.declared, calc.Kan(t))
}

func (yt *yakuTester) pon(t tile.Tile) {
	yt.open(calc.Pon(t))
}

func (yt *yakuTester) open(m calc.Meld) {
	yt.declared = append(yt.declared, calc.Open(m))
}

func (yt *yakuTester) chi(t tile.Tile) {
	yt.open(calc.Chi(t))
}
