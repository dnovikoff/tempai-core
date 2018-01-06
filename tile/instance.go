package tile

type Instance int

type CopyId int

const NullCopy CopyId = -1

func NewInstanceFromInt(x int) Instance {
	return Instance(x)
}

func (this Instance) IntValue() int {
	return int(this)
}

func NewInstance(tile Tile, c CopyId) Instance {
	val := int(tile*4) + int(c)
	return NewInstanceFromInt(val)
}

func (this Instance) Tile() Tile {
	return Tile(this.IntValue() / 4)
}

func (this Instance) CopyId() CopyId {
	return CopyId(this.IntValue() % 4)
}

func (this Instance) Next() Instance {
	val := (this.IntValue() + 4) % int(Count*4)
	return NewInstanceFromInt(val)
}

func (this Instance) Prev() Instance {
	val := this.IntValue() - 4
	if val < 0 {
		val += int(End * 4)
	}
	return NewInstanceFromInt(val)
}

func (this Instance) Less(rhs Instance) bool {
	return this.IntValue() < rhs.IntValue()
}

const (
	InstanceNull = -1
)

func (this Instance) IsNull() bool {
	return this.IntValue() == -1
}

type Instances []Instance

func (a Instances) Len() int           { return len(a) }
func (a Instances) Less(i, j int) bool { return a[i].Less(a[j]) }
func (a Instances) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (this Instances) Tiles() (ret Tiles) {
	ret = make(Tiles, len(this))
	for k, v := range this {
		ret[k] = v.Tile()
	}
	return
}

func (this Instance) Advance(x int) Instance {
	return NewInstanceFromInt(this.IntValue() + x*4)
}

func (this Instance) StringOrNull() string {
	if this.IsNull() {
		return "!"
	}
	return this.Tile().String()
}

func (this Instances) String() string {
	return this.Tiles().String()
}
