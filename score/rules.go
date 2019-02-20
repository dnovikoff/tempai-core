package score

import (
	"github.com/dnovikoff/tempai-core/yaku"
)

type Money int

type Rules interface {
	ManganRound() bool
	KazoeYakuman() bool
	IsDoubleYakuman(yaku.Yakuman) bool
	YakumanSum() bool

	Honba() Money
}

type RulesStruct struct {
	IsManganRound  bool
	IsKazoeYakuman bool
	DoubleYakumans map[yaku.Yakuman]bool
	IsYakumanSum   bool

	HonbaValue Money
}

var _ Rules = &RulesStruct{}

func (r *RulesStruct) ManganRound() bool {
	return r.IsManganRound
}

func (r *RulesStruct) KazoeYakuman() bool {
	return r.IsKazoeYakuman
}

func (r *RulesStruct) IsDoubleYakuman(y yaku.Yakuman) bool {
	return r.DoubleYakumans[y]
}

func (r *RulesStruct) YakumanSum() bool {
	return r.IsYakumanSum
}

func (r *RulesStruct) Honba() Money {
	return r.HonbaValue
}

func GetHonbaMoney(r Rules, honba Honba) Money {
	return Money(honba) * r.Honba()
}
