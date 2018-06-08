package calc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dnovikoff/tempai-core/compact"
)

type testResults struct {
	records []string
}

func (r *testResults) Record(d *ResultData) {
	if !d.Validator.Validate(d.Closed) {
		return
	}
	if len(d.Left) > 2 {
		return
	}
	s := fmt.Sprintf("%v %v %v", DebugMeld(d.Pair), DebugMelds(d.Closed.Clone()), d.Left)
	r.records = append(r.records, s)
}

func TestCalculatorTempai(t *testing.T) {
	tg := compact.NewTestGenerator(t)
	hand := tg.CompactFromString("4555m123789s333z")
	tr := &testResults{}
	opts := GetOptions(SetResults(tr))
	fm := FilterMelds(hand, CreateComplete().Clone())
	Calculate(fm, hand, opts)
	assert.Equal(t, []string{
		"nil [123s 789s 555m 333z] 4m",
		"55m [123s 789s 333z] 45m",
		"33z [123s 789s 555m] 4m3z",
	}, tr.records)
}

func TestCalculatorShanten(t *testing.T) {
	tg := compact.NewTestGenerator(t)
	hand := tg.CompactFromString("1123m")
	tr := &testResults{}
	opts := GetOptions(SetResults(tr), Opened(3))
	fm := FilterMelds(hand, CreateAll())
	Calculate(fm, hand, opts)
	assert.Equal(t, []string{
		"nil [12m (3m)] 13m",
		"nil [13m (2m)] 12m",
		"nil [23m (14m)] 11m",
		"nil [11m (1m)] 23m",
		"nil [123m] 1m",
		"11m [23m (14m)] ",
	}, tr.records)
}
