package calc

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/tile"
)

type Counters []counterBlock

func (c Counters) CopyFrom(x Counters) {
	for k, v := range x {
		c[k] = v
	}
}

func NewCounters() Counters {
	return make(Counters, 4)
}

func (c Counters) Dec(t tile.Tile, x int) bool {
	b, i := split(t)
	return c[b].dec(i, x)
}

func (c Counters) Dec2(t tile.Tile, shift uint) bool {
	b, i := split(t)
	return c[b].dec2(i, shift)
}

func (c Counters) Dec3(t tile.Tile) bool {
	b, i := split(t)
	return c[b].dec3(i)
}

func (c Counters) Set(t tile.Tile, x int) {
	b, i := split(t)
	c[b].set(i, x)
}

func (c Counters) Get(t tile.Tile) int {
	b, i := split(t)
	return c[b].get(i)
}

func (c Counters) Tiles() compact.Tiles {
	var out compact.Tiles
	for bi, b := range c {
		n := 0
		for b > 0 {
			v := int(b & 7)
			if v != 0 {
				out = out.Set(tile.TileBegin + tile.Tile(bi*9+n))
			}
			n++
			b >>= 3
		}
	}
	return out
}

func (c Counters) Empty() bool {
	for _, v := range c {
		if v != 0 {
			return false
		}
	}
	return true
}

func (c Counters) write(buf tile.Tiles) tile.Tiles {
	index := 0
	for bi, b := range c {
		n := 0
		for b > 0 {
			v := int(b & 7)
			for y := 0; y < v; y++ {
				buf[index] = tile.TileBegin + tile.Tile(bi*9+n)
				index++
			}
			n++
			b >>= 3
		}
	}
	return buf[:index]
}

func (c Counters) Invert() {
	i := len(c)
	for k := range c[:i-1] {
		c[k].invert(9)
	}
	c[i-1].invert(7)
}

func split(t tile.Tile) (int, uint) {
	t -= tile.TileBegin
	return int(t / 9), uint(t % 9)
}

func CountersFromInstances(i compact.Instances) Counters {
	x := NewCounters()
	i.Each(func(mask compact.Mask) bool {
		x.Set(mask.Tile(), mask.Count())
		return true
	})
	return x
}

type Validator struct {
	c Counters
}

func (fv Validator) Empty(t tile.Tile) bool {
	return fv.c.Get(t) == 0
}

func (fv Validator) Validate(melds Melds) bool {
	for k, v := range melds {
		w := v.Waits()
		if w == nil {
			continue
		}
		for _, v := range w {
			cnt := fv.c.Get(v)
			if cnt == 0 {
				continue
			}
			fv.c.Set(v, cnt-1)
			result := fv.Validate(melds[k+1:])
			fv.c.Set(v, cnt)
			if result {
				return true
			}
		}
		return false
	}
	return true
}
