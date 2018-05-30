package yaku

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dnovikoff/tempai-core/tile"
)

func NewRedYakuTester(t *testing.T, in string) *YakuTester {
	tester := NewYakuTester(t, in)
	tester.ctx.Rules = RulesTenhouRed()
	return tester
}

func TestYakuDora(t *testing.T) {
	assert.Equal(t, "2 = YakuChun: 1, YakuDora: 1", NewYakuTester(t, "123p345s789m4477z").IDora("9p").Ron(tile.Red).String())
	assert.Equal(t, "1 = YakuChun: 1", NewYakuTester(t, "123p345s789m4477z").IDora("3p").Ron(tile.Red).String())
	assert.Equal(t, "6 = YakuChun: 1, YakuDora: 5", NewYakuTester(t, "123p345s789m4477z").IDora("12p2s3z").Ron(tile.Red).String())

	assert.Nil(t, NewYakuTester(t, "123p345s789m44z79p").IDora("9p").Ron(tile.Pin8))
}

func TestYakuDoraRiichiOpensUra(t *testing.T) {
	assert.Equal(t, "3 = YakuChun: 1, YakuDora: 2", NewYakuTester(t, "123p345s789m4477z").IDora("12p").IUra("45s").Ron(tile.Red).String())
	assert.Equal(t, "5 = YakuChun: 1, YakuRiichi: 1, YakuDora: 2, YakuUraDora: 1", NewYakuTester(t, "123p345s789m4477z").IDora("12p").IUra("45s").Riichi().Ron(tile.Red).String())
}

func TestYakuDoraAka(t *testing.T) {
	assert.Equal(t, "2 = YakuChun: 1, YakuAkaDora: 1", NewRedYakuTester(t, "123p345s789m4477z").IDora("12z").Ron(tile.Red).String())

	// Aka is in dead
	assert.Equal(t, "1 = YakuChun: 1", NewRedYakuTester(t, "123p34s789m44777z").IDora("12z").IUra("5s").Ron(tile.Sou5).String())

	assert.Equal(t, "3 = YakuChun: 1, YakuAkaDora: 1, YakuDora: 1", NewRedYakuTester(t, "123p345s789m4477z").IDora("4s").Ron(tile.Red).String())
}
