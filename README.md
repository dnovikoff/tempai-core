## About
This is riichi mahjong Golang package

## Index

### base
Contents basic object need for other packages: Opponent, Wind, Round

### tile
Tile representation and string in tenhou format parser-generator

### compact
bit-based representation for tile collections

### meld
Chi, Pon, Kan objects

### shanten
Calculators for shanten and tempai

### yaku
Calculate yaku, based on calculated forms and context

### score
Score calculation and scoring rules

### examples
Here you can find short usage examples

## Installation

`go get bitbucket.org/dnovikoff/tempai-core/...`

## Quick Start

All examples could be found in example folder

Calculate shanten number
```go
generator := compact.NewTileGenerator()
tiles, _ := generator.CompactFromString("3567m5677p268s277z")
results := shanten.CalculateShanten(tiles, 0, nil)
fmt.Printf("Hand is %s\n", tiles.Instances())
fmt.Printf("Shanten value is: %v\n", results.Value)
fmt.Printf("Regular hand improves: %s (%v)\n", results.RegularImproves.Tiles(), results.RegularImproves.Count())
```

Calculate tempai
```go
generator := compact.NewTileGenerator()
tiles, _ := generator.CompactFromString("789m4466678p234s")
results := shanten.CalculateTempai(tiles, nil)
fmt.Printf("Hand is %s\n", tiles.Instances())
fmt.Printf("Waits are %s\n", results.Waits().Tiles())
```

Calculate tile effectivity
```go
generator := compact.NewTileGenerator()
tiles, _ := generator.CompactFromString("5677m4456899p25s3z")
results := shanten.CalculateEffectivity(tiles, 0, nil)
fmt.Printf("Hand is %s\n", tiles.Instances())
best := results.Best()
fmt.Printf("Best tiles is %v\n", best.Tile)
fmt.Printf("Best shanten: %v\n", best.Shanten.Value)
```

Calculate han+fu value
```go
rules := score.RulesEMA
s := rules.GetScore(4, 22, 0)
fmt.Printf("Hand value is %v.%v (%v)\n", s.Han, s.Fu, s.Fu.Round())
fmt.Printf("Dealer ron: %v\n", s.PayRonDealer)
fmt.Printf("Dealer tsumo: %v all\n", s.PayTsumoDealer)
fmt.Printf("Ron: %v\n", s.PayRon)
fmt.Printf("Tsumo: %v/%v\n", s.PayTsumoDealer, s.PayTsumo)
changes := s.GetChanges(base.WindEast, base.WindEast, 0)
fmt.Printf("Total for dealer tsumo is %v\n", changes.TotalWin())
```

Calculate yaku for hand
```go
generator := compact.NewTileGenerator()
tiles, _ := generator.CompactFromString("33z123m456p66778s")
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
```