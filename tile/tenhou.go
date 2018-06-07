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
	if len(str) == 0 {
		return nil, nil
	}
	tmp := make(Tiles, 0, len(str))
	index := 0
	t := TileEnd
	max := '0'
	for k, v := range str {
		r := rune(v)
		switch r {
		case 's':
			t = Sou1
		case 'm':
			t = Man1
		case 'p':
			t = Pin1
		case 'z':
			t = East
			if max > '7' {
				return nil, stackerr.Newf("Unexpected value '%s' for type '%s'", max, v)
			}
		default:
			if r < '1' || r > '9' {
				return nil, stackerr.Newf("Unexpected symbol '%s' at position %v", v, k)
			}
			if r > max {
				max = r
			}
		}
		if t != TileEnd {
			if index == k {
				return nil, stackerr.Newf("Empty range at %v", index)
			}
			for _, val := range str[index:k] {
				tmp = append(tmp, Tile(int(rune(val)-'1'))+t)
			}
			index = k + 1
			max = '0'
			t = TileEnd
		}
	}
	if index != len(str) {
		return nil, stackerr.Newf("Expected to end with a letter")
	}
	return tmp, nil
}
