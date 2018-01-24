package main

import (
	"fmt"
	"log"

	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/shanten"
)

func main() {
	generator := compact.NewTileGenerator()
	// https://tempai.net/en/eff/5677m4456899p25s3z
	tiles, err := generator.CompactFromString("5677m4456899p25s3z")
	if err != nil {
		log.Fatal(err)
	}
	results := shanten.CalculateEffectivity(tiles, 0, nil)
	fmt.Printf("Hand is %s\n", tiles.Instances())
	best := results.Best()
	fmt.Printf("Best tiles is %v\n", best.Tile)
	fmt.Printf("Best shanten: %v\n", best.Shanten.Value)
}
