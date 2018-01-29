package tempai

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/meld"
	"github.com/dnovikoff/tempai-core/tile"
)

type TempaiMelds []meld.Melds
type IndexedResult map[tile.Tile]TempaiMelds

func (this TempaiMelds) UsedTiles() compact.Tiles {
	if len(this) == 0 {
		return 0
	}
	i := compact.NewInstances()
	for _, v := range this {
		v.AddTo(i)
	}
	return i.GetFull()
}

func (this TempaiMelds) Index() IndexedResult {
	ret := make(IndexedResult)
	for _, m := range this {
		meldsWaits(m).Each(func(k tile.Tile) bool {
			ret[k] = append(ret[k], m)
			return true
		})
	}
	return ret
}

func (this TempaiMelds) Waits() compact.Tiles {
	used := this.UsedTiles()
	waits := compact.Tiles(0)
	for _, v := range this {
		waits |= meldsWaits(v).Sub(used)
	}
	return waits
}

func (this IndexedResult) Waits() compact.Tiles {
	ret := compact.Tiles(0)
	for k, _ := range this {
		ret = ret.Set(k)
	}
	return ret
}

func meldsWaits(m meld.Melds) compact.Tiles {
	tiles := compact.Tiles(0)
	for _, v := range m {
		tiles = tiles.Merge(v.Waits())
	}
	return tiles
}
