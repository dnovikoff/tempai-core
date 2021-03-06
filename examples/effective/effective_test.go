package effective_test

import (
	"fmt"
	"log"

	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/effective"
)

func ExampleCalculate() {
	generator := compact.NewTileGenerator()
	// https://tempai.net/en/eff/5677m4456899p25s3z
	tiles, err := generator.CompactFromString("5677m4456899p25s3z")
	if err != nil {
		log.Fatal(err)
	}
	results := effective.Calculate(tiles)
	fmt.Printf("Hand is %s\n", tiles.Instances())
	best := results.Sorted(tiles).Best()
	fmt.Printf("Best to drop is %v\n", best.Tile)
	fmt.Printf("Best shanten: %v\n", best.Shanten.Total.Value)
	// Output:
	// Hand is 5677m4456899p25s3z
	// Best to drop is 3z
	// Best shanten: 3
}
