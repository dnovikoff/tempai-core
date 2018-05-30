package score

func RulesEMA() *RulesStruct {
	return &RulesStruct{
		IsManganRound:   false,
		IsKazoeYakuman:  false,
		IsYakumanDouble: false,
		IsYakumanSum:    false,
		HonbaValue:      100,
	}
}

func RulesTenhou() *RulesStruct {
	return &RulesStruct{
		IsManganRound:   false,
		IsKazoeYakuman:  true,
		IsYakumanDouble: false,
		IsYakumanSum:    true,
		HonbaValue:      100,
	}
}

func RulesJPMLA() *RulesStruct {
	return &RulesStruct{
		IsManganRound:   false,
		IsKazoeYakuman:  false,
		IsYakumanDouble: false,
		IsYakumanSum:    false,
		HonbaValue:      100,
	}
}

func RulesJPMLB() *RulesStruct {
	return &RulesStruct{
		IsManganRound:   true,
		IsKazoeYakuman:  false,
		IsYakumanDouble: false,
		IsYakumanSum:    false,
		HonbaValue:      100,
	}
}
