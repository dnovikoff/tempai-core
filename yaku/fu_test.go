package yaku

import (
	"testing"

	"github.com/dnovikoff/tempai-core/hand/calc"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/stretchr/testify/assert"
)

func TestFu(t *testing.T) {
	for _, v := range []struct {
		name     string
		expected FuPoints
		meld     calc.Meld
	}{
		{"middle opened pon", 2, calc.Open(calc.Pon(tile.Man2))},
		{"middle closed pon", 4, calc.Pon(tile.Man2)},
		{"middle opened kan", 8, calc.Open(calc.Kan(tile.Man2))},
		{"middle closed kan", 16, calc.Kan(tile.Man2)},

		{"terminal opened pon", 4, calc.Open(calc.Pon(tile.Pin9))},
		{"terminal closed pon", 8, calc.Pon(tile.Pin9)},
		{"terminal opened kan", 16, calc.Open(calc.Kan(tile.Pin9))},
		{"terminal closed kan", 32, calc.Kan(tile.Pin9)},

		{"honor opened pon", 4, calc.Open(calc.Pon(tile.Red))},
		{"honor closed pon", 8, calc.Pon(tile.Red)},
		{"honor opened kan", 16, calc.Open(calc.Kan(tile.Red))},
		{"honor closed kan", 32, calc.Kan(tile.Red)},
	} {
		t.Run(v.name, func(t *testing.T) {
			assert.Equal(t, v.expected, GetMeldFu(v.meld))
		})
	}
}
