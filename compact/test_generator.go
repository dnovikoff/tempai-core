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

func (this *TestGenerator) TilesLeft() tile.Instances {
	return this.impl.TilesLeft()
}

func (this *TestGenerator) InstancesFromString(str string) tile.Instances {
	x, err := this.impl.InstancesFromString(str)
	require.NoError(this.t, err)
	return x
}
