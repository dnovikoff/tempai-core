package yaku

import (
	"github.com/dnovikoff/tempai-core/tile"
)

func RulesTenhouRed() *RulesStruct {
	return &RulesStruct{
		IsOpenTanyao:         true,
		IsHaiteiFromLiveOnly: true,
		RenhouLimit:          LimitNone,
		AkaDoras:             []tile.Instance{tile.Man5.Instance(0), tile.Pin5.Instance(0), tile.Sou5.Instance(0)},
		IsIpatsu:             true,
		IsUra:                true,
		IsGreenRequired:      false,
		IsRinshanFu:          true,
	}
}

func RulesEMA() *RulesStruct {
	return &RulesStruct{
		IsOpenTanyao:         true,
		IsHaiteiFromLiveOnly: true,
		RenhouLimit:          LimitMangan,
		IsIpatsu:             true,
		IsUra:                true,
		IsGreenRequired:      false,
		IsRinshanFu:          true,
	}
}

func RulesJPMPLA() *RulesStruct {
	return &RulesStruct{
		IsOpenTanyao:         true,
		IsHaiteiFromLiveOnly: true,
		RenhouLimit:          LimitMangan,
		IsIpatsu:             false,
		IsUra:                false,
		IsGreenRequired:      true,
		IsRinshanFu:          false,
	}
}

func RulesJPMPLB() *RulesStruct {
	return RulesEMA()
}
