package main

import (
	"fmt"
	"log"

	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/shanten"
)

func main() {
	generator := compact.NewTileGenerator()
	tiles, err := generator.CompactFromString("3567m5677p268s277z")
	if err != nil {
		log.Fatal(err)
	}
	// We tell the calculator, that there are 0 opened melds.
	res := shanten.Calculate(tiles)
	fmt.Printf("Hand is %s\n", tiles.Instances())

	fmt.Printf("Regular shanten value is: %v\n", res.Regular.Value)
	fmt.Printf("Pairs shanten value is: %v\n", res.Pairs.Value)
	fmt.Printf("Kokushi shanten value is: %v\n", res.Kokushi.Value)
	fmt.Printf("Total shanten value is: %v\n", res.Total.Value)

	uke := res.Total.CalculateUkeIre(compact.NewTotals().Merge(tiles))
	fmt.Printf("Total uke ire: %v/%v\n", uke.UniqueCount(), uke.Count())
	fmt.Printf("Hand improves: %s\n", res.Total.Improves.Tiles())
}
