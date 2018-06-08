package calc

import (
	"sort"

	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/tile"
)

type kokushi struct {
	meld
	hole tile.Tile
}

var kokushi13waits = compact.TerminalOrHonor.Tiles()

var kokushi13meld = &kokushi{meld{
	tags:   TagKoksuhi13 | TagKokushi,
	waits:  kokushi13waits,
	cwaits: compact.TerminalOrHonor,
	tiles:  kokushi13waits,
}, tile.TileNull}

func KokushiComplete(p tile.Tile) Meld {
	hand := kokushi13waits.Clone()
	hand = append(hand, p)
	sort.Sort(hand)
	return &kokushi{meld{
		t:     p,
		tags:  TagKokushi | TagComplete,
		tiles: hand,
	}, tile.TileNull}
}

func Kokushi13() Meld {
	return kokushi13meld
}

func KokushiMeld(pair tile.Tile, wait tile.Tile) Meld {
	x := kokushi13waits.Clone()
	for k, v := range x {
		if v == wait {
			x[k] = pair
		}
	}
	sort.Sort(x)
	return &kokushi{meld{
		tags:   TagKokushi,
		waits:  tile.Tiles{wait},
		cwaits: compact.FromTile(wait),
		t:      pair,
		tiles:  x,
	}, wait}
}

func (*kokushi) Extract(Counters) bool {
	panic("kokushi.Extract not implemented")
}

func (k *kokushi) Complete(t tile.Tile) Meld {
	if !k.cwaits.Check(t) {
		return nil
	}
	if k.hole == t {
		return KokushiComplete(k.t)
	}
	return KokushiComplete(t)
}
