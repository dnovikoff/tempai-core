package tile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypes(t *testing.T) {
	assert.Equal(t, TypeMan, Man4.Type())
	assert.Equal(t, TypeSou, Sou9.Type())
	assert.Equal(t, TypePin, Pin1.Type())
	assert.Equal(t, TypeWind, West.Type())
	assert.Equal(t, TypeDragon, White.Type())
	assert.Equal(t, TypeNull, TileNull.Type())
	assert.Equal(t, TypeNull, TileEnd.Type())

	assert.Equal(t, TypeNull, Tile(-900).Type())
	assert.Equal(t, TypeNull, Tile(1400).Type())
}

func TestTypeTile(t *testing.T) {
	assert.Equal(t, Man4, TypeMan.Tile(4))
	assert.Equal(t, East, TypeWind.Tile(1))
}

func TestTypeRunes(t *testing.T) {
	assert.Equal(t, 'm', TypeRune(TypeMan))
	assert.Equal(t, 'z', TypeRune(TypeDragon))
	assert.Equal(t, 'z', TypeRune(TypeWind))
	assert.Equal(t, '-', TypeRune(TypeNull))
	assert.Equal(t, '-', TypeRune(400))
	assert.Equal(t, '-', TypeRune(-256))
}
