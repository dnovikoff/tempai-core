package main

import (
	"fmt"
	"log"

	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/shanten"
)

func main() {
	generator := compact.NewTileGenerator()
	tiles, err := generator.CompactFromString("3567m5677p268s277z")
	if err != nil {
		log.Fatal(err)
	}
	// We tell the calculator, that there are 0 opened melds.
	results := shanten.CalculateShanten(tiles, 0, nil)
	fmt.Printf("Hand is %s\n", tiles.Instances())
	fmt.Printf("Shanten value is: %v\n", results.Value)
	fmt.Printf("Regular hand improves: %s (%v)\n", results.RegularImproves.Tiles(), results.RegularImproves.Count())
}
