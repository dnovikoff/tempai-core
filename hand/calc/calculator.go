package calc

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/meld"
	"github.com/dnovikoff/tempai-core/tile"
)

type Results interface {
	CheckMinuses(minuses int) bool
	Record(melds meld.Melds, tiles compact.Instances, totals compact.Totals)
}

type calculator struct {
	tiles   compact.Instances
	stack   *meldStack
	options *Options

	totals  compact.Totals
	sets    int
	minuses int

	baseMelds meld.BaseMelds
}

func Calculate(startMelds meld.BaseMelds, tiles compact.Instances, opts *Options) {
	newCalculator(startMelds, tiles, opts).run()
}

func newCalculator(startMelds meld.BaseMelds, tiles compact.Instances, opts *Options) *calculator {
	c := &calculator{}
	c.tiles = tiles
	c.baseMelds = startMelds
	c.stack = newMeldStack(7)
	c.options = opts
	totals := compact.NewTotals().Merge(c.options.Used).Merge(tiles)
	c.totals = totals

	return c
}

func (c *calculator) res() Results {
	return c.options.Results
}

func (c *calculator) opened() int {
	return c.options.Opened
}

func (c *calculator) record() {
	c.res().Record(c.stack.getMelds(), c.tiles, c.totals)
}

func (c *calculator) run() {
	parts := c.baseMelds.Filter(c.tiles, c.totals.FreeTiles())
	c.stack.reset()
	c.sets = c.opened()
	c.calculateImpl(parts)

	c.sets = c.opened() - 1
	c.tiles.Each(func(mask compact.Mask) bool {
		if mask.Count() < 2 {
			return true
		}
		m := c.push(meld.NewPairFromMask(mask).Meld())
		if m != 0 {
			c.calculateImpl(parts)
			c.pop(m)
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

func (c *calculator) push(m meld.Meld) meld.Meld {
	missing := getMissing(m)
	if missing > 0 && !c.res().CheckMinuses(missing+c.minuses) {
		return 0
	}
	fixed := m.Rebase(c.tiles)
	if fixed == 0 {
		return 0
	}

	c.sets++
	c.minuses += missing
	fixed.ExtractFrom(c.tiles)
	c.stack.push(fixed)
	return fixed
}

func (c *calculator) pop(m meld.Meld) {
	c.stack.pop()
	c.sets--
	c.minuses -= getMissing(m)
	m.AddTo(c.tiles)
}

func (c *calculator) tryMeld(m meld.Meld, parts meld.Melds) bool {
	m = c.push(m)
	if m == 0 {
		return false
	}
	w := m.Waits()
	if w.IsEmpty() {
		c.calculateImpl(parts)
	} else {
		base := m.Base()
		w.EachRange(base, base+3, func(t tile.Tile) bool {
			if c.totals.IsFull(t) {
				return true
			}
			c.totals.Add(t, 1)
			c.calculateImpl(parts)
			c.totals.Add(t, -1)
			return true
		})
	}
	c.pop(m)
	return true
}

func (c *calculator) calculateImpl(parts meld.Melds) {
	if c.sets > 3 {
		c.record()
		return
	}
	one := false

	for k, meld := range parts {
		// Do not change order - must be calculated
		one = c.tryMeld(meld, parts[k:]) || one
	}
	if !one {
		c.record()
	}
}
