package tile

import (
	"fmt"

	"github.com/facebookgo/stackerr"
)

func TenhoTypeToString(t Type) (ret string) {
	if t >= TypeWind {
		return "z"
	}
	return string(t.Rune())
}

func TilesToTenhouString(tiles Tiles) (ret string) {
	if len(tiles) == 0 {
		return
	}
	prevT := tiles[0].Type()
	if prevT == TypeDragon {
		prevT = TypeWind
	}
	for _, v := range tiles {
		t := v.Type()
		n := v.NumberInSequence()
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
	return
}

func NewTilesFromString(str string) (ret Tiles, err error) {
	tmp := make(Tiles, 0, len(str))
	if len(str) == 0 {
		return Tiles{}, nil
	}
	index := 0
	t := End
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
				err = stackerr.Newf("Unexpected value %v for type %v", string(max), string(v))
				return
			}
		default:
			if r < '1' || r > '9' {
				err = stackerr.Newf("Unexpected symbol %v at position %v", string(v), k)
				return
			}
			if r > max {
				max = r
			}
		}
		if t != End {
			if index == k {
				err = stackerr.Newf("Empty range at %v", index)
				return
			}
			for _, val := range str[index:k] {
				tmp = append(tmp, Tile(int(rune(val)-'1'))+t)
			}
			index = k + 1
			max = '0'
			t = End
		}
	}

	if index != len(str) {
		err = stackerr.Newf("Expected to end with a letter")
		return
	}
	ret = tmp
	return
}
