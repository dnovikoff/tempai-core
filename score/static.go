package score

var RulesEMA = Rules{
	StartMoney:    0,
	ManganRound:   false,
	KazoeYakuman:  false,
	DoubleYakuman: false,
	YakumanSum:    false,
	Uma:           []Money{15000, 5000, -5000, -15000},
	Oka:           0,
	HonbaValue:    100,
}

var RulesTenhou = Rules{
	StartMoney:    25000,
	ManganRound:   false,
	KazoeYakuman:  true,
	DoubleYakuman: false,
	YakumanSum:    true,
	Uma:           []Money{20000, 10000, -10000, -20000},
	Oka:           20000,
	HonbaValue:    100,
}

var RulesJPMLA = Rules{
	StartMoney:    30000,
	ManganRound:   false,
	KazoeYakuman:  false,
	DoubleYakuman: false,
	YakumanSum:    false,
	// TODO: special UMA
	Uma:        []Money{15000, 5000, -5000, -15000},
	Oka:        0,
	HonbaValue: 100,
}

var RulesJPMLB = Rules{
	StartMoney:    30000,
	ManganRound:   true,
	KazoeYakuman:  false,
	DoubleYakuman: false,
	YakumanSum:    false,
	Uma:           []Money{15000, 5000, -5000, -15000},
	Oka:           0,
	HonbaValue:    100,
}
