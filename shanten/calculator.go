package shanten

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/meld"
	"github.com/dnovikoff/tempai-core/tile"
)

type CalculateResults interface {
	CheckMinuses(minuses int) bool
	Record(melds meld.Melds, tiles compact.Instances, totals compact.Totals)
}

type Calculator struct {
	tiles compact.Instances
	stack *meldStack

	results CalculateResults
	totals  compact.Totals
	sets    int
	opened  int
	minuses int

	baseMelds meld.BaseMelds
}

func NewCalculator(startMelds meld.BaseMelds, tiles compact.Instances, used compact.Instances, opened int) *Calculator {
	this := &Calculator{}
	this.tiles = tiles
	this.baseMelds = startMelds
	this.opened = opened
	this.stack = newMeldStack(7)
	totals := compact.NewTotals().Merge(used).Merge(tiles)
	this.totals = totals

	return this
}

func (this *Calculator) ResetResult(result CalculateResults) {
	this.results = result
}

func (this *Calculator) record() {
	this.results.Record(this.stack.Melds(), this.tiles, this.totals)
}

func (this *Calculator) Calculate() {
	parts := this.baseMelds.Filter(this.tiles, this.totals.FreeTiles())
	this.stack.Reset()
	this.sets = this.opened
	this.CalculateImpl(parts)

	this.sets = this.opened - 1
	this.tiles.Each(func(mask compact.Mask) bool {
		if mask.Count() < 2 {
			return true
		}
		m := this.push(meld.NewPairFromMask(mask).Meld())
		if m != 0 {
			this.CalculateImpl(parts)
			this.pop(m)
		}
		return true
	})
}

func getMissing(m meld.Meld) int {
	if m != 0 && !m.IsComplete() {
		return 1
	}
	return 0
}

func (this *Calculator) push(m meld.Meld) meld.Meld {
	missing := getMissing(m)
	if missing > 0 && !this.results.CheckMinuses(missing+this.minuses) {
		return 0
	}
	fixed := m.Rebase(this.tiles)
	if fixed == 0 {
		return 0
	}

	this.sets++
	this.minuses += missing
	fixed.ExtractFrom(this.tiles)
	this.stack.Push(fixed)
	return fixed
}

func (this *Calculator) pop(m meld.Meld) {
	this.stack.Pop()
	this.sets--
	this.minuses -= getMissing(m)
	m.AddTo(this.tiles)
}

func (this *Calculator) tryMeld(m meld.Meld, parts meld.Melds) bool {
	m = this.push(m)
	if m == 0 {
		return false
	}
	w := m.Waits()
	if w.IsEmpty() {
		this.CalculateImpl(parts)
	} else {
		base := m.Base()
		w.EachRange(base, base+3, func(t tile.Tile) bool {
			if this.totals.IsFull(t) {
				return true
			}
			this.totals.Add(t, 1)
			this.CalculateImpl(parts)
			this.totals.Add(t, -1)
			return true
		})
	}
	this.pop(m)
	return true
}

func (this *Calculator) CalculateImpl(parts meld.Melds) {
	if this.sets > 3 {
		this.record()
		return
	}
	one := false

	for k, meld := range parts {
		// Do not change order - must be calculated
		one = this.tryMeld(meld, parts[k:]) || one
	}
	if !one {
		this.record()
	}
}
