package tempai_test

import (
	"fmt"
	"log"

	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/tempai"
)

func ExampleCalculate() {
	generator := compact.NewTileGenerator()
	tiles, err := generator.CompactFromString("789m4466678p234s")
	if err != nil {
		log.Fatal(err)
	}
	results := tempai.Calculate(tiles)
	fmt.Printf("Hand is %s\n", tiles.Instances())
	fmt.Printf("Waits are %s\n", tempai.GetWaits(results).Tiles())
	// Output:
	// Hand is 789m4466678p234s
	// Waits are 469p
}
