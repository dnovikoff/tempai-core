package calc

import (
	"github.com/dnovikoff/tempai-core/tile"
)

type Melds []Meld

func (m Melds) Clone() Melds {
	x := make(Melds, len(m))
	for k, v := range m {
		x[k] = v
	}
	return x
}

func createAllPon() Melds {
	x := make(Melds, 0, tile.TileCount)
	for i := tile.TileBegin; i < tile.TileEnd; i++ {
		x = append(x, Pon(i))
	}
	return x
}

func createAllPonParts() Melds {
	x := make(Melds, 0, tile.TileCount)
	for i := tile.TileBegin; i < tile.TileEnd; i++ {
		x = append(x, PonPart(i))
	}
	return x
}

func createAllChi() Melds {
	x := make(Melds, 0, int(tile.SequenceEnd-tile.SequenceBegin))
	for i := tile.SequenceBegin; i < tile.SequenceEnd; i++ {
		c := Chi(i)
		if c == nil {
			continue
		}
		x = append(x, c)
	}
	return x
}

func createAllChiParts() Melds {
	x := make(Melds, 0, int(tile.SequenceEnd-tile.SequenceBegin)*3)
	for i := tile.SequenceBegin; i < tile.SequenceEnd; i++ {
		m := ChiPart1(i)
		if m != nil {
			x = append(x, m)
		}
		m = ChiPart2(i)
		if m != nil {
			x = append(x, m)
		}
	}
	return x
}

func CreateComplete() Melds {
	return append(createAllChi(),
		createAllPon()...)
}

func CreateParts() Melds {
	return append(
		createAllChiParts(),
		createAllPonParts()...)
}

func CreateAll() Melds {
	return append(
		CreateParts(),
		CreateComplete()...)
}
