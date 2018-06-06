package compact

import (
	"errors"

	"github.com/dnovikoff/tempai-core/tile"
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

func (this *Generator) Instance(t tile.Tile) (ret tile.Instance) {
	m := this.tiles.GetMask(t)
	if m.IsEmpty() {
		return tile.InstanceNull
	}
	ret = m.First()
	this.tiles.Remove(ret)
	return
}

func (this *Generator) CompactFromString(str string) (Instances, error) {
	tiles, err := this.InstancesFromString(str)
	if err != nil {
		return nil, err
	}
	return NewInstances().Add(tiles), nil
}

func (this *Generator) InstancesFromString(str string) (tile.Instances, error) {
	tiles, err := tile.NewTilesFromString(str)
	if err != nil {
		return nil, err
	}
	ret := this.Tiles(tiles)
	if ret == nil {
		return nil, errors.New("Incorrect input string")
	}
	return ret, nil
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
