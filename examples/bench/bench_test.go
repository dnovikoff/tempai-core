package bench_test

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/effective"
	"github.com/dnovikoff/tempai-core/hand/shanten"
	"github.com/dnovikoff/tempai-core/hand/tempai"
	"github.com/dnovikoff/tempai-core/tile"
)

func BenchmarkShanten(b *testing.B) {
	repeat := b.N
	data := make([]compact.Instances, repeat)
	// prepare
	source := rand.NewSource(123)
	rnd := rand.New(source)
	instances := compact.AllInstances().Instances()
	for k := range data {
		shuffle(rnd, instances)
		data[k] = compact.NewInstances().Add(instances[:13])
	}
	b.ResetTimer()
	for _, v := range data {
		shanten.Calculate(v)
	}

	b.StopTimer()
	b.Log("Repeat: ", repeat)
	b.Log("RPS: ", float64(repeat)/b.Elapsed().Seconds())
}

func BenchmarkTempai(b *testing.B) {
	repeat := b.N
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

	b.ResetTimer()
	for _, v := range data {
		r := tempai.Calculate(v)
		if r != nil {
			cnt++
		}
	}

	b.StopTimer()
	b.Log("Repeat: ", repeat)
	b.Log("Tempai hand count: ", cnt)
	b.Log("RPS: ", float64(repeat)/b.Elapsed().Seconds())
}

func BenchmarkEffective(b *testing.B) {
	repeat := b.N
	data := make([]compact.Instances, repeat)
	// prepare
	source := rand.NewSource(123)
	rnd := rand.New(source)
	instances := compact.AllInstances().Instances()
	for k := range data {
		shuffle(rnd, instances)
		data[k] = compact.NewInstances().Add(instances[:14])
	}

	b.ResetTimer()
	for _, v := range data {
		effective.Calculate(v)
	}

	b.StopTimer()
	b.Log("Repeat: ", repeat)
	b.Log("RPS: ", float64(repeat)/b.Elapsed().Seconds())
}

func shuffle(r *rand.Rand, x sort.Interface) {
	r.Shuffle(x.Len(), func(i, j int) { x.Swap(i, j) })
}
