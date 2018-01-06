package compact

import (
	"errors"

	"bitbucket.org/dnovikoff/tempai-core/tile"
)

type Generator struct {
	tiles Instances
}

func NewTileGenerator() *Generator {
	return &Generator{NewAllInstances()}
}

func (this *Generator) TilesLeft() tile.Instances {
	return this.tiles.Instances()
}

func (this *Generator) InstancePtr(t tile.Tile) (ret *tile.Instance) {
	x := this.Instance(t)
	return &x
}

func (this *Generator) Instance(t tile.Tile) (ret tile.Instance) {
	m := this.tiles.GetMask(t)
	if m.IsEmpty() {
		return tile.InstanceNull
	}
	ret = m.First()
	this.tiles.Remove(ret)
	return
}

func (this *Generator) CompactFromString(str string) (ret Instances, err error) {
	tiles, err := tile.NewTilesFromString(str)
	if err != nil {
		return
	}
	x := this.Tiles(tiles)
	ret = NewInstances().Add(x)
	return
}

func (this *Generator) InstancesFromString(str string) (ret tile.Instances, err error) {
	tiles, err := tile.NewTilesFromString(str)
	if err != nil {
		return
	}
	ret = this.Tiles(tiles)
	if ret == nil {
		err = errors.New("Incorrect input string")
	}
	return
}

func (this *Generator) Tiles(tiles tile.Tiles) (ret tile.Instances) {
	ret = make(tile.Instances, len(tiles))
	for i, t := range tiles {
		r := this.Instance(t)
		if r == tile.InstanceNull {
			return nil
		}
		ret[i] = r
	}
	return
}
