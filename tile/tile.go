package tile

// Tile is a representation of an uniqe tile. THere are 34 unique tiles.
// Tile values starts from value of `1` and ends with value of `34`.
// Value `0` (TileNull) used for `undefined` (not set).
// Tiles order is the following:
// - `123456789` `Man` (numbers 1-9)
// - `123456789` `Pin` (numbers 10-18)
// - `123456789` `Sou` (numbers 19-27)
// - `East`, `South`, `West`, `East` (numbers 28-31)
// - `White`, `Green`, `Red` (numbers 32-34)
// Take note, that for the last group values starts from `White` (not `Red`, as stated in some manuals)
// See also Instance.
type Tile int

// Tile numberation starts from 1
const (
	TileNull Tile = iota
	Man1
	Man2
	Man3
	Man4
	Man5
	Man6
	Man7
	Man8
	Man9

	Pin1
	Pin2
	Pin3
	Pin4
	Pin5
	Pin6
	Pin7
	Pin8
	Pin9

	Sou1
	Sou2
	Sou3
	Sou4
	Sou5
	Sou6
	Sou7
	Sou8
	Sou9

	East
	South
	West
	North

	White
	Green
	Red

	TileEnd
)

const (
	TileCount     = int(TileEnd - TileBegin)
	SequenceBegin = TileBegin
	SequenceEnd   = East
	TileBegin     = Man1
)

func (t Tile) Type() Type {
	tp := Type(t)
	switch {
	case t < TileBegin:
		return TypeNull
	case tp < TypePin:
		return TypeMan
	case tp < TypeSou:
		return TypePin
	case tp < TypeWind:
		return TypeSou
	case tp < TypeDragon:
		return TypeWind
	case t < TileEnd:
		return TypeDragon
	}
	return TypeNull
}

func (t Tile) Number() int {
	return int(t) - int(t.Type()) + 1
}

// Indicates used for dora indicators to choose dora tile
func (t Tile) Indicates() Tile {
	next := t + 1
	if t.Type() != next.Type() {
		return Tile(t.Type())
	}
	return next
}

func (t Tile) String() string {
	return Tiles{t}.String()
}

func (t Tile) Instance(c CopyID) Instance {
	return newInstance(t, c)
}

type Tiles []Tile

func (t Tiles) Contains(x Tile) bool {
	for _, v := range t {
		if x == v {
			return true
		}
	}
	return false
}

func (t Tiles) Clone() Tiles {
	x := make(Tiles, len(t))
	for k, v := range t {
		x[k] = v
	}
	return x
}

func (t Tiles) String() string {
	return TilesToTenhouString(t)
}

func (t Tiles) Len() int           { return len(t) }
func (t Tiles) Less(i, j int) bool { return t[i] < t[j] }
func (t Tiles) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
