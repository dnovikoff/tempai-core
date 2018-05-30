package score

type Money int

type Rules interface {
	ManganRound() bool
	KazoeYakuman() bool
	YakumanDouble() bool
	YakumanSum() bool
	Honba() Money
}

type RulesStruct struct {
	IsManganRound   bool
	IsKazoeYakuman  bool
	IsYakumanDouble bool
	IsYakumanSum    bool
	HonbaValue      Money
}

var _ Rules = &RulesStruct{}

func (r *RulesStruct) ManganRound() bool {
	return r.IsManganRound
}

func (r *RulesStruct) KazoeYakuman() bool {
	return r.IsKazoeYakuman
}

func (r *RulesStruct) YakumanDouble() bool {
	return r.IsYakumanDouble
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
