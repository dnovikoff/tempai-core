package tile

import (
	"fmt"

	"github.com/facebookgo/stackerr"
)

func TenhoTypeToString(t Type) (ret string) {
	if t >= TypeWind {
		return "z"
	}
	return string(TypeRune(t))
}

func TilesToTenhouString(tiles Tiles) string {
	if len(tiles) == 0 {
		return ""
	}
	prevT := tiles[0].Type()
	if prevT == TypeDragon {
		prevT = TypeWind
	}
	ret := ""
	for _, v := range tiles {
		t := v.Type()
		n := v.Number()
		if t == TypeDragon {
			t = TypeWind
			n += 4
		}
		if t != prevT {
			ret += TenhoTypeToString(prevT)
			prevT = t
		}

		ret += fmt.Sprintf("%v", n)
	}
	ret += TenhoTypeToString(prevT)
	return ret
}

func NewTilesFromString(str string) (Tiles, error) {
	p := &parser{}
	p.Parse(str)
	if p.err != nil {
		return nil, p.err
	}
	return p.result, nil
}

type parser struct {
	index  int
	max    rune
	err    error
	input  string
	result []Tile
}

func (p *parser) fill(k int, t Tile) {
	if p.index == k {
		p.err = stackerr.Newf("Empty range at %v", p.index)
		return
	}
	for _, val := range p.input[p.index:k] {
		p.result = append(p.result, Tile(int(rune(val)-'1'))+t)
	}
	p.index = k + 1
	p.max = '0'
}

func (p *parser) Parse(str string) {
	if len(str) == 0 {
		return
	}
	p.input = str
	p.result = make(Tiles, 0, len(str))
	p.index = 0
	p.max = '0'
	p.err = nil
	for k, v := range str {
		r := rune(v)
		switch r {
		case 's':
			p.fill(k, Sou1)
		case 'm':
			p.fill(k, Man1)
		case 'p':
			p.fill(k, Pin1)
		case 'z':
			if p.max > '7' {
				p.err = stackerr.Newf("Unexpected value '%s' for type '%s'", string(p.max), string(v))
			} else {
				p.fill(k, East)
			}
		default:
			if r < '1' || r > '9' {
				p.err = stackerr.Newf("Unexpected symbol '%s' at position %v", string(v), k)
			} else if r > p.max {
				p.max = r
			}
		}
		if p.err != nil {
			return
		}
	}
	if p.index != len(str) {
		p.err = stackerr.Newf("Expected to end with a letter")
	}
	return
}
