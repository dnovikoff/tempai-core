package yaku

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dnovikoff/tempai-core/tile"
)

func NewRedYakuTester(t *testing.T, in string) *yakuTester {
	tester := newYakuTester(t, in)
	tester.ctx.Rules = RulesTenhouRed()
	return tester
}

func TestYakuDora(t *testing.T) {
	assert.Equal(t, "2 = YakuChun: 1, YakuDora: 1", newYakuTester(t, "123p345s789m4477z").iDora("9p").ron(tile.Red).String())
	assert.Equal(t, "1 = YakuChun: 1", newYakuTester(t, "123p345s789m4477z").iDora("3p").ron(tile.Red).String())
	assert.Equal(t, "6 = YakuChun: 1, YakuDora: 5", newYakuTester(t, "123p345s789m4477z").iDora("12p2s3z").ron(tile.Red).String())

	assert.Nil(t, newYakuTester(t, "123p345s789m44z79p").iDora("9p").ron(tile.Pin8))
}

func TestYakuDoraRiichiOpensUra(t *testing.T) {
	assert.Equal(t, "3 = YakuChun: 1, YakuDora: 2", newYakuTester(t, "123p345s789m4477z").iDora("12p").iUra("45s").ron(tile.Red).String())
	assert.Equal(t, "5 = YakuChun: 1, YakuRiichi: 1, YakuDora: 2, YakuUraDora: 1", newYakuTester(t, "123p345s789m4477z").iDora("12p").iUra("45s").riichi().ron(tile.Red).String())
}

func TestYakuDoraAka(t *testing.T) {
	assert.Equal(t, "2 = YakuChun: 1, YakuAkaDora: 1", NewRedYakuTester(t, "123p345s789m4477z").iDora("12z").ron(tile.Red).String())

	// Aka is in dead
	assert.Equal(t, "1 = YakuChun: 1", NewRedYakuTester(t, "123p34s789m44777z").iDora("12z").iUra("5s").ron(tile.Sou5).String())

	assert.Equal(t, "3 = YakuChun: 1, YakuAkaDora: 1, YakuDora: 1", NewRedYakuTester(t, "123p345s789m4477z").iDora("4s").ron(tile.Red).String())
}
