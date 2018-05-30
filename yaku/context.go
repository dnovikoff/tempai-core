package yaku

import (
	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/tile"
)

// Context is createed to avoid using other game objects in calculations
type Context struct {
	Tile        tile.Instance
	SelfWind    base.Wind
	RoundWind   base.Wind
	DoraTiles   tile.Tiles
	UraTiles    tile.Tiles
	Rules       Rules
	IsTsumo     bool
	IsRiichi    bool
	IsIpatsu    bool
	IsDaburi    bool
	IsLastTile  bool
	IsRinshan   bool
	IsFirstTake bool
	IsChankan   bool
}

func (this Context) ShouldAddUras() bool {
	return this.IsRiichi && this.Rules.Ura()
}

func (this Context) ShouldAddIpatsu() bool {
	return this.IsIpatsu && this.Rules.Ipatsu()
}

func (this Context) IsRon() bool {
	return !this.IsTsumo
}

func IndicatorsToDoraTiles(in tile.Instances) tile.Tiles {
	return TileIndicatorsToDoraTiles(in.Tiles())
}

func TileIndicatorsToDoraTiles(in tile.Tiles) tile.Tiles {
	out := make(tile.Tiles, len(in))
	for k, v := range in {
		out[k] = v.Indicates()
	}
	return out
}
