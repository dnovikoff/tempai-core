package effective

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortId(t *testing.T) {
	assert.True(t, newSortId(12, 11, 2).betterThan(newSortId(11, 11, 2)))
	assert.False(t, newSortId(10, 11, 2).betterThan(newSortId(11, 11, 2)))

	assert.True(t, newSortId(11, 11, 1).betterThan(newSortId(11, 11, 2)))
	assert.False(t, newSortId(11, 11, 3).betterThan(newSortId(11, 11, 2)))

	assert.True(t, newSortId(11, 12, 2).betterThan(newSortId(11, 11, 2)))
	assert.False(t, newSortId(11, 10, 2).betterThan(newSortId(11, 11, 2)))

	assert.True(t, newSortId(80, 1, 2).betterThan(newSortId(1, 90, 2)))
}
