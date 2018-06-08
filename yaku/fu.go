package yaku

import (
	"github.com/dnovikoff/tempai-core/hand/calc"
)

func GetMeldFu(m calc.Meld) FuPoints {
	tags := m.Tags()
	if tags.CheckAny(calc.TagKanchan | calc.TagPenchan | calc.TagTanki) {
		return 2
	}
	if !tags.CheckAny(calc.TagPon) {
		return 0
	}
	var value FuPoints = 2
	if tags.CheckAny(calc.TagTerminal | calc.TagHonor) {
		value = 4
	}
	if !tags.CheckAny(calc.TagOpened) {
		value *= 2
	}
	if tags.CheckAny(calc.TagKan) {
		value *= 4
	}
	return value
}
