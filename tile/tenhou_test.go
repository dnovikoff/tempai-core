package tile

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTenhouStringify(t *testing.T) {
	assert.Equal(t, "", TilesToTenhouString(nil))
	assert.Equal(t, "", TilesToTenhouString(Tiles{}))
	assert.Equal(t, "4s12p3z9p", TilesToTenhouString(Tiles{Sou4, Pin1, Pin2, West, Pin9}))
}

func TestTenhouTilesFromString(t *testing.T) {
	tst := func(str string) Tiles {
		ret, err := NewTilesFromString(str)
		require.NoError(t, err)
		return ret
	}

	assert.Nil(t, tst(""))
	assert.Equal(t, Tiles{Man1}, tst("1m"))
	assert.Equal(t, Tiles{Sou2, Sou3}, tst("23s"))
	assert.Equal(t, Tiles{Sou2, Sou3}, tst("2s3s"))
	assert.Equal(t, Tiles{Man1, Pin2, Sou1, East}, tst("1m2p1s1z"))
}

func TestTenhouTilesFromStringErrors(t *testing.T) {
	tst := func(str string) string {
		ret, err := NewTilesFromString(str)
		require.Nil(t, ret)
		require.Error(t, err)
		return err.Error()
	}

	assert.Contains(t, tst("1"), "Expected to end with a letter")
	assert.Contains(t, tst("1ps"), "Empty range at 2")
	assert.Contains(t, tst("p"), "Empty range at 0")
	assert.Contains(t, tst("0p"), "Unexpected symbol 0 at position 0")
	assert.Contains(t, tst("8z"), "Unexpected value 8 for type z")
}
