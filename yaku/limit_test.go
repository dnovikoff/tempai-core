package yaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimitShostString(t *testing.T) {
	assert.Equal(t, "Mangan", LimitMangan.ShortString())
	assert.Equal(t, "Haneman", LimitHaneman.ShortString())
	assert.Equal(t, "Baiman", LimitBaiman.ShortString())
	assert.Equal(t, "Sanbaiman", LimitSanbaiman.ShortString())
	assert.Equal(t, "Yakuman", LimitYakuman.ShortString())
}
