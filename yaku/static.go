package yaku

import (
	"bitbucket.org/dnovikoff/tempai-core/tile"
)

var RulesTenhouRed = Rules{
	OpenTanyao:           true,
	HaiteiIsFromLiveOnly: true,
	Renhou:               LimitNone,
	AkaDoras:             []tile.Instance{tile.Man5.Instance(0), tile.Pin5.Instance(0), tile.Sou5.Instance(0)},
	Ipatsu:               true,
	Ura:                  true,
}

var RulesEMA = Rules{
	OpenTanyao:           true,
	HaiteiIsFromLiveOnly: true,
	Renhou:               LimitMangan,
	Ipatsu:               true,
	Ura:                  true,
}

var RulesJPMPLA = Rules{
	OpenTanyao:           true,
	HaiteiIsFromLiveOnly: true,
	Renhou:               LimitMangan,
	Ipatsu:               false,
	Ura:                  false,
}

var RulesJPMPLB = RulesEMA
