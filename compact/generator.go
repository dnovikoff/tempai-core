package compact

import (
	"errors"

	"github.com/dnovikoff/tempai-core/tile"
)

type Generator struct {
	instances Instances
}

func NewTileGenerator() *Generator {
	return &Generator{AllInstances()}
}

func (g *Generator) TilesLeft() tile.Instances {
	return g.instances.Instances()
}

func (g *Generator) Instance(t tile.Tile) tile.Instance {
	m := g.instances.GetMask(t)
	if m.IsEmpty() {
		return tile.InstanceNull
	}
	ret := m.First()
	g.instances.Remove(ret)
	return ret
}

func (g *Generator) CompactFromString(str string) (Instances, error) {
	tiles, err := g.InstancesFromString(str)
	if err != nil {
		return nil, err
	}
	return NewInstances().Add(tiles), nil
}

func (g *Generator) InstancesFromString(str string) (tile.Instances, error) {
	tiles, err := tile.NewTilesFromString(str)
	if err != nil {
		return nil, err
	}
	ret := g.Tiles(tiles)
	if ret == nil {
		return nil, errors.New("Incorrect input string")
	}
	return ret, nil
}

func (g *Generator) Tiles(tiles tile.Tiles) tile.Instances {
	ret := make(tile.Instances, len(tiles))
	for i, t := range tiles {
		r := g.Instance(t)
		if r == tile.InstanceNull {
			return nil
		}
		ret[i] = r
	}
	return ret
}
