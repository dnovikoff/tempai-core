package effective

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
