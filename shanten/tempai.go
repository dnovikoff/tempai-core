package shanten

import (
	"bitbucket.org/dnovikoff/tempai-core/compact"
	"bitbucket.org/dnovikoff/tempai-core/meld"
	"bitbucket.org/dnovikoff/tempai-core/tile"
)

type TempaiMelds []meld.Melds

type TempaiResult struct {
	Melds TempaiMelds

	declared meld.Melds
}

type IndexedTempaiResult map[tile.Tile]TempaiMelds

func meldsWaits(m meld.Melds) compact.Tiles {
	tiles := compact.Tiles(0)
	for _, v := range m {
		tiles = tiles.Merge(v.Waits())
	}
	return tiles
}

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

func (this TempaiMelds) Index() IndexedTempaiResult {
	ret := make(IndexedTempaiResult)
	for _, m := range this {
		meldsWaits(m).Each(func(k tile.Tile) bool {
			ret[k] = append(ret[k], m)
			return true
		})
	}
	return ret
}

func (this IndexedTempaiResult) Waits() compact.Tiles {
	ret := compact.Tiles(0)
	for k, _ := range this {
		ret = ret.Set(k)
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

func (this *TempaiResult) CalculatePairs(tiles compact.Instances) {
	var wait meld.Meld
	pairs := make(meld.Melds, 0, 7)
	tiles.Each(func(m compact.Mask) bool {
		c := m.Count()
		if c == 1 {
			wait = meld.NewTanki(m.First()).Meld()
		} else if c > 1 {

			pairs = append(pairs, meld.NewPairFromMask(m).Meld())
		} else {
			return false
		}
		return true
	})
	if wait.IsNull() || len(pairs) != 6 {
		return
	}
	this.Melds = append(this.Melds, append(pairs, wait))
}

func (this *TempaiResult) Calculate13(tiles compact.Instances) {
	melds := make(meld.Melds, 0, 14)
	tanki := make(meld.Melds, 0, 14)

	var hole meld.Meld
	var pair meld.Meld
	for _, t := range tile.KokushiTiles {
		mask := tiles.GetMask(t)
		first := mask.First()
		switch mask.Count() {
		case 0:
			hole = meld.NewHole(t).Meld()
		case 1:
			tanki = append(tanki, meld.NewTanki(first).Meld())
			melds = append(melds, meld.NewOne(first).Meld())
		case 2:
			pair = meld.NewPairFromMask(mask).Meld()
		default:
			return
		}
	}

	if hole.IsNull() {
		if !pair.IsNull() || len(tanki) != 13 {
			return
		}
		this.Melds = append(this.Melds, tanki)
	} else {
		if pair.IsNull() || len(melds) != 11 {
			return
		}
		this.Melds = append(this.Melds, append(melds, pair, hole))
	}
}

func (this *TempaiResult) CheckMinuses(minuses int) bool {
	return true
}

func (this *TempaiResult) Record(melds meld.Melds, tiles compact.Instances, totals compact.Totals) {
	if len(this.declared)+len(melds) != 4 {
		return
	}
	last := meld.ExtractLastMeld(tiles)
	if last == 0 {
		return
	}

	// Validate
	if last.Waits().Each(func(t tile.Tile) bool {
		if !totals.IsFull(t) {
			return false
		}
		return true
	}) {
		return
	}

	ret := make(meld.Melds, 0, 5)
	ret = append(ret, this.declared...)
	ret = append(ret, melds...)
	ret = append(ret, last)

	this.Melds = append(this.Melds, ret)
}

func getMeldsInstances(in meld.Melds) compact.Instances {
	ret := compact.NewInstances()
	for _, v := range in {
		v.AddTo(ret)
	}
	return ret
}

func NewTempai(closed compact.Instances, declared meld.Melds) *Calculator {
	opened := len(declared)
	if opened*3+closed.Count() != 13 {
		return nil
	}

	return NewCalculator(meld.AllTempaiMelds, closed, getMeldsInstances(declared), len(declared))
}

func CalculateTempai(closed compact.Instances, declared meld.Melds) TempaiMelds {
	x := &TempaiResult{declared: declared}
	t := NewTempai(closed, declared)
	if t == nil {
		return nil
	}
	t.ResetResult(x)
	t.Calculate()
	if len(declared) == 0 {
		x.CalculatePairs(closed)
		x.Calculate13(closed)
	}
	return x.Melds
}

func CheckTempai(closed compact.Instances, declared meld.Melds) bool {
	// TODO: optimize
	x := CalculateTempai(closed, declared)
	return len(x) > 0
}

// TODO: solve with effectivity
func GetTempaiTiles(closed compact.Instances, declared meld.Melds) compact.Tiles {
	if len(declared)*3+closed.Count() != 14 {
		return 0
	}
	result := compact.Tiles(0)

	closed.EachTile(func(t tile.Tile) bool {
		i := closed.RemoveTile(t)
		if CheckTempai(closed, declared) {
			result = result.Set(t)
		}
		closed.Set(i)
		return true
	})
	return result
}
