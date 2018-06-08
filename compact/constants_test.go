package compact

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dnovikoff/tempai-core/tile"
	"github.com/stretchr/testify/assert"
)

func TestAllTiles(t *testing.T) {
	tile := tile.TileBegin
	next := func() string {
		ret := fmt.Sprintf("[%d] %v=%v", tile, tile, tileStatus(tile))
		tile++
		return ret
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

func TestConst(t *testing.T) {
	tst := func(x Tiles) string {
		return x.Tiles().String()
	}

	assert.Equal(t, "19m19p19s", tst(Terminal))
	assert.Equal(t, "1234z", tst(Wind))
	assert.Equal(t, "567z", tst(Dragon))

	assert.Equal(t, "123456789m", tst(Man))
	assert.Equal(t, "123456789p", tst(Pin))
	assert.Equal(t, "123456789s", tst(Sou))
	assert.Equal(t, "123456789m123456789p123456789s", tst(Sequence))

	assert.Equal(t, "1234567z", tst(Honor))
	assert.Equal(t, "19m19p19s1234567z", tst(TerminalOrHonor))

	assert.Equal(t, "2345678m2345678p2345678s", tst(Middle))
}

// Terminal = Sou1 | Sou9 | Man1 | Man9 | Pin1 | Pin9
// Wind     = East | South | West | North
// Dragon   = White | Green | Red

// Man = Man1 | Man2 | Man3 | Man4 | Man5 | Man6 | Man7 | Man8 | Man9
// Pin = Pin1 | Pin2 | Pin3 | Pin4 | Pin5 | Pin6 | Pin7 | Pin8 | Pin9
// Sou = Sou1 | Sou2 | Sou3 | Sou4 | Sou5 | Sou6 | Sou7 | Sou8 | Sou9

// Sequence = Man | Pin | Sou

// Honor           = Wind | Dragon
// TerminalOrHonor = Terminal | Honor

// AllTiles = TileEnd - 1
// Middle   = AllTiles ^ TerminalOrHonor

// GreenYakuman = Sou2 | Sou3 | Sou4 | Sou6 | Sou8 | Green

func runeBool(l rune, in bool) string {
	if in {
		return strings.ToUpper(string(l))
	}
	return string(l)
}

func tileStatus(t tile.Tile) string {
	ret := ""
	ret += runeBool('s', Sequence.Check(t))
	ret += runeBool('m', Middle.Check(t))
	ret += runeBool('t', Terminal.Check(t))
	ret += runeBool('h', Honor.Check(t))
	ret += runeBool('o', TerminalOrHonor.Check(t))
	return ret
}
