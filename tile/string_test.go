package tile

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	tst := func(in Tile) string {
		return Tiles{in}.String()
	}
	for i := 1; i <= 9; i++ {
		assert.Equal(t, strconv.Itoa(i)+"m", tst(Tile(i-1)+Man1))
		assert.Equal(t, strconv.Itoa(i)+"p", tst(Tile(i-1)+Pin1))
		assert.Equal(t, strconv.Itoa(i)+"s", tst(Tile(i-1)+Sou1))
	}
	assert.Equal(t, "1z", tst(East))
	assert.Equal(t, "2z", tst(South))
	assert.Equal(t, "3z", tst(West))
	assert.Equal(t, "4z", tst(North))

	assert.Equal(t, "5z", tst(White))
	assert.Equal(t, "6z", tst(Green))
	assert.Equal(t, "7z", tst(Red))
}
