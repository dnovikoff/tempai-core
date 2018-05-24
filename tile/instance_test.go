package tile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstanceSimple(t *testing.T) {
	tl := Man1
	assert.EqualValues(t, 1, tl)
	assert.EqualValues(t, 1, tl.Instance(0))
	assert.EqualValues(t, 2, tl.Instance(1))
	assert.EqualValues(t, 3, tl.Instance(2))
	assert.EqualValues(t, 4, tl.Instance(3))
}
