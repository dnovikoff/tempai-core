package meld

import (
	"bitbucket.org/dnovikoff/tempai-core/compact"
	"bitbucket.org/dnovikoff/tempai-core/tile"
)

type BaseMelds Melds

func (this BaseMelds) Filter(possible compact.Instances, freeTiles compact.Tiles) Melds {
	filtered := make(Melds, 0, len(this))
	for _, v := range this {
		w := v.Waits()
		if !w.IsEmpty() && (w & freeTiles).IsEmpty() {
			continue
		}
		if v.Rebase(possible) != 0 {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func createAllBase(args ...func(i tile.Tile) Meld) (ret BaseMelds) {
	push := func(meld Meld) {
		if meld == 0 {
			return
		}
		ret = append(ret, meld)
	}
	for _, f := range args {
		for i := tile.Begin; i < tile.End; i++ {
			push(f(i))
		}
	}
	return
}

func createMergedBase(lhs, rhs BaseMelds) BaseMelds {
	x := make(BaseMelds, 0, len(lhs)+len(rhs))
	x = append(x, lhs...)
	x = append(x, rhs...)
	return x
}

func baseChi(t tile.Tile) Meld {
	return NewSeq(t, 0, 0, 0).Meld()
}

func baseChi2(t tile.Tile) Meld {
	return NewSeq(t, 0, 0, HoleCopy).Meld()
}

func baseKanchan(t tile.Tile) Meld {
	return NewSeq(t, 0, HoleCopy, 0).Meld()
}

func basePon(t tile.Tile) Meld {
	return NewPon(t, 3).Meld()
}

func basePonPart(t tile.Tile) Meld {
	return NewPonPart(t, 2, 3).Meld()
}

var AllPonMelds = createAllBase(basePon)
var AllTempaiMelds = createAllBase(baseChi, basePon)

// Seems like ryanment must be last, because of perf issues
var AllShantenMelds = createAllBase(baseChi, basePon, basePonPart, baseKanchan, baseChi2)
var allPartMelds = createAllBase(basePonPart, baseKanchan, baseChi2)
var AllPartWithKanMelds = createMergedBase(AllPonMelds, allPartMelds)
