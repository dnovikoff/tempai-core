package shanten

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/meld"
	"github.com/dnovikoff/tempai-core/tile"
)

type sSortId uint32

type ShantenResults struct {
	FormalyImprovedBy compact.Tiles
	RegularImproves   compact.Tiles
	UkeIre            compact.Totals
	Value             int
	Opened            int
	checked           compact.Tiles

	sortId sSortId
}

func NewShantenResults(opened int) *ShantenResults {
	return &ShantenResults{Value: 14, Opened: opened}
}

func NewSortId(uCount, tCount, value int) sSortId {
	val := (13 - sSortId(value))
	val = (val << 8) | sSortId(uCount)
	val = (val << 8) | sSortId(tCount)

	return val
}

func (this sSortId) BetterThan(other sSortId) bool {
	return this > other
}

func (this *ShantenResults) SortId() sSortId {
	return this.sortId
}

func (this *ShantenResults) Clone() *ShantenResults {
	return &ShantenResults{
		this.FormalyImprovedBy,
		this.RegularImproves,
		this.UkeIre.Clone(),
		this.Value,
		this.Opened,
		0,
		this.sortId,
	}
}

func (this *ShantenResults) Recount(total compact.Totals) *ShantenResults {
	uke := compact.NewTotals()
	uCount := 0
	tCount := 0
	this.FormalyImprovedBy.Each(func(t tile.Tile) bool {
		c := 4 - total.Get(t)
		if c <= 0 {
			return true
		}
		uke.Set(t, c)
		uCount++
		tCount += c
		return true
	})
	this.UkeIre = uke

	val := (13 - sSortId(this.Value))
	val = (val << 8) | sSortId(uCount)
	val = (val << 8) | sSortId(tCount)

	this.sortId = NewSortId(uCount, tCount, this.Value)
	return this
}

func (this *ShantenResults) CheckMinuses(minuses int) bool {
	// A set without 1 tile is still tempai
	return minuses <= this.Value+1
}

func reducesShantenBy(m meld.Meld) int {
	if m == 0 {
		return 0
	}
	if m.Type() == meld.TypePair {
		if meld.Pair(m).IsComplete() {
			return 1
		} else {
			return 0
		}
	}
	if m.IsComplete() {
		return 2
	}
	return 1
}

func (this *ShantenResults) Record(melds meld.Melds, tiles compact.Instances, totals compact.Totals) {
	sets := this.Opened
	value := 8 - sets*2
	var improves compact.Tiles
	havePair := false
	for _, v := range melds {
		if v.Type() == meld.TypePair {
			havePair = true
		} else {
			sets++
		}
		value -= reducesShantenBy(v)
		improves |= v.Waits()
	}

	if value > this.Value {
		// avoid useless calculations
		return
	} else if value < this.Value {
		this.checked = 0
	}

	fullSets := sets > 3

	if !havePair {
		tiles.Each(func(m compact.Mask) bool {
			if !totals.IsFull(m.Tile()) {
				improves = improves.Set(m.Tile())
			}
			return true
		})
	}

	if !fullSets {
		toCheck := tiles.UniqueTiles() & (^this.checked)
		try := func(central, improve, wait tile.Tile) {
			if (central.Type() != improve.Type()) || (central.Type() != wait.Type()) {
				return
			}

			if totals.IsFull(wait) || totals.IsFull(improve) {
				return
			}
			improves = improves.Set(improve)
		}
		toCheck.EachRange(tile.Begin, tile.SequenceEnd, func(t tile.Tile) bool {
			try(t, t-2, t-1)
			try(t, t-1, t-2)
			try(t, t-1, t+1)
			try(t, t+1, t-1)
			try(t, t+1, t+2)
			try(t, t+2, t+1)
			return true
		})
		this.checked |= toCheck
	}

	this.Add(value, improves, true)
}

func (this *ShantenResults) Add(value int, improves compact.Tiles, isReg bool) {
	if improves == 0 {
		return
	}
	if value > this.Value {
		return
	} else if value < this.Value {
		this.FormalyImprovedBy = 0
		this.RegularImproves = 0
		this.Value = value
		// this.checked = 0
	}
	this.FormalyImprovedBy |= improves
	if isReg {
		this.RegularImproves |= improves
	}
}

func NewShantentCalculator(tiles compact.Instances, used compact.Instances, opened int) *Calculator {
	return NewCalculator(meld.AllShantenMelds, tiles, used, opened)
}

func (this *ShantenResults) CalculatePairs(tiles compact.Instances) {
	value := 6
	improves := compact.Tiles(0)
	tiles.Each(func(m compact.Mask) bool {
		c := m.Count()
		if c == 1 {
			improves = improves.Set(m.Tile())
		} else if c > 1 {
			value--
		}
		return true
	})
	this.Add(value, improves, false)
}

func (this *ShantenResults) Calculate13(tiles compact.Instances) {
	havePair := false
	count := 0
	missing := compact.Tiles(0)
	for _, t := range tile.KokushiTiles {
		switch tiles.GetMask(t).Count() {
		case 0:
			missing = missing.Set(t)
		case 1:
			count++
		default:
			count++
			havePair = true
		}
	}

	if !havePair {
		missing = compact.KokushiTiles
	} else {
		count++
	}
	this.Add(13-count, missing, false)
}

func CalculateShantenByMelds(tiles compact.Instances, melds meld.Melds) *ShantenResults {
	used := compact.NewInstances()
	melds.AddTo(used)
	return CalculateShanten(tiles, len(melds), used)
}

func CalculateShanten(tiles compact.Instances, opened int, used compact.Instances) *ShantenResults {
	results := NewShantenResults(opened)
	calc := NewShantentCalculator(tiles, used, opened)
	calc.ResetResult(results)
	calc.Calculate()
	if opened == 0 {
		results.CalculatePairs(tiles)
		results.Calculate13(tiles)
	}
	return results.Recount(calc.totals)
}

func CalculateShantenBoth(tiles compact.Instances, opened int, used compact.Instances) (all *ShantenResults, regular *ShantenResults) {
	all = NewShantenResults(opened)
	calc := NewShantentCalculator(tiles, used, opened)
	calc.ResetResult(all)
	calc.Calculate()
	regular = all.Clone().Recount(calc.totals)
	if opened == 0 {
		all.CalculatePairs(tiles)
		all.Calculate13(tiles)
	}
	all = all.Recount(calc.totals)
	return
}
