package yaku

// Different prefs for showing to user

type WindMergeType int

const (
	MergeWindNoMerge WindMergeType = iota
	MergeWindToYakuhai
	MergeWindToTile
	MergeWindToType // Self/Round
)

type YakuPrefs struct {
	MergeDragons   bool
	MergeWinds     WindMergeType
	MergeWindRound bool
	MergeDora      bool
}

var DefaultYakuPref = &YakuPrefs{}
var CompactYakuPref = &YakuPrefs{MergeDragons: true, MergeWinds: MergeWindToYakuhai, MergeDora: true}
var CompactYakuPref2 = &YakuPrefs{MergeDragons: true, MergeWinds: MergeWindToTile, MergeDora: true}
var CompactYakuPref3 = &YakuPrefs{MergeDragons: true, MergeWinds: MergeWindToType, MergeDora: true}

func (this *YakuPrefs) FormatYaku(in YakuSet) YakuSet {
	out := YakuSet{}

	for k, v := range in {
		switch k {
		case YakuHaku, YakuHatsu, YakuChun:
			if this.MergeDragons {
				out[YakuYakuhai] += v
				continue
			}
		case YakuTonSelf,
			YakuNanSelf,
			YakuSjaSelf,
			YakuPeiSelf,
			YakuTonRound,
			YakuNanRound,
			YakuSjaRound,
			YakuPeiRound:
			switch this.MergeWinds {
			case MergeWindNoMerge:
			case MergeWindToYakuhai:
				out[YakuYakuhai] += v
				continue
			case MergeWindToTile:
				if k >= YakuTonRound {
					k -= 4
				}
				k -= 4
				out[k] += v
				continue
			case MergeWindToType:
				if k >= YakuTonRound {
					out[YakuWindRound] += v
				} else {
					out[YakuWindSelf] += v
				}
				continue
			}
		case YakuUraDora, YakuAkaDora:
			if this.MergeDora && v != 0 {
				out[YakuDora] += v
				continue
			}
		}
		out[k] += v
	}

	return out
}
