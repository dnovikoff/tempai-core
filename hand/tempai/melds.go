package tempai

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/meld"
	"github.com/dnovikoff/tempai-core/tile"
)

type TempaiMelds []meld.Melds
type IndexedResult map[tile.Tile]TempaiMelds

func (tm TempaiMelds) UsedTiles() compact.Tiles {
	if len(tm) == 0 {
		return 0
	}
	i := compact.NewInstances()
	for _, v := range tm {
		v.AddTo(i)
	}
	return i.GetFull()
}

func (tm TempaiMelds) Index() IndexedResult {
	if len(tm) == 0 {
		return nil
	}
	ret := make(IndexedResult)
	for _, m := range tm {
		meldsWaits(m).Each(func(k tile.Tile) bool {
			ret[k] = append(ret[k], m)
			return true
		})
	}
	return ret
}

func (tm TempaiMelds) Waits() compact.Tiles {
	used := tm.UsedTiles()
	waits := compact.Tiles(0)
	for _, v := range tm {
		waits |= meldsWaits(v).Sub(used)
	}
	return waits
}

func (tm IndexedResult) Waits() compact.Tiles {
	ret := compact.Tiles(0)
	for k := range tm {
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
