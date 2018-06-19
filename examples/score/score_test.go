package score_test

import (
	"fmt"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/score"
)

func ExampleGetScore() {
	s := score.GetScore(score.RulesEMA(), 4, 22, 0)
	fmt.Printf("Hand value is %v.%v (%v)\n", s.Han, s.Fu, s.Fu.Round())
	fmt.Printf("Dealer ron: %v\n", s.PayRonDealer)
	fmt.Printf("Dealer tsumo: %v all\n", s.PayTsumoDealer)
	fmt.Printf("Ron: %v\n", s.PayRon)
	fmt.Printf("Tsumo: %v/%v\n", s.PayTsumoDealer, s.PayTsumo)
	changes := s.GetChanges(base.WindEast, base.WindEast, 0)
	fmt.Printf("Total for dealer tsumo is %v\n", changes.TotalWin())
	// Output:
	// Hand value is 4.22 (30)
	// Dealer ron: 11600
	// Dealer tsumo: 3900 all
	// Ron: 7700
	// Tsumo: 3900/2000
	// Total for dealer tsumo is 11700
}
