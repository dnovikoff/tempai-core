package yaku

import (
	"bitbucket.org/dnovikoff/tempai-core/tile"
)

type Rules struct {
	OpenTanyao           bool            `json:"open-tanyao"`
	AkaDoras             []tile.Instance `json:"akas"`
	Renhou               Limit           `json:"renhou"`
	HaiteiIsFromLiveOnly bool            `json:"haitei-is-from-live-only"`
	Ura                  bool            `json:"ura"`
	Ipatsu               bool            `json:"ipatsu"`
}

func (this Rules) CheckAka(i tile.Instance) bool {
	if len(this.AkaDoras) == 0 {
		return false
	}
	for _, v := range this.AkaDoras {
		if v == i {
			return true
		}
	}
	return false
}
