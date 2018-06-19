package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"sort"
	"time"

	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/effective"
	"github.com/dnovikoff/tempai-core/hand/shanten"
	"github.com/dnovikoff/tempai-core/hand/tempai"
	"github.com/dnovikoff/tempai-core/tile"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var fShanten = flag.Bool("shanten", false, "shanten test")
var fTempai = flag.Bool("tempai", false, "tempai test")
var fEff = flag.Bool("eff", false, "effective test")
var fCount = flag.Int("count", 10000, "Count of tests")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if !(*fShanten || *fTempai || *fEff) {
		*fShanten = true
		*fTempai = true
		*fEff = true
	}
	if *fShanten {
		testShanten()
	}
	if *fTempai {
		testTempai()
	}
	if *fEff {
		testEffective()
	}
}

func testShanten() {
	repeat := *fCount
	data := make([]compact.Instances, repeat)
	// prepare
	source := rand.NewSource(123)
	rnd := rand.New(source)
	instances := compact.AllInstances().Instances()
	for k := range data {
		shuffle(rnd, instances)
		data[k] = compact.NewInstances().Add(instances[:13])
	}
	start := time.Now()
	for _, v := range data {
		shanten.Calculate(v)

	}
	elapsed := time.Since(start)
	fmt.Println("================== Test shanten")
	fmt.Printf("Repeat: %v\n", repeat)
	fmt.Printf("Elapsed: %v\n", elapsed)
	fmt.Printf("Estemated speed: %v per second\n", float64(repeat)/elapsed.Seconds())
}

func testTempai() {
	repeat := *fCount
	data := make([]compact.Instances, repeat)
	// prepare
	source := rand.NewSource(123)
	rnd := rand.New(source)
	instances := compact.AllInstancesFromTo(tile.Sou1, tile.Sou9+1).Instances()
	for k := range data {
		shuffle(rnd, instances)
		data[k] = compact.NewInstances().Add(instances[:13])
	}
	cnt := 0
	start := time.Now()
	for _, v := range data {
		r := tempai.Calculate(v)
		if r != nil {
			cnt++
		}
	}
	elapsed := time.Since(start)
	fmt.Println("================== Test tempai")
	fmt.Printf("Repeat: %v\n", repeat)
	fmt.Printf("Elapsed: %v\n", elapsed)
	fmt.Printf("Estemated speed: %v per second\n", float64(repeat)/elapsed.Seconds())
	fmt.Printf("Tempai hand count: %v\n", cnt)
}

func testEffective() {
	repeat := (*fCount) / 10
	data := make([]compact.Instances, repeat)
	// prepare
	source := rand.NewSource(123)
	rnd := rand.New(source)
	instances := compact.AllInstances().Instances()
	for k := range data {
		shuffle(rnd, instances)
		data[k] = compact.NewInstances().Add(instances[:14])
	}
	start := time.Now()
	for _, v := range data {
		effective.Calculate(v)

	}
	elapsed := time.Since(start)
	fmt.Println("================== Test effectivity")
	fmt.Printf("Repeat: %v\n", repeat)
	fmt.Printf("Elapsed: %v\n", elapsed)
	fmt.Printf("Estemated speed: %v per second\n", float64(repeat)/elapsed.Seconds())
}

func shuffle(r *rand.Rand, x sort.Interface) {
	r.Shuffle(x.Len(), func(i, j int) { x.Swap(i, j) })
}
