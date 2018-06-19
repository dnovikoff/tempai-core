package tempai

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/calc"
)

//go:generate stringer -type=Type
type Type int

const (
	TypeNone Type = iota
	TypeRegular
	TypePairs
	TypeKokushi
)

type TempaiResult struct {
	Type  Type
	Pair  calc.Meld
	Melds calc.Melds
	Last  calc.Meld
	Waits compact.Tiles
}

type TempaiResults struct {
	Hand     compact.Instances
	Declared calc.Melds
	Results  []*TempaiResult
}

func GetWaits(in *TempaiResults) compact.Tiles {
	if in == nil {
		return 0
	}
	c := compact.Tiles(0)
	for _, v := range in.Results {
		c |= v.Waits
	}
	return c
}

type result struct {
	original  compact.Instances
	tmp       compact.Instances
	results   []*TempaiResult
	validator calc.Validator
}

func (r *result) Record(in *calc.ResultData) {
	switch in.Sets {
	case 3:
		if in.Pair == nil {
			return
		}
	case 4:
	default:
		return
	}
	var last calc.Meld
	i := in.Left
	t1 := i[0]
	if in.Pair == nil {
		last = calc.Tanki(t1)
	} else {
		t2 := i[1]
		if t1 == t2 {
			last = calc.PonPart(t1)
		} else {
			diff := t2 - t1
			switch diff {
			case 1:
				last = calc.ChiPart1(t1)
			case 2:
				last = calc.ChiPart2(t1)
			}
		}
	}
	if last == nil {
		return
	}
	if !in.Validate(last) {
		return
	}
	waits := compact.Tiles(0)
	for _, v := range last.Waits() {
		if !in.Validator.Empty(v) {
			waits |= compact.FromTile(v)
		}
	}
	if waits.IsEmpty() {
		return
	}
	r.results = append(r.results, &TempaiResult{
		Type:  TypeRegular,
		Melds: in.Closed.Clone(),
		Pair:  in.Pair,
		Last:  last,
		Waits: waits,
	})
}
