package tile

type Tile int
type Type int

const (
	Man1 Tile = iota
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

	End
)

const (
	TypeMan     Type = Type(Man1)
	TypePin     Type = Type(Pin1)
	TypeSou     Type = Type(Sou1)
	TypeWind    Type = Type(East)
	TypeDragon  Type = Type(White)
	TypeEnd     Type = Type(End)
	Count            = int(End) * 4
	SequenceEnd      = East
	Begin            = Man1
)

func createKokushTiles() Tiles {
	result := make(Tiles, 0, 13)
	for t := Begin; t < End; t++ {
		if !t.IsTerminalOrHonor() {
			continue
		}
		result = append(result, t)
	}
	return result
}

var KokushiTiles = createKokushTiles()

func (this Type) Rune() rune {
	switch this {
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

func (this Type) TileNumber(num int) Tile {
	return this.Tile(num - 1)
}

func (this Type) Tile(shift int) Tile {
	return Tile(this) + Tile(shift)
}

func NewTile(t Type, number int) Tile {
	return Tile(t) + Tile(number) - 1
}

func (this Tile) IsNull() bool {
	return this == End
}

func (this Tile) Type() Type {
	t := Type(this)
	switch {
	case this < Begin:
		return TypeEnd
	case t < TypePin:
		return TypeMan
	case t < TypeSou:
		return TypePin
	case t < TypeWind:
		return TypeSou
	case t < TypeDragon:
		return TypeWind
	case t < TypeEnd:
		return TypeDragon
	}
	return TypeEnd
}

func (this Tile) NumberInSequence() int {
	return int(this) - int(this.Type()) + 1
}

func (this Tile) IsHonor() bool {
	return Type(this) >= TypeWind
}

func (this Tile) IsSequence() bool {
	return !this.IsHonor()
}

func (this Tile) IsTerminal() bool {
	return this.IsSequence() && (this.NumberInSequence() == 1 || this.NumberInSequence() == 9)
}

func (this Tile) IsTerminalOrHonor() bool {
	return this.IsHonor() || this.IsTerminal()
}

func (this Tile) IsMiddle() bool {
	return this.IsSequence() && !this.IsTerminal()
}

// For all green yakuman
func (this Tile) IsGreen() bool {
	switch this {
	case Green, Sou2, Sou3, Sou4, Sou6, Sou8:
		return true
	}
	return false
}

func (this Tile) Indicates() Tile {
	next := this + 1
	if this.Type() != next.Type() {
		return Tile(this.Type())
	}
	return next
}

func (this Tile) DistanceAbs(rhs Tile) (ret int) {
	ret = int(this - rhs)
	if ret < 0 {
		ret *= -1
	}
	return
}

func (this Tile) StringOrNull() string {
	if this.IsNull() {
		return "!"
	}
	return this.String()
}

func (this Tile) String() string {
	return Tiles{this}.String()
}

func (this Tile) Instance(c CopyId) Instance {
	return NewInstance(this, c)
}

type Tiles []Tile

func (this Tiles) CheckLinear(t Tile) bool {
	for _, v := range this {
		if t == v {
			return true
		}
	}
	return false
}

func (this Tiles) Clone() Tiles {
	x := make(Tiles, len(this))
	for k, v := range this {
		x[k] = v
	}
	return x
}

func (this Tiles) String() (ret string) {
	return TilesToTenhouString(this)
}

func (this Tile) VisitImproves(f func(Tile)) {
	if !this.IsSequence() {
		f(this)
		return
	}
	for x := this - 2; x < this+3; x++ {
		if this.Type() == x.Type() {
			f(x)
		}
	}

}

func (a Tiles) Len() int           { return len(a) }
func (a Tiles) Less(i, j int) bool { return a[i] < a[j] }
func (a Tiles) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
