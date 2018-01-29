package effective

import "github.com/dnovikoff/tempai-core/tile"

type sSortId uint32

func newSortId(uCount, tCount, value int) sSortId {
	val := (13 - sSortId(value))
	val = (val << 8) | sSortId(uCount)
	val = (val << 8) | sSortId(tCount)

	return val
}

func (this sSortId) betterThan(other sSortId) bool {
	return this > other
}

type ResultsSorted []*Result

func (this ResultsSorted) Best() *Result {
	if len(this) == 0 {
		return nil
	}
	return this[0]
}

func (this ResultsSorted) Len() int {
	return len(this)
}

func (this ResultsSorted) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func tileSortId(t tile.Tile) int {
	return (tilePriority(t) << 8) | int(t)
}

func tileLess(l, r tile.Tile) bool {
	return tileSortId(l) < tileSortId(r)
}

func effLess(l, r *Result) bool {
	lhs, rhs := l.sortId, r.sortId
	if lhs.betterThan(rhs) {
		return true
	} else if rhs.betterThan(lhs) {
		return false
	}
	return tileLess(l.Tile, r.Tile)
}

func (this ResultsSorted) Less(i, j int) bool {
	return effLess(this[i], this[j])
}

func (this ResultsSorted) First() (ret *Result) {
	if len(this) == 0 {
		return
	}
	return this[0]
}
