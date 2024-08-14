## About

[![CI](https://github.com/dnovikoff/tempai-core/workflows/CI/badge.svg?branch=master&event=push)](https://github.com/dnovikoff/tempai-core/actions?query=workflow%3ACI)
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

`examples/shanten/shanten_test.go`

```go
generator := compact.NewTileGenerator()
tiles, _ := generator.CompactFromString("3567m5677p268s77z")
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
Hand is 3567m5677p268s77z
Regular shanten value is: 2
Pairs shanten value is: 4
Kokushi shanten value is: 11
Total shanten value is: 2
Total uke ire: 18/63
Hand improves: 123458m456789p12347s7z
```

### Calculate tempai

Tempai results could be transformed into yaku results -> han/fu value -> score values

`examples/tempai/tempai_test.go`

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

`examples/effective/effective_test.go`

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

`examples/score/score_test.go`

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

`examples/yaku/yaku_test.go`

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

### Benchmarks
Benchmarks are located at `examples/bench`.
I've already made some code improvements, based on profiling for simular test code.
Although there could be more space for improve, the results on my machine seems quite fine to start with.

`cd ./examples/bench/ && go test -bench=. -benchtime 5s -benchmem -run notest`

Output:
```
goos: linux
goarch: amd64
pkg: github.com/dnovikoff/tempai-core/examples/bench
cpu: 12th Gen Intel(R) Core(TM) i7-1260P
BenchmarkShanten-16               419035             13633 ns/op            5924 B/op         29 allocs/op
--- BENCH: BenchmarkShanten-16
    bench_test.go:32: Repeat:  1
    bench_test.go:33: RPS:  33561.55188615921
    bench_test.go:32: Repeat:  100
    bench_test.go:33: RPS:  60719.62470414363
    bench_test.go:32: Repeat:  10000
    bench_test.go:33: RPS:  69839.23838410915
    bench_test.go:32: Repeat:  419035
    bench_test.go:33: RPS:  73352.4715081205
BenchmarkTempai-16                976167              6007 ns/op            3514 B/op         52 allocs/op
--- BENCH: BenchmarkTempai-16
    bench_test.go:58: Repeat:  1
    bench_test.go:59: Tempai hand count:  0
    bench_test.go:60: RPS:  30308.540946838817
    bench_test.go:58: Repeat:  100
    bench_test.go:59: Tempai hand count:  48
    bench_test.go:60: RPS:  79856.06742407485
    bench_test.go:58: Repeat:  10000
    bench_test.go:59: Tempai hand count:  4826
    bench_test.go:60: RPS:  119202.33090145081
    bench_test.go:58: Repeat:  715213
        ... [output truncated]
BenchmarkEffective-16              48742            123743 ns/op           21308 B/op        315 allocs/op
--- BENCH: BenchmarkEffective-16
    bench_test.go:81: Repeat:  1
    bench_test.go:82: RPS:  12198.543493906827
    bench_test.go:81: Repeat:  100
    bench_test.go:82: RPS:  8245.353351732016
    bench_test.go:81: Repeat:  10000
    bench_test.go:82: RPS:  8123.973175995651
    bench_test.go:81: Repeat:  48742
    bench_test.go:82: RPS:  8081.242170967126
PASS
ok      github.com/dnovikoff/tempai-core/examples/bench 24.435s
```
