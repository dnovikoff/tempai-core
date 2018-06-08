package calc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dnovikoff/tempai-core/tile"
)

func TestMeldsComletePon(t *testing.T) {
	assert.Equal(t, []string{
		"111m", "222m", "333m",
		"444m", "555m", "666m",
		"777m", "888m", "999m",
		"111p", "222p", "333p",
		"444p", "555p", "666p",
		"777p", "888p", "999p",
		"111s", "222s", "333s",
		"444s", "555s", "666s",
		"777s", "888s", "999s",
		"111z", "222z", "333z", "444z",
		"555z", "666z", "777z",
	}, DebugMelds(createAllPon().Clone()))
}

func TestMeldsComleteChi(t *testing.T) {
	assert.Equal(t, []string{
		"123m", "234m", "345m",
		"456m", "567m", "678m", "789m",
		"123p", "234p", "345p",
		"456p", "567p", "678p", "789p",
		"123s", "234s", "345s",
		"456s", "567s", "678s", "789s",
	}, DebugMelds(createAllChi().Clone()))
}

func TestMeldsPonParts(t *testing.T) {
	assert.Equal(t, []string{
		"11m (1m)", "22m (2m)", "33m (3m)",
		"44m (4m)", "55m (5m)", "66m (6m)",
		"77m (7m)", "88m (8m)", "99m (9m)",
		"11p (1p)", "22p (2p)", "33p (3p)",
		"44p (4p)", "55p (5p)", "66p (6p)",
		"77p (7p)", "88p (8p)", "99p (9p)",
		"11s (1s)", "22s (2s)", "33s (3s)",
		"44s (4s)", "55s (5s)", "66s (6s)",
		"77s (7s)", "88s (8s)", "99s (9s)",
		"11z (1z)", "22z (2z)", "33z (3z)", "44z (4z)",
		"55z (5z)", "66z (6z)", "77z (7z)",
	}, DebugMelds(createAllPonParts().Clone()))
}

func TestChiParts(t *testing.T) {
	assert.Equal(t, []string{
		"12m (3m)", "13m (2m)",
		"23m (14m)", "24m (3m)",
		"34m (25m)", "35m (4m)",
		"45m (36m)", "46m (5m)",
		"56m (47m)", "57m (6m)",
		"67m (58m)", "68m (7m)",
		"78m (69m)", "79m (8m)",
		"89m (7m)",

		"12p (3p)", "13p (2p)",
		"23p (14p)", "24p (3p)",
		"34p (25p)", "35p (4p)",
		"45p (36p)", "46p (5p)",
		"56p (47p)", "57p (6p)",
		"67p (58p)", "68p (7p)",
		"78p (69p)", "79p (8p)",
		"89p (7p)",

		"12s (3s)", "13s (2s)",
		"23s (14s)", "24s (3s)",
		"34s (25s)", "35s (4s)",
		"45s (36s)", "46s (5s)",
		"56s (47s)", "57s (6s)",
		"67s (58s)", "68s (7s)",
		"78s (69s)", "79s (8s)",
		"89s (7s)",
	}, DebugMelds(createAllChiParts().Clone()))
}

func TestChiPartsComplete(t *testing.T) {
	c := ChiPart1(tile.Man5)
	assert.Equal(t, "56m (47m)", DebugMeld(c))
	assert.Equal(t, "456m", DebugMeld(c.Complete(tile.Man4)))
	assert.Equal(t, "567m", DebugMeld(c.Complete(tile.Man7)))

	assert.Nil(t, c.Complete(tile.Man5))
	assert.Nil(t, c.Complete(tile.Man6))
	assert.Nil(t, c.Complete(tile.Man3))
	assert.Nil(t, c.Complete(tile.Sou4))
}

func TestPonPartsComplete(t *testing.T) {
	c := PonPart(tile.Man5)
	assert.Equal(t, "55m (5m)", DebugMeld(c))
	assert.Equal(t, "555m", DebugMeld(c.Complete(tile.Man5)))

	assert.Nil(t, c.Complete(tile.Man6))
	assert.Nil(t, c.Complete(tile.Man3))
	assert.Nil(t, c.Complete(tile.Sou4))
}

func TestTankiComplete(t *testing.T) {
	c := Tanki(tile.Man5)
	assert.Equal(t, "5m (5m)", DebugMeld(c))
	assert.Equal(t, "55m", DebugMeld(c.Complete(tile.Man5)))

	assert.Nil(t, c.Complete(tile.Man6))
	assert.Nil(t, c.Complete(tile.Man3))
	assert.Nil(t, c.Complete(tile.Sou4))
}
