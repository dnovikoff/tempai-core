package calc

import "github.com/dnovikoff/tempai-core/meld"

type meldStack struct {
	melds meld.Melds
	index int
}

func newMeldStack(capacity int) *meldStack {
	this := &meldStack{}
	this.melds = make(meld.Melds, capacity)
	return this
}

func (this *meldStack) size() int {
	return this.index
}

func (this *meldStack) reset() {
	this.index = 0
}

func (this *meldStack) push(meld meld.Meld) {
	this.melds[this.index] = meld
	this.index++
}

func (this *meldStack) getMelds() meld.Melds {
	return this.melds[:this.index]
}

func (this *meldStack) back() meld.Meld {
	return this.melds[this.index-1]
}

func (this *meldStack) pop() {
	this.index--
}
