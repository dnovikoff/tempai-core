package yaku

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/meld"
	"github.com/dnovikoff/tempai-core/shanten"
	"github.com/dnovikoff/tempai-core/tile"
)

type YakuTester struct {
	hand     compact.Instances
	declared meld.Melds
	ctx      *Context
	tg       *compact.Generator
	t        *testing.T
}

func NewYakuTester(t *testing.T, in string) (this *YakuTester) {
	this = &YakuTester{t: t}
	tg := compact.NewTileGenerator()
	hand, err := tg.CompactFromString(in)
	require.NoError(t, err)
	this.hand = hand
	this.ctx = &Context{
		Rules: &RulesEMA,
	}
	this.tg = tg
	return
}

func (this *YakuTester) Riichi() *YakuTester {
	this.ctx.IsRiichi = true
	return this
}

func (this *YakuTester) IDora(str string) *YakuTester {
	tiles, err := this.tg.InstancesFromString(str)
	require.NoError(this.t, err, str)
	this.ctx.DoraTiles = IndicatorsToDoraTiles(tiles)
	return this
}

func (this *YakuTester) IUra(str string) *YakuTester {
	tiles, err := this.tg.InstancesFromString(str)
	require.NoError(this.t, err, str)
	this.ctx.UraTiles = IndicatorsToDoraTiles(tiles)
	return this
}

func (this *YakuTester) tempai() shanten.TempaiMelds {
	cnt := this.hand.Count() + len(this.declared)*3
	require.Equal(this.t, 13, cnt, this.hand.Instances().String())
	t := shanten.CalculateTempai(this.hand, this.declared)
	require.NotEmpty(this.t, t)
	return t
}

func (this *YakuTester) win(t tile.Tile) *YakuResult {
	i := this.tempai().Index()
	this.ctx.Tile = this.tg.Instance(t)
	win := Win(i, this.ctx)
	return win
}

func (this *YakuTester) Ron(t tile.Tile) *YakuResult {
	this.ctx.IsTsumo = false
	return this.win(t)
}

func (this *YakuTester) TempaiTiles() string {
	return this.tempai().Waits().Tiles().String()
}

func (this *YakuTester) Tsumo(t tile.Tile) *YakuResult {
	this.ctx.IsTsumo = true
	return this.win(t)
}

func (this *YakuTester) Kan(t tile.Tile) {
	require.Equal(this.t, 4, this.hand.GetCount(t))
	kan := meld.NewKan(t, 0)
	kan.ExtractFrom(this.hand)
	this.declared = append(this.declared, kan.Meld())
}

func (this *YakuTester) Pon(t tile.Tile, o base.Opponent) {
	this.Declare(meld.NewPonPart(t, 0, 1), t, o)
}

func (this *YakuTester) Declare(m meld.Interface, t tile.Tile, o base.Opponent) {
	fixed := m.Rebase(this.hand)
	require.False(this.t, fixed.IsNull())
	i := this.tg.Instance(t)
	require.False(this.t, i.IsNull())
	opened := fixed.Interface().Open(i, o)
	require.False(this.t, fixed.IsNull())
	fixed.ExtractFrom(this.hand)
	this.declared = append(this.declared, opened)
}
