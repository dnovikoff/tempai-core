package yaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYakuCompact(t *testing.T) {
	yakus := YakuSet{}

	yakus[YakuRiichi] = 1
	yakus[YakuDora] = 2
	yakus[YakuAkaDora] = 3
	yakus[YakuUraDora] = 4
	yakus[YakuChun] = 1
	yakus[YakuHatsu] = 1
	yakus[YakuHaku] = 1

	yakus[YakuTonSelf] = 1
	yakus[YakuTonRound] = 1
	yakus[YakuNanSelf] = 1
	yakus[YakuNanRound] = 1
	yakus[YakuSjaRound] = 1
	yakus[YakuSjaSelf] = 1
	yakus[YakuPeiRound] = 1
	yakus[YakuPeiSelf] = 1

	assert.Equal(t, "YakuAkaDora: 3, YakuChun: 1, YakuDora: 2, YakuHaku: 1, YakuHatsu: 1, YakuNanRound: 1, YakuNanSelf: 1, YakuPeiRound: 1, YakuPeiSelf: 1, YakuRiichi: 1, YakuSjaRound: 1, YakuSjaSelf: 1, YakuTonRound: 1, YakuTonSelf: 1, YakuUraDora: 4", yakus.String())
	assert.Equal(t, "YakuDora: 9, YakuRiichi: 1, YakuYakuhai: 11", CompactYakuPref.FormatYaku(yakus).String())
	assert.Equal(t, "YakuDora: 9, YakuNan: 2, YakuPei: 2, YakuRiichi: 1, YakuSja: 2, YakuTon: 2, YakuYakuhai: 3", CompactYakuPref2.FormatYaku(yakus).String())
	assert.Equal(t, "YakuDora: 9, YakuRiichi: 1, YakuWindRound: 4, YakuWindSelf: 4, YakuYakuhai: 3", CompactYakuPref3.FormatYaku(yakus).String())
}
