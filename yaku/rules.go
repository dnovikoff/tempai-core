package yaku

import (
	"github.com/dnovikoff/tempai-core/tile"
)

type Rules interface {
	OpenTanyao() bool
	CheckAka(i tile.Instance) bool
	Renhou() Limit
	HaiteiFromLiveOnly() bool
	Ura() bool
	Ipatsu() bool
	GreenRequired() bool
	RinshanFu() bool
}

type RulesStruct struct {
	IsOpenTanyao         bool
	AkaDoras             []tile.Instance
	RenhouLimit          Limit
	IsHaiteiFromLiveOnly bool
	IsUra                bool
	IsIpatsu             bool
	IsGreenRequired      bool
	IsRinshanFu          bool
}

var _ Rules = &RulesStruct{}

func (r *RulesStruct) GreenRequired() bool {
	return r.IsGreenRequired
}

func (r *RulesStruct) RinshanFu() bool {
	return r.IsRinshanFu
}

func (r *RulesStruct) OpenTanyao() bool {
	return r.IsOpenTanyao
}

func (r *RulesStruct) Renhou() Limit {
	return r.RenhouLimit
}

func (r *RulesStruct) HaiteiFromLiveOnly() bool {
	return r.IsHaiteiFromLiveOnly
}

func (r *RulesStruct) Ura() bool {
	return r.IsUra
}

func (r *RulesStruct) Ipatsu() bool {
	return r.IsIpatsu
}

func (r *RulesStruct) CheckAka(i tile.Instance) bool {
	if len(r.AkaDoras) == 0 {
		return false
	}
	for _, v := range r.AkaDoras {
		if v == i {
			return true
		}
	}
	return false
}
