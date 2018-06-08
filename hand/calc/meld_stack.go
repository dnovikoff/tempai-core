package calc

type meldStack struct {
	melds Melds
	index int
}

func newMeldStack(capacity int) *meldStack {
	return &meldStack{
		melds: make(Melds, capacity),
	}
}

func (ms *meldStack) size() int {
	return ms.index
}

func (ms *meldStack) reset() {
	ms.index = 0
}

func (ms *meldStack) push(meld Meld) {
	ms.melds[ms.index] = meld
	ms.index++
}

func (ms *meldStack) getMelds() Melds {
	return ms.melds[:ms.index]
}

func (ms *meldStack) back() Meld {
	return ms.melds[ms.index-1]
}

func (ms *meldStack) pop() {
	ms.index--
}
