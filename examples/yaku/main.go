package main

import (
	"fmt"
	"log"

	"bitbucket.org/dnovikoff/tempai-core/compact"
	"bitbucket.org/dnovikoff/tempai-core/shanten"
	"bitbucket.org/dnovikoff/tempai-core/tile"
	"bitbucket.org/dnovikoff/tempai-core/yaku"
)

func main() {
	generator := compact.NewTileGenerator()
	tiles, err := generator.CompactFromString("33z123m456p66778s")
	if err != nil {
		log.Fatal(err)
	}
	winTile := generator.Instance(tile.Sou5)

	results := shanten.CalculateTempai(tiles, nil).Index()
	ctx := &yaku.Context{
		Tile:      winTile,
		Rules:     &yaku.RulesEMA,
		IsTsumo:   true,
		IsChankan: true,
	}
	yakuResult := yaku.Win(results, ctx)
	fmt.Printf("%v\n", yakuResult.Yaku.String())
	fmt.Printf("Value: %v.%v\n", yakuResult.Sum(), yakuResult.Fus.Sum())
}
