package yaku

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/calc"
	"github.com/dnovikoff/tempai-core/hand/tempai"
	"github.com/dnovikoff/tempai-core/meld"
	"github.com/dnovikoff/tempai-core/tile"
)

type yakuTester struct {
	hand     compact.Instances
	declared meld.Melds
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

func (yt *yakuTester) tempai() tempai.TempaiMelds {
	cnt := yt.hand.Count() + len(yt.declared)*3
	require.Equal(yt.t, 13, cnt, yt.hand.Instances().String())
	t := tempai.Calculate(yt.hand, calc.Melds(yt.declared))
	require.NotEmpty(yt.t, t)
	return t
}

func (yt *yakuTester) win(t tile.Tile) *Result {
	i := yt.tempai().Index()
	yt.ctx.Tile = yt.tg.Instance(t)
	win := Win(i, yt.ctx)
	return win
}

func (yt *yakuTester) ron(t tile.Tile) *Result {
	yt.ctx.IsTsumo = false
	return yt.win(t)
}

func (yt *yakuTester) tempaiTiles() string {
	return yt.tempai().Waits().Tiles().String()
}

func (yt *yakuTester) tsumo(t tile.Tile) *Result {
	yt.ctx.IsTsumo = true
	return yt.win(t)
}

func (yt *yakuTester) kan(t tile.Tile) {
	require.Equal(yt.t, 4, yt.hand.GetCount(t))
	kan := meld.NewKan(t.Instance(0))
	kan.ExtractFrom(yt.hand)
	yt.declared = append(yt.declared, kan.Meld())
}

func (yt *yakuTester) pon(t tile.Tile, o base.Opponent) {
	yt.declare(meld.NewPonPart(t, 0, 1), t, o)
}

func (yt *yakuTester) declare(m meld.Interface, t tile.Tile, o base.Opponent) {
	fixed := m.Rebase(yt.hand)
	require.False(yt.t, fixed.IsNull())
	i := yt.tg.Instance(t)
	require.NotEqual(yt.t, tile.InstanceNull, i)
	opened := fixed.Interface().Open(i, o)
	require.False(yt.t, fixed.IsNull())
	fixed.ExtractFrom(yt.hand)
	yt.declared = append(yt.declared, opened)
}
