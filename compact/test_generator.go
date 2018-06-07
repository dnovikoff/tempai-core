package compact

import (
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/stretchr/testify/require"
)

type TestGenerator struct {
	impl *Generator
	t    require.TestingT
}

func NewTestGenerator(t require.TestingT) *TestGenerator {
	return &TestGenerator{NewTileGenerator(), t}
}

func (tg *TestGenerator) TilesLeft() tile.Instances {
	return tg.impl.TilesLeft()
}

func (tg *TestGenerator) InstancesFromString(str string) tile.Instances {
	x, err := tg.impl.InstancesFromString(str)
	require.NoError(tg.t, err)
	return x
}
