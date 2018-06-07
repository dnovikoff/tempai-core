package yaku

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dnovikoff/tempai-core/meld"
)

type YakuSet map[Yaku]HanPoints
type YakumanSet map[Yakuman]int

type Result struct {
	Melds   meld.Melds
	Yaku    YakuSet
	Yakuman YakumanSet
	Bonuses YakuSet
	Fus     Fus

	IsClosed bool
}

func newResult(melds meld.Melds) *Result {
	res := &Result{}
	res.Melds = melds
	res.Yaku = make(YakuSet, 16)
	res.Bonuses = make(YakuSet, 16)
	res.Yakuman = make(YakumanSet, 16)
	return res
}

func (r *Result) Sum() HanPoints {
	x := r.Yaku.Sum()
	if x == 0 {
		return 0
	}
	return x + r.Bonuses.Sum()
}

func (r *Result) String() string {
	if len(r.Yakuman) > 0 {
		return r.Yakuman.String()
	}
	if len(r.Yaku) > 0 {
		x := r.Yaku.String()
		if len(r.Bonuses) > 0 {
			x += ", " + r.Bonuses.String()
		}
		return fmt.Sprintf("%v = %v", r.Sum(), x)
	}
	return "No yaku"
}

func (r *Result) setValues(k Yaku, opened, closed HanPoints) {
	if r.Yaku[k] == 0 {
		return
	}

	if r.IsClosed {
		r.Yaku[k] = closed
	} else {
		r.Yaku[k] = opened
	}
	if r.Yaku[k] == 0 {
		delete(r.Yaku, k)
	}
}

type FuInfo struct {
	Meld   meld.Meld
	Fu     Fu
	Points FuPoints
}

func (y YakuSet) Sum() HanPoints {
	sum := HanPoints(0)
	for _, v := range y {
		sum += HanPoints(v)
	}
	return sum
}

func (y YakuSet) String() string {
	results := make([]string, 0, len(y))
	for k, v := range y {
		results = append(results, fmt.Sprintf("%v: %v", k, v))
	}
	sort.Strings(results)
	return strings.Join(results, ", ")
}

func (ys YakumanSet) String() string {
	results := make([]string, 0, len(ys))
	for k, v := range ys {
		str := fmt.Sprintf("%v:%v", k.String(), v)
		results = append(results, str)
	}
	sort.Strings(results)
	return strings.Join(results, ", ")
}

type Fus []*FuInfo

func (f Fus) Sum() FuPoints {
	var res FuPoints
	for _, v := range f {
		res += v.Points
	}
	return res
}

func (f Fus) String() string {
	parts := make([]string, 0, len(f))
	for _, v := range f {
		part := fmt.Sprintf("%v(%v)", v.Points, v.Fu)
		if !v.Meld.IsNull() {
			part += "[" + meld.DebugDescribe(v.Meld) + "]"
		}
		parts = append(parts, part)
	}
	return fmt.Sprintf("%v = %v", f.Sum(), strings.Join(parts, " + "))
}
