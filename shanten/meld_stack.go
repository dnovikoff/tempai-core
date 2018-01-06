package shanten

import "bitbucket.org/dnovikoff/tempai-core/meld"

type meldStack struct {
	melds meld.Melds
	index int
}

func newMeldStack(capacity int) *meldStack {
	this := &meldStack{}
	this.melds = make(meld.Melds, capacity)
	return this
}

func (this *meldStack) Size() int {
	return this.index
}

func (this *meldStack) Reset() {
	this.index = 0
}

func (this *meldStack) Push(meld meld.Meld) {
	this.melds[this.index] = meld
	this.index++
}

func (this *meldStack) Melds() meld.Melds {
	return this.melds[:this.index]
}

func (this *meldStack) Back() meld.Meld {
	return this.melds[this.index-1]
}

func (this *meldStack) Pop() {
	this.index--
}
