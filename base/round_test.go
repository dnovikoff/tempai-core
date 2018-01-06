package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRound(t *testing.T) {
	first := FirstRound
	next := func() RoundID {
		first = first.Next()
		return first
	}
	assert.Equal(t, NewRound(WindEast, WindSouth), next())
	assert.Equal(t, NewRound(WindEast, WindWest), next())
	assert.Equal(t, NewRound(WindEast, WindNorth), next())
	assert.Equal(t, NewRound(WindSouth, WindEast), next())
	assert.Equal(t, NewRound(WindSouth, WindSouth), next())
	assert.Equal(t, NewRound(WindSouth, WindWest), next())
	assert.Equal(t, NewRound(WindSouth, WindNorth), next())
	assert.Equal(t, NewRound(WindWest, WindEast), next())
}
