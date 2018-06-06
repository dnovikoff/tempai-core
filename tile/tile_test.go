package tile

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOneTile(t *testing.T) {
	tile := Man1
	assert.Equal(t, "1 1 1", fmt.Sprintf("%d %d %v", tile, tile.Type(), tile.Number()))
	tile = White
	assert.Equal(t, "32 32 1", fmt.Sprintf("%d %d %v", tile, tile.Type(), tile.Number()))
}

func TestConvert(t *testing.T) {
	for k := TileBegin; k < TileEnd; k++ {
		for c := CopyID(0); c < 4; c++ {
			t.Run(k.String()+"_"+strconv.Itoa(int(c)), func(t *testing.T) {
				i := k.Instance(c)
				assert.Equal(t, k, i.Tile())
				assert.Equal(t, c, i.CopyID())
			})
		}
	}
}

func TestAllTiles(t *testing.T) {
	tile := TileBegin
	next := func() (ret string) {
		ret = fmt.Sprintf("[%d] %v=%v", tile, tile, tileStatus(tile))
		tile++
		return
	}
	assert.Equal(t, "[1] 1m=SmThO", next())
	assert.Equal(t, "[2] 2m=SMtho", next())
	assert.Equal(t, "[3] 3m=SMtho", next())
	assert.Equal(t, "[4] 4m=SMtho", next())
	assert.Equal(t, "[5] 5m=SMtho", next())
	assert.Equal(t, "[6] 6m=SMtho", next())
	assert.Equal(t, "[7] 7m=SMtho", next())
	assert.Equal(t, "[8] 8m=SMtho", next())
	assert.Equal(t, "[9] 9m=SmThO", next())

	// pin
	assert.Equal(t, "[10] 1p=SmThO", next())
	assert.Equal(t, "[11] 2p=SMtho", next())
	assert.Equal(t, "[12] 3p=SMtho", next())
	assert.Equal(t, "[13] 4p=SMtho", next())
	assert.Equal(t, "[14] 5p=SMtho", next())
	assert.Equal(t, "[15] 6p=SMtho", next())
	assert.Equal(t, "[16] 7p=SMtho", next())
	assert.Equal(t, "[17] 8p=SMtho", next())
	assert.Equal(t, "[18] 9p=SmThO", next())

	// sou
	assert.Equal(t, "[19] 1s=SmThO", next())
	assert.Equal(t, "[20] 2s=SMtho", next())
	assert.Equal(t, "[21] 3s=SMtho", next())
	assert.Equal(t, "[22] 4s=SMtho", next())
	assert.Equal(t, "[23] 5s=SMtho", next())
	assert.Equal(t, "[24] 6s=SMtho", next())
	assert.Equal(t, "[25] 7s=SMtho", next())
	assert.Equal(t, "[26] 8s=SMtho", next())
	assert.Equal(t, "[27] 9s=SmThO", next())

	// wind
	assert.Equal(t, "[28] 1z=smtHO", next())
	assert.Equal(t, "[29] 2z=smtHO", next())
	assert.Equal(t, "[30] 3z=smtHO", next())
	assert.Equal(t, "[31] 4z=smtHO", next())

	// dragon
	assert.Equal(t, "[32] 5z=smtHO", next())
	assert.Equal(t, "[33] 6z=smtHO", next())
	assert.Equal(t, "[34] 7z=smtHO", next())
}

func TestTileStatic(t *testing.T) {
	assert.Equal(t, 34, TileCount)
	assert.Equal(t, 136, InstanceCount)
	assert.EqualValues(t, 0, TypeNull)
}

func TestInstanceToTile(t *testing.T) {
	assert.EqualValues(t, 1, Man1)
	assert.EqualValues(t, 1, Man1.Instance(0))
	assert.EqualValues(t, 1, Man1.Instance(0).Tile())
}

func TestTileIndicates(t *testing.T) {
	assert.Equal(t, Man2, Man1.Indicates())
	assert.Equal(t, Man3, Man2.Indicates())
	assert.Equal(t, Man4, Man3.Indicates())
	assert.Equal(t, Man5, Man4.Indicates())
	assert.Equal(t, Man6, Man5.Indicates())
	assert.Equal(t, Man7, Man6.Indicates())
	assert.Equal(t, Man8, Man7.Indicates())
	assert.Equal(t, Man9, Man8.Indicates())
	assert.Equal(t, Man1, Man9.Indicates())

	assert.Equal(t, Pin2, Pin1.Indicates())
	assert.Equal(t, Pin3, Pin2.Indicates())
	assert.Equal(t, Pin4, Pin3.Indicates())
	assert.Equal(t, Pin5, Pin4.Indicates())
	assert.Equal(t, Pin6, Pin5.Indicates())
	assert.Equal(t, Pin7, Pin6.Indicates())
	assert.Equal(t, Pin8, Pin7.Indicates())
	assert.Equal(t, Pin9, Pin8.Indicates())
	assert.Equal(t, Pin1, Pin9.Indicates())

	assert.Equal(t, Sou2, Sou1.Indicates())
	assert.Equal(t, Sou3, Sou2.Indicates())
	assert.Equal(t, Sou4, Sou3.Indicates())
	assert.Equal(t, Sou5, Sou4.Indicates())
	assert.Equal(t, Sou6, Sou5.Indicates())
	assert.Equal(t, Sou7, Sou6.Indicates())
	assert.Equal(t, Sou8, Sou7.Indicates())
	assert.Equal(t, Sou9, Sou8.Indicates())
	assert.Equal(t, Sou1, Sou9.Indicates())

	assert.Equal(t, South, East.Indicates())
	assert.Equal(t, West, South.Indicates())
	assert.Equal(t, North, West.Indicates())
	assert.Equal(t, East, North.Indicates())

	assert.Equal(t, Green, White.Indicates())
	assert.Equal(t, Red, Green.Indicates())
	assert.Equal(t, White, Red.Indicates())
}

func TestTilesContains(t *testing.T) {
	ts := Tiles{Pin4, White, Sou3}
	assert.True(t, ts.Contains(Pin4))
	assert.True(t, ts.Contains(White))

	assert.False(t, ts.Contains(Pin3))
	assert.False(t, ts.Contains(Pin5))
}

func TestTileImproves(t *testing.T) {
}

func runeBool(l rune, in bool) string {
	if in {
		return strings.ToUpper(string(l))
	}
	return string(l)
}

func tileStatus(t Tile) (ret string) {
	ret += runeBool('s', t.IsSequence())
	ret += runeBool('m', t.IsMiddle())
	ret += runeBool('t', t.IsTerminal())
	ret += runeBool('h', t.IsHonor())
	ret += runeBool('o', t.IsTerminalOrHonor())
	return
}
