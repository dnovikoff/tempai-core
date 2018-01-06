package yaku

type Yakuman int

//go:generate stringer -type=Yakuman
const (
	YakumanKokushi Yakuman = iota
	YakumanKokushi13
	YakumanSuukantsu
	YakumanSuuankou
	YakumanSuuankouTanki
	YakumanDaisangen
	YakumanShousuushi
	YakumanDaisuushi
	YakumanRyuuiisou
	YakumanTsuiisou
	YakumanChinrouto
	YakumanChuurenpooto
	YakumanChuurenpooto9
	YakumanTenhou
	YakumanChihou
	YakumanRenhou
)

type Yaku int

//go:generate stringer -type=Yaku
const (
	YakuNone Yaku = iota
	YakuRiichi
	YakuDaburi
	YakuIppatsu
	YakuTsumo

	YakuTanyao
	YakuChanta
	YakuJunchan
	YakuHonrouto

	// yakuhai and variants
	YakuYakuhai
	YakuHaku
	YakuHatsu
	YakuChun
	YakuWindRound
	YakuWindSelf

	// Do not change wind order
	YakuTon
	YakuNan
	YakuSja
	YakuPei
	YakuTonSelf
	YakuNanSelf
	YakuSjaSelf
	YakuPeiSelf
	YakuTonRound
	YakuNanRound
	YakuSjaRound
	YakuPeiRound

	YakuChiitoi
	YakuToitoi
	YakuSanankou
	YakuSankantsu
	YakuSanshoku
	YakuShousangen
	YakuPinfu
	YakuIppeiko
	YakuRyanpeikou
	YakuItsuu
	YakuSanshokuDoukou
	YakuHonitsu
	YakuChinitsu

	YakuDora
	YakuUraDora
	YakuAkaDora

	YakuRenhou
	YakuHaitei
	YakuHoutei
	YakuRinshan

	YakuChankan
)

type Fu int

//go:generate stringer -type=Fu
const (
	FuBase Fu = iota
	FuBaseClosedRon
	FuBase7
	FuSet
	FuTsumo
	FuOther
	FuNoOpenFu
	FuBadWait
	FuPair
)
