package score

var RulesEMA = Rules{
	ManganRound:   false,
	KazoeYakuman:  false,
	DoubleYakuman: false,
	YakumanSum:    false,
	HonbaValue:    100,
}

var RulesTenhou = Rules{
	ManganRound:   false,
	KazoeYakuman:  true,
	DoubleYakuman: false,
	YakumanSum:    true,
	HonbaValue:    100,
}

var RulesJPMLA = Rules{
	ManganRound:   false,
	KazoeYakuman:  false,
	DoubleYakuman: false,
	YakumanSum:    false,
	HonbaValue:    100,
}

var RulesJPMLB = Rules{
	ManganRound:   true,
	KazoeYakuman:  false,
	DoubleYakuman: false,
	YakumanSum:    false,
	HonbaValue:    100,
}
