package calc

type counterBlock int64

func (b *counterBlock) set(i uint, c int) {
	*b = (*b) & (^(counterBlock(7) << (i * 3)))
	*b = (*b) | counterBlock((c&7)<<(i*3))
}

func (b *counterBlock) get(i uint) int {
	return (int(*b) >> (i * 3)) & 7
}

func (b *counterBlock) invert() {
	for i := uint(0); i < 9; i++ {
		b.set(i, 4-b.get(i))
	}
}

func (b *counterBlock) dec(i uint, c int) bool {
	c = b.get(i) - c
	if c < 0 {
		return false
	}
	b.set(i, c)
	return true
}

func (b *counterBlock) dec2(i uint, shift uint) bool {
	x1 := b.get(i)
	x2 := b.get(i + shift)
	if x1 < 1 || x2 < 1 {
		return false
	}
	b.set(i, x1-1)
	b.set(i+shift, x2-1)
	return true
}

func (b *counterBlock) dec3(i uint) bool {
	x1 := b.get(i)
	x2 := b.get(i + 1)
	x3 := b.get(i + 2)
	if x1 < 1 || x2 < 1 || x3 < 1 {
		return false
	}
	b.set(i, x1-1)
	b.set(i+1, x2-1)
	b.set(i+2, x3-1)
	return true
}
