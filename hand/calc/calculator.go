package calc

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/tile"
)

type ResultData struct {
	Validator Validator
	Left      tile.Tiles
	Closed    Melds
	Pair      Meld
	Sets      int
}

type Results interface {
	Record(*ResultData)
}

type calculator struct {
	stack   *meldStack
	options *Options
	hand    Counters
	buf     tile.Tiles

	states []*state
	ResultData
}

func Calculate(startMelds Melds, tiles compact.Instances, opts *Options) {
	newCalculator(tiles, opts).run(startMelds)
}

func newCalculator(tiles compact.Instances, opts *Options) *calculator {
	c := &calculator{}
	c.hand = countersFromInstances(tiles)
	c.buf = make(tile.Tiles, 14)
	c.stack = newMeldStack(4)
	c.options = opts
	tmp := compact.NewInstances().Merge(c.options.Used).Merge(tiles)
	c.Validator.c = countersFromInstances(tmp)
	c.Validator.c.Invert()
	for _, v := range opts.Declared {
		if !v.Extract(c.Validator.c) {
			return nil
		}
	}
	c.states = newStates()
	return c
}

func (c *calculator) res() Results {
	return c.options.Results
}

func (c *calculator) opened() int {
	return c.options.Opened
}

func (c *ResultData) Validate(last Meld) bool {
	if last == nil {
		return c.Validator.Validate(c.Closed)
	}
	melds := make(Melds, 0, len(c.Closed)+1)
	melds = append(melds, c.Closed...)
	melds = append(melds, last)
	return c.Validator.Validate(melds)
}

func (c *calculator) record() {
	melds := c.stack.getMelds()
	c.Closed = melds
	c.Left = c.hand.write(c.buf)
	c.res().Record(&c.ResultData)
}

func (c *calculator) save() *state {
	x := c.states[c.Sets+1]
	x.save(c)
	return x
}

func (c *calculator) saveTop() *state {
	x := c.states[0]
	x.save(c)
	return x
}

func (c *calculator) run(parts Melds) {
	// parts := c.filterMelds(baseMelds) //.Filter(c.totals.FreeTiles())
	c.stack.reset()
	c.Sets = c.opened()
	c.subRun(parts)
	c.runPairs(parts)
}

func (c *calculator) runPairs(parts Melds) {
	state := c.saveTop()
	for i := tile.TileBegin; i < tile.TileEnd; i++ {
		if !c.hand.Dec(i, 2) {
			continue
		}
		c.Pair = Pair(i)
		c.subRun(parts)
		state.recover(c)
	}
}

func (c *calculator) push(m Meld) bool {
	if !m.Extract(c.hand) {
		return false
	}
	c.Sets++
	c.stack.push(m)
	return true
}

func (c *calculator) subRun(parts Melds) {
	cnt := c.Sets
	if c.Pair != nil {
		cnt++
	}
	if c.Sets > 4 {
		return
	} else if c.Sets == 4 {
		c.record()
		return
	}
	bestResult := true
	state := c.save()
	for k, m := range parts {
		if !c.push(m) {
			continue
		}
		bestResult = false
		idx := k
		if m.Tags().CheckAny(TagPon) {
			idx++
		}
		c.subRun(parts[idx:])
		state.recover(c)
		c.stack.pop()
	}
	if bestResult {
		c.record()
	}
}

func FilterMelds(tiles compact.Instances, m Melds) Melds {
	c := countersFromInstances(tiles)
	x := make(Melds, 0, len(m))
	tmp := NewCounters()
	tmp.copyFrom(c)
	for _, v := range m {
		if v.Extract(tmp) {
			x = append(x, v)
			tmp.copyFrom(c)
		}
	}
	return x
}

type state struct {
	data    Counters
	minuses int
	sets    int
}

func newStates() []*state {
	x := make([]*state, 5)
	for k := range x {
		x[k] = newState()
	}
	return x
}

func newState() *state {
	return &state{
		data: NewCounters(),
	}
}

func (s *state) save(c *calculator) {
	s.data.copyFrom(c.hand)
	s.sets = c.Sets
}

func (s *state) recover(c *calculator) {
	c.hand.copyFrom(s.data)
	c.Sets = s.sets
}
