package tile

type Instance int

const (
	InstanceBegin Instance = 1
	InstanceNull  Instance = 0
	InstanceCount          = TileCount * 4
	NullCopy      CopyId   = -1
	AnyCopy       CopyId   = 0
)

type CopyId int

func NewInstance(tile Tile, c CopyId) Instance {
	val := (tile-Begin)*4 + Tile(c)
	return Instance(val) + InstanceBegin
}

func (this Instance) Tile() Tile {
	return Tile(this-InstanceBegin)/4 + Begin
}

func (this Instance) CopyId() CopyId {
	return CopyId(this-InstanceBegin) % 4
}

func (this Instance) Less(rhs Instance) bool {
	return this < rhs
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

func (this Instances) String() string {
	return this.Tiles().String()
}
