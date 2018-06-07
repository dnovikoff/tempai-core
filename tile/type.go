package tile

type Type int

const (
	TypeMan    Type = Type(Man1)
	TypePin    Type = Type(Pin1)
	TypeSou    Type = Type(Sou1)
	TypeWind   Type = Type(East)
	TypeDragon Type = Type(White)
	TypeNull   Type = Type(TileNull)
)

func (t Type) Tile(num int) Tile {
	return Tile(t) + Tile(num-1)
}

// TypeRune used for stringifying tiles
func TypeRune(t Type) rune {
	switch t {
	case TypeMan:
		return 'm'
	case TypePin:
		return 'p'
	case TypeSou:
		return 's'
	case TypeDragon, TypeWind:
		return 'z'
	}
	return '-'
}
