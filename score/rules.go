package score

type Money int

type Rules struct {
	ManganRound   bool  `json:"mangan-round"`
	KazoeYakuman  bool  `json:"kazoe-yakuman"`
	DoubleYakuman bool  `json:"double-yakuman"`
	YakumanSum    bool  `json:"yakuman-sum"`
	HonbaValue    Money `json:"honba-value"`
}

func (this Rules) GetHonbaMoney(honba Honba) Money {
	return Money(honba) * this.HonbaValue
}
