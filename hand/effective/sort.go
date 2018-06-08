package effective

import "github.com/dnovikoff/tempai-core/tile"

type sSortId uint32

func newSortId(uCount, tCount, value int) sSortId {
	val := (13 - sSortId(value))
	val = (val << 8) | sSortId(uCount)
	val = (val << 8) | sSortId(tCount)

	return val
}

func (id sSortId) betterThan(other sSortId) bool {
	return id > other
}

type ResultsSorted []*Result

func (r ResultsSorted) Best() *Result {
	if len(r) == 0 {
		return nil
	}
	return r[0]
}

func (r ResultsSorted) Len() int {
	return len(r)
}

func (r ResultsSorted) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
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

func (r ResultsSorted) Less(i, j int) bool {
	return effLess(r[i], r[j])
}

func (r ResultsSorted) First() *Result {
	if len(r) == 0 {
		return nil
	}
	return r[0]
}
