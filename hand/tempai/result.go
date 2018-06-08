package tempai

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/meld"
	"github.com/dnovikoff/tempai-core/tile"
)

type result struct {
	Melds TempaiMelds

	declared meld.Melds
}

func (*result) CheckMinuses(minuses int) bool {
	return true
}

func (r *result) Record(melds meld.Melds, tiles compact.Instances, totals compact.Totals) {
	if len(r.declared)+len(melds) != 4 {
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
	ret = append(ret, r.declared...)
	ret = append(ret, melds...)
	ret = append(ret, last)

	r.Melds = append(r.Melds, ret)
}

func getMeldsInstances(in meld.Melds) compact.Instances {
	ret := compact.NewInstances()
	for _, v := range in {
		v.AddTo(ret)
	}
	return ret
}
