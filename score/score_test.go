package score

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/yaku"
)

func testScoreBase(r Rules, han yaku.HanPoints, fu yaku.FuPoints, honba Honba) (ret string) {
	score := GetScore(r, han, fu, honba)
	ret = fmt.Sprintf("ron = %v/%v, tsumo = %v/%v", score.PayRon, score.PayRonDealer, score.PayTsumo, score.PayTsumoDealer)
	if score.Special != yaku.LimitNone {
		ret += " " + score.Special.String()
	}
	return
}

func testScore(han yaku.HanPoints, fu yaku.FuPoints, honba Honba) string {
	return testScoreBase(RulesTenhou(), han, fu, honba)
}

func TestAlltestScore(t *testing.T) {
	tst := testScore
	assert.Equal(t, "ron = 1000/1500, tsumo = 300/500", tst(1, 30, 0))
	assert.Equal(t, "ron = 1300/2000, tsumo = 400/700", tst(1, 40, 0))
	assert.Equal(t, "ron = 1600/2400, tsumo = 400/800", tst(1, 50, 0))
	assert.Equal(t, "ron = 2000/2900, tsumo = 500/1000", tst(1, 60, 0))
	assert.Equal(t, "ron = 2300/3400, tsumo = 600/1200", tst(1, 70, 0))
	assert.Equal(t, "ron = 2600/3900, tsumo = 700/1300", tst(1, 80, 0))
	assert.Equal(t, "ron = 2900/4400, tsumo = 800/1500", tst(1, 90, 0))
	assert.Equal(t, "ron = 3200/4800, tsumo = 800/1600", tst(1, 100, 0))
	assert.Equal(t, "ron = 3600/5300, tsumo = 900/1800", tst(1, 110, 0))
	assert.Equal(t, "ron = 1300/2000, tsumo = 400/700", tst(2, 20, 0))
	assert.Equal(t, "ron = 1600/2400, tsumo = 400/800", tst(2, 25, 0))
	assert.Equal(t, "ron = 2000/2900, tsumo = 500/1000", tst(2, 26, 0))
	assert.Equal(t, "ron = 2000/2900, tsumo = 500/1000", tst(2, 30, 0))
	assert.Equal(t, "ron = 2600/3900, tsumo = 700/1300", tst(2, 40, 0))
	assert.Equal(t, "ron = 3200/4800, tsumo = 800/1600", tst(2, 50, 0))
	assert.Equal(t, "ron = 3900/5800, tsumo = 1000/2000", tst(2, 60, 0))
	assert.Equal(t, "ron = 4500/6800, tsumo = 1200/2300", tst(2, 70, 0))
	assert.Equal(t, "ron = 5200/7700, tsumo = 1300/2600", tst(2, 80, 0))
	assert.Equal(t, "ron = 5800/8700, tsumo = 1500/2900", tst(2, 90, 0))
	assert.Equal(t, "ron = 6400/9600, tsumo = 1600/3200", tst(2, 100, 0))
	assert.Equal(t, "ron = 7100/10600, tsumo = 1800/3600", tst(2, 110, 0))
	assert.Equal(t, "ron = 3200/4800, tsumo = 800/1600", tst(3, 25, 0))
	assert.Equal(t, "ron = 2600/3900, tsumo = 700/1300", tst(3, 20, 0))
	assert.Equal(t, "ron = 3900/5800, tsumo = 1000/2000", tst(3, 30, 0))
	assert.Equal(t, "ron = 5200/7700, tsumo = 1300/2600", tst(3, 40, 0))
	assert.Equal(t, "ron = 6400/9600, tsumo = 1600/3200", tst(3, 50, 0))
	assert.Equal(t, "ron = 7700/11600, tsumo = 2000/3900", tst(3, 60, 0))
	assert.Equal(t, "ron = 8000/12000, tsumo = 2000/4000 LimitMangan", tst(3, 70, 0))
	assert.Equal(t, "ron = 6400/9600, tsumo = 1600/3200", tst(4, 25, 0))
	assert.Equal(t, "ron = 5200/7700, tsumo = 1300/2600", tst(4, 20, 0))
	assert.Equal(t, "ron = 7700/11600, tsumo = 2000/3900", tst(4, 30, 0))
	assert.Equal(t, "ron = 8000/12000, tsumo = 2000/4000 LimitMangan", tst(4, 40, 0))
	assert.Equal(t, "ron = 8000/12000, tsumo = 2000/4000 LimitMangan", tst(5, 0, 0))
	assert.Equal(t, "ron = 12000/18000, tsumo = 3000/6000 LimitHaneman", tst(6, 0, 0))
	assert.Equal(t, "ron = 12000/18000, tsumo = 3000/6000 LimitHaneman", tst(7, 0, 0))
	assert.Equal(t, "ron = 16000/24000, tsumo = 4000/8000 LimitBaiman", tst(8, 0, 0))
	assert.Equal(t, "ron = 16000/24000, tsumo = 4000/8000 LimitBaiman", tst(9, 0, 0))
	assert.Equal(t, "ron = 16000/24000, tsumo = 4000/8000 LimitBaiman", tst(10, 0, 0))
	assert.Equal(t, "ron = 24000/36000, tsumo = 6000/12000 LimitSanbaiman", tst(11, 0, 0))
	assert.Equal(t, "ron = 24000/36000, tsumo = 6000/12000 LimitSanbaiman", tst(12, 0, 0))
	assert.Equal(t, "ron = 32000/48000, tsumo = 8000/16000 LimitYakuman", tst(13, 0, 0))
}

func TestAlltestScoreOtherRules(t *testing.T) {
	r := RulesTenhou()
	r.IsManganRound = true
	r.IsKazoeYakuman = false

	tst := func(han yaku.HanPoints, fu yaku.FuPoints, honba Honba) string {
		return testScoreBase(r, han, fu, honba)
	}
	assert.Equal(t, "ron = 1000/1500, tsumo = 300/500", tst(1, 30, 0))
	assert.Equal(t, "ron = 1300/2000, tsumo = 400/700", tst(1, 40, 0))
	assert.Equal(t, "ron = 1600/2400, tsumo = 400/800", tst(1, 50, 0))
	assert.Equal(t, "ron = 2000/2900, tsumo = 500/1000", tst(1, 60, 0))
	assert.Equal(t, "ron = 2300/3400, tsumo = 600/1200", tst(1, 70, 0))
	assert.Equal(t, "ron = 2600/3900, tsumo = 700/1300", tst(1, 80, 0))
	assert.Equal(t, "ron = 2900/4400, tsumo = 800/1500", tst(1, 90, 0))
	assert.Equal(t, "ron = 3200/4800, tsumo = 800/1600", tst(1, 100, 0))
	assert.Equal(t, "ron = 3600/5300, tsumo = 900/1800", tst(1, 110, 0))
	assert.Equal(t, "ron = 1600/2400, tsumo = 400/800", tst(2, 25, 0))
	assert.Equal(t, "ron = 1300/2000, tsumo = 400/700", tst(2, 20, 0))
	assert.Equal(t, "ron = 2000/2900, tsumo = 500/1000", tst(2, 30, 0))
	assert.Equal(t, "ron = 2600/3900, tsumo = 700/1300", tst(2, 40, 0))
	assert.Equal(t, "ron = 3200/4800, tsumo = 800/1600", tst(2, 50, 0))
	assert.Equal(t, "ron = 3900/5800, tsumo = 1000/2000", tst(2, 60, 0))
	assert.Equal(t, "ron = 4500/6800, tsumo = 1200/2300", tst(2, 70, 0))
	assert.Equal(t, "ron = 5200/7700, tsumo = 1300/2600", tst(2, 80, 0))
	assert.Equal(t, "ron = 5800/8700, tsumo = 1500/2900", tst(2, 90, 0))
	assert.Equal(t, "ron = 6400/9600, tsumo = 1600/3200", tst(2, 100, 0))
	assert.Equal(t, "ron = 7100/10600, tsumo = 1800/3600", tst(2, 110, 0))
	assert.Equal(t, "ron = 3200/4800, tsumo = 800/1600", tst(3, 25, 0))
	assert.Equal(t, "ron = 2600/3900, tsumo = 700/1300", tst(3, 20, 0))
	assert.Equal(t, "ron = 3900/5800, tsumo = 1000/2000", tst(3, 30, 0))
	assert.Equal(t, "ron = 5200/7700, tsumo = 1300/2600", tst(3, 40, 0))
	assert.Equal(t, "ron = 6400/9600, tsumo = 1600/3200", tst(3, 50, 0))
	assert.Equal(t, "ron = 8000/12000, tsumo = 2000/4000 LimitMangan", tst(3, 60, 0))
	assert.Equal(t, "ron = 8000/12000, tsumo = 2000/4000 LimitMangan", tst(3, 70, 0))
	assert.Equal(t, "ron = 6400/9600, tsumo = 1600/3200", tst(4, 25, 0))
	assert.Equal(t, "ron = 5200/7700, tsumo = 1300/2600", tst(4, 20, 0))
	assert.Equal(t, "ron = 8000/12000, tsumo = 2000/4000 LimitMangan", tst(4, 30, 0))
	assert.Equal(t, "ron = 8000/12000, tsumo = 2000/4000 LimitMangan", tst(4, 40, 0))
	assert.Equal(t, "ron = 8000/12000, tsumo = 2000/4000 LimitMangan", tst(5, 0, 0))
	assert.Equal(t, "ron = 12000/18000, tsumo = 3000/6000 LimitHaneman", tst(6, 0, 0))
	assert.Equal(t, "ron = 12000/18000, tsumo = 3000/6000 LimitHaneman", tst(7, 0, 0))
	assert.Equal(t, "ron = 16000/24000, tsumo = 4000/8000 LimitBaiman", tst(8, 0, 0))
	assert.Equal(t, "ron = 16000/24000, tsumo = 4000/8000 LimitBaiman", tst(9, 0, 0))
	assert.Equal(t, "ron = 16000/24000, tsumo = 4000/8000 LimitBaiman", tst(10, 0, 0))
	assert.Equal(t, "ron = 24000/36000, tsumo = 6000/12000 LimitSanbaiman", tst(11, 0, 0))
	assert.Equal(t, "ron = 24000/36000, tsumo = 6000/12000 LimitSanbaiman", tst(12, 0, 0))
	assert.Equal(t, "ron = 24000/36000, tsumo = 6000/12000 LimitSanbaiman", tst(13, 0, 0))
}

func TestRoundCalculated(t *testing.T) {
	assert.Equal(t, "ron = 1000/1500, tsumo = 300/500", testScore(1, 28, 0))
	assert.Equal(t, "ron = 1300/2000, tsumo = 400/700", testScore(1, 31, 0))
	assert.Equal(t, "ron = 1600/2400, tsumo = 400/800", testScore(1, 49, 0))
}

func TestSpecial(t *testing.T) {
	assert.Equal(t, "ron = 32000/48000, tsumo = 8000/16000 LimitYakuman", testScore(24, 0, 0))
}

func TestRenchan(t *testing.T) {
	assert.Equal(t, "ron = 1600/2100, tsumo = 500/700", testScore(1, 30, 2))
}

func TestOtherHonbaRule(t *testing.T) {
	r := RulesTenhou()
	r.HonbaValue = 500
	assert.Equal(t, "ron = 2500/3000, tsumo = 800/1000", testScoreBase(r, 1, 30, 1))
}

func TestFuRound(t *testing.T) {
	round := func(fu yaku.FuPoints) int {
		return int(fu.Round())
	}
	assert.Equal(t, 25, round(25))
	assert.Equal(t, 10, round(10))
	assert.Equal(t, 20, round(12))
	assert.Equal(t, 30, round(24))
	assert.Equal(t, 30, round(26))
	assert.Equal(t, 30, round(30))
	assert.Equal(t, 40, round(32))
}

func TestResultCounting(t *testing.T) {
	r := RulesEMA()
	result := &yaku.YakuResult{
		Yaku: yaku.YakuSet{
			yaku.YakuChanta: 1,
			yaku.YakuDora:   5,
		},
	}

	assert.Equal(t, Score{
		PayRon:         12000,
		PayRonDealer:   18000,
		PayTsumo:       3000,
		PayTsumoDealer: 6000,
		Special:        yaku.LimitHaneman,
		Han:            6,
		Fu:             0,
	}, GetScoreByResult(r, result, 0))
}

func TestResultCountingYakuman(t *testing.T) {
	r := RulesEMA()
	r.IsYakumanSum = false
	r.IsYakumanDouble = false
	result := &yaku.YakuResult{
		Yakuman: yaku.YakumanSet{
			yaku.YakumanDaisangen: 1,
			yaku.YakumanDaisuushi: 2,
			yaku.YakumanTsuiisou:  1,
		},
	}

	assert.Equal(t, Score{
		PayRon:         32000,
		PayRonDealer:   48000,
		PayTsumo:       8000,
		PayTsumoDealer: 16000,
		Special:        yaku.LimitYakuman,
	}, GetScoreByResult(r, result, 0))

	r.IsYakumanSum = true
	r.IsYakumanDouble = false
	assert.Equal(t, Score{
		PayRon:         32000 * 3,
		PayRonDealer:   48000 * 3,
		PayTsumo:       8000 * 3,
		PayTsumoDealer: 16000 * 3,
		Special:        yaku.LimitYakuman,
	}, GetScoreByResult(r, result, 0))

	r.IsYakumanSum = true
	r.IsYakumanDouble = true
	assert.Equal(t, Score{
		PayRon:         32000 * 4,
		PayRonDealer:   48000 * 4,
		PayTsumo:       8000 * 4,
		PayTsumoDealer: 16000 * 4,
		Special:        yaku.LimitYakuman,
	}, GetScoreByResult(r, result, 0))

	r.IsYakumanSum = false
	r.IsYakumanDouble = true
	assert.Equal(t, Score{
		PayRon:         32000 * 2,
		PayRonDealer:   48000 * 2,
		PayTsumo:       8000 * 2,
		PayTsumoDealer: 16000 * 2,
		Special:        yaku.LimitYakuman,
	}, GetScoreByResult(r, result, 0))
}

func TestChangesArray(t *testing.T) {
	changes := ScoreChanges{
		base.WindEast:  2,
		base.WindSouth: 4,
		base.WindWest:  6,
		base.WindNorth: 8,
	}

	assert.Equal(t, []Money{2, 4, 6, 8}, changes.ArrayFrom(base.WindEast, 4))
	assert.Equal(t, []Money{4, 6, 8, 2}, changes.ArrayFrom(base.WindSouth, 4))
	assert.Equal(t, []Money{6, 8, 2, 4}, changes.ArrayFrom(base.WindWest, 4))
	assert.Equal(t, []Money{8, 2, 4, 6}, changes.ArrayFrom(base.WindNorth, 4))
}

func TestReducedChangesArray(t *testing.T) {
	changes := ScoreChanges{
		base.WindWest: 6,
	}

	assert.Equal(t, []Money{0, 0, 6, 0}, changes.ArrayFrom(base.WindEast, 4))
	assert.Equal(t, []Money{0, 6, 0, 0}, changes.ArrayFrom(base.WindSouth, 4))
	assert.Equal(t, []Money{6, 0, 0, 0}, changes.ArrayFrom(base.WindWest, 4))
	assert.Equal(t, []Money{0, 0, 0, 6}, changes.ArrayFrom(base.WindNorth, 4))
}

func TestChangesMerge(t *testing.T) {
	assert.Equal(t, []Money{-4, 22, 48, 22},
		MergeChangeArrays(
			[]Money{2, 4, 6, 8},
			[]Money{-6, 18, 42, 14}))
}
