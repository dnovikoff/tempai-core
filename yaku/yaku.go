package yaku

type Yakuman int

//go:generate stringer -type=Yakuman
// Yakuman numbers are now fixed and should not be changed
const (
	YakumanNone          Yakuman = 0
	YakumanKokushi       Yakuman = 1
	YakumanKokushi13     Yakuman = 2
	YakumanSuukantsu     Yakuman = 3
	YakumanSuuankou      Yakuman = 4
	YakumanSuuankouTanki Yakuman = 5
	YakumanDaisangen     Yakuman = 6
	YakumanShousuushi    Yakuman = 7
	YakumanDaisuushi     Yakuman = 8
	YakumanRyuuiisou     Yakuman = 9
	YakumanTsuiisou      Yakuman = 10
	YakumanChinrouto     Yakuman = 11
	YakumanChuurenpooto  Yakuman = 12
	YakumanChuurenpooto9 Yakuman = 13
	YakumanTenhou        Yakuman = 14
	YakumanChihou        Yakuman = 15
	YakumanRenhou        Yakuman = 16
)

type Yaku int

//go:generate stringer -type=Yaku
// Yaku numbers are now fixed and should not be changed
const (
	YakuNone           Yaku = 0
	YakuRiichi         Yaku = 1
	YakuDaburi         Yaku = 2
	YakuIppatsu        Yaku = 3
	YakuTsumo          Yaku = 4
	YakuTanyao         Yaku = 5
	YakuChanta         Yaku = 6
	YakuJunchan        Yaku = 7
	YakuHonrouto       Yaku = 8
	YakuYakuhai        Yaku = 9
	YakuHaku           Yaku = 10
	YakuHatsu          Yaku = 11
	YakuChun           Yaku = 12
	YakuWindRound      Yaku = 13
	YakuWindSelf       Yaku = 14
	YakuTon            Yaku = 15
	YakuNan            Yaku = 16
	YakuSja            Yaku = 17
	YakuPei            Yaku = 18
	YakuTonSelf        Yaku = 19
	YakuNanSelf        Yaku = 20
	YakuSjaSelf        Yaku = 21
	YakuPeiSelf        Yaku = 22
	YakuTonRound       Yaku = 23
	YakuNanRound       Yaku = 24
	YakuSjaRound       Yaku = 25
	YakuPeiRound       Yaku = 26
	YakuChiitoi        Yaku = 27
	YakuToitoi         Yaku = 28
	YakuSanankou       Yaku = 29
	YakuSankantsu      Yaku = 30
	YakuSanshoku       Yaku = 31
	YakuShousangen     Yaku = 32
	YakuPinfu          Yaku = 33
	YakuIppeiko        Yaku = 34
	YakuRyanpeikou     Yaku = 35
	YakuItsuu          Yaku = 36
	YakuSanshokuDoukou Yaku = 37
	YakuHonitsu        Yaku = 38
	YakuChinitsu       Yaku = 39
	YakuDora           Yaku = 40
	YakuUraDora        Yaku = 41
	YakuAkaDora        Yaku = 42
	YakuRenhou         Yaku = 43
	YakuHaitei         Yaku = 44
	YakuHoutei         Yaku = 45
	YakuRinshan        Yaku = 46
	YakuChankan        Yaku = 47
)

type Fu int

//go:generate stringer -type=Fu
const (
	FuNone Fu = iota
	FuBase
	FuBaseClosedRon
	FuBase7
	FuSet
	FuTsumo
	FuOther
	FuNoOpenFu
	FuBadWait
	FuPair
)
