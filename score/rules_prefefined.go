package score

import (
	"github.com/dnovikoff/tempai-core/yaku"
)

func RulesEMA() *RulesStruct {
	return &RulesStruct{
		IsManganRound:  false,
		IsKazoeYakuman: false,
		IsYakumanSum:   false,
		HonbaValue:     100,
	}
}

func RulesTenhou() *RulesStruct {
	return &RulesStruct{
		IsManganRound:  false,
		IsKazoeYakuman: true,
		IsYakumanSum:   true,
		HonbaValue:     100,
	}
}

func RulesJPMLA() *RulesStruct {
	return &RulesStruct{
		IsManganRound:  false,
		IsKazoeYakuman: false,
		IsYakumanSum:   false,
		HonbaValue:     100,
	}
}

func RulesJPMLB() *RulesStruct {
	return &RulesStruct{
		IsManganRound:  true,
		IsKazoeYakuman: false,
		IsYakumanSum:   false,
		HonbaValue:     100,
	}
}

func DefaultDoubleYakumans() map[yaku.Yakuman]bool {
	return map[yaku.Yakuman]bool{
		yaku.YakumanChuurenpooto9: true,
		yaku.YakumanKokushi13:     true,
		yaku.YakumanSuuankouTanki: true,
		yaku.YakumanDaisuushi:     true,
	}
}
