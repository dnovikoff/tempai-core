package shanten_test

import (
	"fmt"
	"log"

	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/shanten"
)

func ExampleCalculate() {
	generator := compact.NewTileGenerator()
	tiles, err := generator.CompactFromString("3567m5677p268s77z")
	if err != nil {
		log.Fatal(err)
	}
	if tiles.CountBits() != 13 {
		log.Fatal("Expected 13 tiles, but got", tiles.CountBits())
	}
	res := shanten.Calculate(tiles)
	fmt.Printf("Hand is %s\n", tiles.Instances())

	fmt.Printf("Regular shanten value is: %v\n", res.Regular.Value)
	fmt.Printf("Pairs shanten value is: %v\n", res.Pairs.Value)
	fmt.Printf("Kokushi shanten value is: %v\n", res.Kokushi.Value)
	fmt.Printf("Total shanten value is: %v\n", res.Total.Value)

	uke := res.Total.CalculateUkeIre(compact.NewTotals().Merge(tiles))
	fmt.Printf("Total uke ire: %v/%v\n", uke.UniqueCount(), uke.Count())
	fmt.Printf("Hand improves: %s\n", res.Total.Improves.Tiles())
	// Output:
	// Hand is 3567m5677p268s77z
	// Regular shanten value is: 2
	// Pairs shanten value is: 4
	// Kokushi shanten value is: 11
	// Total shanten value is: 2
	// Total uke ire: 18/63
	// Hand improves: 123458m456789p12347s7z
}
