## About

[![Build Status](https://travis-ci.org/dnovikoff/tempai-core.svg?branch=master)](https://travis-ci.org/dnovikoff/tempai-core)
[![Coverage Status](https://img.shields.io/codecov/c/github/dnovikoff/tempai-core.svg)](https://codecov.io/gh/dnovikoff/tempai-core)
[![Go Report Card](https://goreportcard.com/badge/github.com/dnovikoff/tempai-core)](https://goreportcard.com/report/github.com/dnovikoff/tempai-core)

This is riichi mahjong Golang package.
Package supports calculation for:

1. Shanten value
2. Effectivity drops
3. Tempai
4. Hand yaku
4. Han/Fu hand value based on yaku
5. Scroing base on han/fu value

The game itself is NOT the purpose of this particular package.

This package is extracted from a larger repo(private) of my code and it is been quite tested. 
I've found a good number of corner-cases to include in my tests.
You can see the example of effetivity calculator, based on this package here: https://tempai.net/en/eff .

I've also validated my code against more than 1 million tenhou Phoenix replays, downloaded from tenhou.net server.
So it seems that the results of calculation could be trusted.

## Installation

`go get github.com/dnovikoff/tempai-core/...`

## Quick Start by exmaples

All examples could be found in `example` folder.
Note that provided examples are simple and do not cover all possible package functionality.
You can also explore `*_test.go` test files to search for more usage examples.

All example hand inputs provided in tenhou-style string format:

1. `123456789s` for Sou
2. `123456789m` for Man
3. `123456789p` for Pin
4. `1234z` for East, South, West, North
5. `567z` for White, Green, Red

### Calculate shanten number

Shanten, Tempai and effectivity calculators support different forms separatly or all-together:

1. Regular hand
2. Seven pairs
3. Kokushi

You can also take into considiration any number of visible tiles to affect results.
Calculating hands with opened melds is also supported.

`go run ./examples/shanten/main.go`

```go
generator := compact.NewTileGenerator()
tiles, _ := generator.CompactFromString("3567m5677p268s277z")
res := shanten.Calculate(tiles)
fmt.Printf("Hand is %s\n", tiles.Instances())

fmt.Printf("Regular shanten value is: %v\n", res.Regular.Value)
fmt.Printf("Pairs shanten value is: %v\n", res.Pairs.Value)
fmt.Printf("Kokushi shanten value is: %v\n", res.Kokushi.Value)
fmt.Printf("Total shanten value is: %v\n", res.Total.Value)

uke := res.Total.CalculateUkeIre(compact.NewTotals().Merge(tiles))
fmt.Printf("Total uke ire: %v/%v\n", uke.UniqueCount(), uke.Count())
fmt.Printf("Hand improves: %s\n", res.Total.Improves.Tiles())
```

Output:
```
Hand is 3567m5677p268s277z
Regular shanten value is: 2
Pairs shanten value is: 4
Kokushi shanten value is: 10
Total shanten value is: 2
Total uke ire: 19/66
Hand improves: 123458m456789p12347s27z
```

### Calculate tempai

Tempai results could be transformed into yaku results -> han/fu value -> score values

`go run ./examples/tempai/main.go`

```go
generator := compact.NewTileGenerator()
tiles, _ := generator.CompactFromString("789m4466678p234s")
results := tempai.Calculate(tiles)
fmt.Printf("Hand is %s\n", tiles.Instances())
fmt.Printf("Waits are %s\n", results.Waits().Tiles())
```

Output:
```
Hand is 789m4466678p234s
Waits are 469p
```

### Calculate tile effectivity

Calculating Uke-Ure value for hand.

`go run ./examples/effective/main.go`

```go
generator := compact.NewTileGenerator()
tiles, _ := generator.CompactFromString("5677m4456899p25s3z")
results := effective.Calculate(tiles)
fmt.Printf("Hand is %s\n", tiles.Instances())
best := results.Sorted(tiles).Best()
fmt.Printf("Best to drop is %v\n", best.Tile)
fmt.Printf("Best shanten: %v\n", best.Shanten.Total.Value)
```

Output:
```
Hand is 5677m4456899p25s3z
Best to drop is 3z
Best shanten: 3
```

### Calculate han+fu value

Package supports configuration of different rulesets.
Configuration options includes:

1. Mangan round (4.30/3.60 could be rounded to mangan)
2. Yakuman summ option
3. Double yakuman option
4. Kazoe Yakuman/Sanbaiman
5. Changing hoba value

Included rulesets for: EMA, JPML-A, JPLML-B, Tenhou.
You can also configure your own ruleset

`go run ./examples/score/main.go`

```go
s := score.GetScore(score.RulesEMA(), 4, 22, 0)
fmt.Printf("Hand value is %v.%v (%v)\n", s.Han, s.Fu, s.Fu.Round())
fmt.Printf("Dealer ron: %v\n", s.PayRonDealer)
fmt.Printf("Dealer tsumo: %v all\n", s.PayTsumoDealer)
fmt.Printf("Ron: %v\n", s.PayRon)
fmt.Printf("Tsumo: %v/%v\n", s.PayTsumoDealer, s.PayTsumo)
changes := s.GetChanges(base.WindEast, base.WindEast, 0)
fmt.Printf("Total for dealer tsumo is %v\n", changes.TotalWin())
```

Output:
```
Hand value is 4.22 (30)
Dealer ron: 11600
Dealer tsumo: 3900 all
Ron: 7700
Tsumo: 3900/2000
Total for dealer tsumo is 11700
```

### Calculate yaku for hand

Package supports configuration of different rulesets.
Configuration options includes:

1. Setting any number of akkadors (not limited by red fives)
2. Renhou could be configured as yakuman or mangan
3. Uradoras could be disabled (for JPML-A)
4. Ipatsu could be disabled (for JPML-A)
5. Haitei could be combined with Rinshan
6. Enabling/Disabling open Tanyao

Included rulesets for: EMA, JPML-A, JPLML-B, Tenhou with red fives.
You can also configure your own ruleset

`go run ./examples/yaku/main.go`

```go
generator := compact.NewTileGenerator()
tiles, _ := generator.CompactFromString("33z123m456p66778s")
winTile := generator.Instance(tile.Sou5)

results := tempai.Calculate(tiles).Index()
ctx := &yaku.Context{
    Tile:      winTile,
    Rules:     yaku.RulesEMA(),
    IsTsumo:   true,
    IsChankan: true,
}
yakuResult := yaku.Win(results, ctx, nil)
fmt.Printf("%v\n", yakuResult.Yaku.String())
fmt.Printf("Value: %v.%v\n", yakuResult.Sum(), yakuResult.Fus.Sum())
```

Output:
```
YakuChankan: 1, YakuPinfu: 1, YakuTsumo: 1
Value: 3.20
```

### Perfomance
There is an `examples/perfomance` folder with some mesaurement program.
I've already made some code improvements, based on profiling for simular test code.
Although there could be more space for improve, the results on my machine seems quite fine to start with.

` go run ./examples/perfomance/main.go`

Output:
```
================== Test shanten
Repeat: 10000
Elapsed: 293.791605ms
Estemated speed: 34037.7322898658 per second
================== Test tempai
Repeat: 10000
Elapsed: 149.649275ms
Estemated speed: 66822.9097668532 per second
Tempai hand count: 4910
================== Test effectivity
Repeat: 1000
Elapsed: 324.172207ms
Estemated speed: 3084.780182898283 per second
```
