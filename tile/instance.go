package tile

import (
	"sort"
)

// Instance is a representation of one of 136 tiles in mahjong, including their copies.
// Instance values starts from value of `1` and ends with value of `136`.
// Value `0` (InstanceNull) used for `undefined` (not set).
// Four sequenced numbers indicates same tile. Ex. 1,2,3,4 is for 1 Man, 5,6,7,8 is for 2 Man, etc.
// Tiles order is the following:
// - `123456789` `Man` (numbers 1-36)
// - `123456789` `Pin` (numbers 37-72)
// - `123456789` `Sou` (numbers 73-108)
// - `East`, `South`, `West`, `East` (numbers 109-124)
// - `White`, `Green`, `Red` (numbers 125-136).
// Take note, that for the last group values starts from `White` (not `Red`, as stated in some manuals)
// See also Tile.
type Instance int

const (
	// InstanceCount is a number of tiles. Results to 136.
	InstanceCount = TileCount * 4
	// InstanceBegin is a first valid tile.
	InstanceBegin Instance = 1
	// InstanceEnd is a tile after the last valid tile.
	// Use for iteration `for i:= InstanceBegin; i < InstanceEnd; i++`
	// or for validating valid ranges.
	InstanceEnd Instance = InstanceBegin + Instance(InstanceCount)
	// InstanceNull is a special value for indicating empty (no instance).
	InstanceNull Instance = 0
	NullCopy     CopyID   = -1
	AnyCopy      CopyID   = 0
)

// CopyID is a number of tile inside a group. Values are [0-3].
// For example there are 4 tiles of 1 Man.
type CopyID int

func newInstance(tile Tile, c CopyID) Instance {
	val := (tile-TileBegin)*4 + Tile(c)
	return Instance(val) + InstanceBegin
}

func (i Instance) Tile() Tile {
	return Tile(i-InstanceBegin)/4 + TileBegin
}

func (i Instance) CopyID() CopyID {
	return CopyID(i-InstanceBegin) % 4
}

type Instances []Instance

var _ sort.Interface = Instances{}

func (i Instances) Len() int               { return len(i) }
func (i Instances) Less(lhs, rhs int) bool { return i[lhs] < i[rhs] }
func (i Instances) Swap(lhs, rhs int)      { i[lhs], i[rhs] = i[rhs], i[lhs] }

func (i Instances) Tiles() Tiles {
	ret := make(Tiles, len(i))
	for k, v := range i {
		ret[k] = v.Tile()
	}
	return ret
}

// String is for debugging purposes and for using in tests.
func (i Instances) String() string {
	return i.Tiles().String()
}

// Clone clones the instances avoiding grow.
func (i Instances) Clone() Instances {
	x := make(Instances, len(i))
	for k, v := range i {
		x[k] = v
	}
	return x
}
