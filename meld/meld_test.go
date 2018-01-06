package meld

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bitbucket.org/dnovikoff/tempai-core/compact"
)

func TestMeldExtract(t *testing.T) {
	tst := func(str string) string {
		g := compact.NewTileGenerator()
		tiles, err := g.CompactFromString(str)
		require.NoError(t, err)
		m := ExtractLastMeld(tiles)
		if m == 0 {
			return "Null"
		}
		ret := m.Instances().String()
		if !m.Waits().IsEmpty() {
			ret += " (" + m.Waits().Tiles().String() + ")"
		}
		return ret
	}
	assert.Equal(t, "12p (3p)", tst("12p"))
	assert.Equal(t, "13p (2p)", tst("13p"))
	assert.Equal(t, "23p (14p)", tst("23p"))
	assert.Equal(t, "89p (7p)", tst("89p"))

	assert.Equal(t, "88p (8p)", tst("88p"))
	assert.Equal(t, "8p (8p)", tst("8p"))

	assert.Equal(t, "Null", tst("888p"))
	assert.Equal(t, "Null", tst("14p"))
	assert.Equal(t, "Null", tst("1p2m"))
}
