package main

import (
	"fmt"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/score"
)

func main() {
	rules := score.RulesEMA
	s := rules.GetScore(4, 22, 0)
	fmt.Printf("Hand value is %v.%v (%v)\n", s.Han, s.Fu, s.Fu.Round())
	fmt.Printf("Dealer ron: %v\n", s.PayRonDealer)
	fmt.Printf("Dealer tsumo: %v all\n", s.PayTsumoDealer)
	fmt.Printf("Ron: %v\n", s.PayRon)
	fmt.Printf("Tsumo: %v/%v\n", s.PayTsumoDealer, s.PayTsumo)
	changes := s.GetChanges(base.WindEast, base.WindEast, 0)
	fmt.Printf("Total for dealer tsumo is %v\n", changes.TotalWin())
}
