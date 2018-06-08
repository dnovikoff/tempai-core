package calc

import "github.com/dnovikoff/tempai-core/meld"

type meldStack struct {
	melds meld.Melds
	index int
}

func newMeldStack(capacity int) *meldStack {
	return &meldStack{
		melds: make(meld.Melds, capacity),
	}
}

func (ms *meldStack) size() int {
	return ms.index
}

func (ms *meldStack) reset() {
	ms.index = 0
}

func (ms *meldStack) push(meld meld.Meld) {
	ms.melds[ms.index] = meld
	ms.index++
}

func (ms *meldStack) getMelds() meld.Melds {
	return ms.melds[:ms.index]
}

func (ms *meldStack) back() meld.Meld {
	return ms.melds[ms.index-1]
}

func (ms *meldStack) pop() {
	ms.index--
}
