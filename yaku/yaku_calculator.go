package yaku

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/meld"
	"github.com/dnovikoff/tempai-core/tile"
)

type FuInfo struct {
	Meld   meld.Meld
	Fu     Fu
	Points FuPoints
}

type Fus []*FuInfo

func (this YakuSet) Sum() HanPoints {
	sum := HanPoints(0)
	for _, v := range this {
		sum += HanPoints(v)
	}
	return sum
}

func (this YakuSet) String() string {
	results := make([]string, 0, len(this))
	for k, v := range this {
		results = append(results, fmt.Sprintf("%v: %v", k, v))
	}
	sort.Strings(results)
	return strings.Join(results, ", ")
}

type YakuSet map[Yaku]HanPoints
type YakumanSet map[Yakuman]int

type YakuResult struct {
	Melds   meld.Melds
	Yaku    YakuSet
	Yakuman YakumanSet
	Bonuses YakuSet
	Fus     Fus

	IsClosed bool
}

func (this YakuResult) Sum() HanPoints {
	x := this.Yaku.Sum()
	if x == 0 {
		return 0
	}
	return x + this.Bonuses.Sum()
}

func (this YakuResult) String() string {
	if len(this.Yakuman) > 0 {
		return this.Yakuman.String()
	}
	if len(this.Yaku) > 0 {
		x := this.Yaku.String()
		if len(this.Bonuses) > 0 {
			x += ", " + this.Bonuses.String()
		}
		return fmt.Sprintf("%v = %v", this.Sum(), x)
	}
	return "No yaku"
}

func (this YakuResult) SetValues(k Yaku, opened, closed HanPoints) {
	if this.Yaku[k] == 0 {
		return
	}

	if this.IsClosed {
		this.Yaku[k] = closed
	} else {
		this.Yaku[k] = opened
	}
	if this.Yaku[k] == 0 {
		delete(this.Yaku, k)
	}
}

func (this YakumanSet) String() string {
	results := make([]string, 0, len(this))
	for k, v := range this {
		str := fmt.Sprintf("%v:%v", k.String(), v)
		results = append(results, str)
	}
	sort.Strings(results)
	return strings.Join(results, ", ")
}

func (this Fus) String() string {
	parts := make([]string, 0, len(this))
	for _, v := range this {
		part := fmt.Sprintf("%v(%v)", v.Points, v.Fu)
		if !v.Meld.IsNull() {
			part += "[" + meld.DebugDescribe(v.Meld) + "]"
		}
		parts = append(parts, part)
	}
	return fmt.Sprintf("%v = %v", this.Sum(), strings.Join(parts, " + "))
}

func NewYakuResult(melds meld.Melds) (ret *YakuResult) {
	ret = &YakuResult{}
	ret.Melds = melds
	ret.Yaku = make(YakuSet, 16)
	ret.Bonuses = make(YakuSet, 16)
	ret.Yakuman = make(YakumanSet, 16)
	return ret
}

type YakuCalculator struct {
	ctx       *Context
	melds     meld.Melds
	seq       meld.Melds
	same      meld.Melds
	result    *YakuResult
	pair      meld.Pair
	win       meld.Meld
	winOpened meld.Meld
	base      tile.Tiles
	isClosed  bool
}

func (this Fus) Sum() (ret FuPoints) {
	for _, v := range this {
		ret += v.Points
	}
	return
}

func (this *YakuCalculator) findWin() meld.Meld {
	c := this.ctx.Tile.Tile()
	for _, v := range this.melds {
		if !v.IsComplete() && v.Waits().Check(c) {
			return v
		}
	}
	return 0
}

func calculateClosed(melds meld.Melds) bool {
	for _, v := range melds {
		if v.Interface().IsOpened() {
			return false
		}
	}
	return true
}

func NewYakuCalculator(ctx *Context, melds meld.Melds) *YakuCalculator {
	ret := &YakuCalculator{}
	ret.ctx = ctx
	ret.melds = melds
	ret.win = ret.findWin()
	ret.winOpened = ret.win.Interface().Open(ctx.Tile, base.Left)
	ret.isClosed = calculateClosed(melds)
	if ret.win.IsTanki() {
		ret.pair = meld.Pair(ret.win)
	}
	for _, v := range melds {
		ret.base = append(ret.base, v.Base())
		if v == ret.win {
			continue
		}
		switch v.Type() {
		case meld.TypeSame:
			ret.same = append(ret.same, v)
		case meld.TypeSeq:
			ret.seq = append(ret.seq, v)
		case meld.TypePair:
			ret.pair = meld.Pair(v)
		}
	}
	switch ret.win.Type() {
	case meld.TypeSame:
		ret.same = append(ret.same, ret.winOpened)
	case meld.TypeSeq:
		ret.seq = append(ret.seq, ret.winOpened)
	}

	return ret
}

func (this *YakuCalculator) finalMeld(m meld.Meld) meld.Meld {
	if m == this.win {
		return this.winOpened
	}
	return m
}

func (this *YakuCalculator) checkClosed(m meld.Meld) bool {
	if m == this.win || m == this.winOpened {
		return this.ctx.IsTsumo
	}
	return !m.Interface().IsOpened()
}

func (this *YakuCalculator) addYaku(yaku Yaku) {
	this.result.Yaku[yaku]++
}

func (this *YakuCalculator) addBonus(yaku Yaku) {
	this.result.Bonuses[yaku]++
}

func (this *YakuCalculator) addFu(fu Fu, count FuPoints) {
	this.addFuMeld(0, fu, count)
}

func (this *YakuCalculator) addFuMeld(meld meld.Meld, fu Fu, count FuPoints) {
	this.result.Fus = append(this.result.Fus, &FuInfo{meld, fu, count})
}

func (this *YakuCalculator) addYakuman2(yaku Yakuman) {
	this.result.Yakuman[yaku] = 2
}

func (this *YakuCalculator) addYakuman(yaku Yakuman) {
	this.result.Yakuman[yaku] = 1
}

func (this *YakuCalculator) isPinfu() bool {
	for _, v := range this.melds {
		i := v.Interface()
		if i.IsBadWait() ||
			i.IsOpened() ||
			v.Type() == meld.TypeSame {
			return false
		}
		base := i.Base()
		if this.ctx.SelfWind.CheckTile(base) ||
			this.ctx.RoundWind.CheckTile(base) ||
			v.Base().Type() == tile.TypeDragon {
			return false
		}
	}
	return true
}

func (this *YakuCalculator) tryFat() {
	count := 0
	closed := 0
	kan := 0

	numbers := make(map[int]int)
	for _, v := range this.same {
		i := v.Interface()
		start := i.Base()
		if start.IsSequence() {
			numbers[start.NumberInSequence()]++
		}
		count++
		if this.checkClosed(v) {
			closed++
		}
		if v.IsKan() {
			kan++
		}
	}

	if kan == 4 {
		this.addYakuman(YakumanSuukantsu)
	}
	if closed == 4 {
		if this.win.IsTanki() {
			this.addYakuman2(YakumanSuuankouTanki)
		} else {
			this.addYakuman(YakumanSuuankou)
		}
	}
	if kan > 2 {
		this.addYaku(YakuSankantsu)
	}
	if closed > 2 {
		this.addYaku(YakuSanankou)
	}
	if count == 4 {
		this.addYaku(YakuToitoi)
	}
	for _, n := range numbers {
		if n == 3 {
			this.addYaku(YakuSanshokuDoukou)
			// Just one possible
			break
		}
	}
}

func (this *YakuCalculator) tryXpeko() {
	if !this.isClosed {
		return
	}
	tst := make(map[tile.Tile]int)
	for _, v := range this.seq {
		tst[v.Interface().Base()]++
	}
	cnt := 0
	for _, v := range tst {
		if v > 1 {
			cnt++
		}
	}
	switch cnt {
	case 1:
		this.addYaku(YakuIppeiko)
	case 2:
		this.addYaku(YakuRyanpeikou)
	}
}

func (this *YakuCalculator) tryYakuhai() {
	dragons := 0
	winds := 0
	const ponValue = 10
	const pairValue = 1

	if this.pair != 0 {
		switch this.pair.Base().Type() {
		case tile.TypeDragon:
			dragons += pairValue
		case tile.TypeWind:
			winds += pairValue
		}
	}

	for _, v := range this.same {
		start := v.Interface().Base()

		switch start.Type() {
		case tile.TypeDragon:
			this.addYaku(YakuHaku + Yaku(start.NumberInSequence()-1))
			dragons += ponValue
		case tile.TypeWind:
			winds += ponValue
			if this.ctx.SelfWind.CheckTile(start) {
				this.addYaku(YakuTonSelf + Yaku(start.NumberInSequence()-1))
			}
			if this.ctx.RoundWind.CheckTile(start) {
				this.addYaku(YakuTonRound + Yaku(start.NumberInSequence()-1))
			}
		}
	}

	switch {
	case dragons == ponValue*3:
		this.addYakuman(YakumanDaisangen)
	case dragons == ponValue*2+pairValue:
		this.addYaku(YakuShousangen)
	case winds == ponValue*4:
		this.addYakuman2(YakumanDaisuushi)
	case winds == ponValue*3+pairValue:
		this.addYakuman(YakumanShousuushi)
	}
}

func (this *YakuCalculator) tryColor() Yaku {
	haveHonor := false
	color := tile.TypeWind
	for _, v := range this.base {
		if v.IsHonor() {
			haveHonor = true
			continue
		}
		if color != tile.TypeWind && color != v.Type() {
			return YakuNone
		}
		color = v.Type()
	}
	y := YakuChinitsu
	if haveHonor {
		y = YakuHonitsu
	}
	this.addYaku(y)
	return y
}

func (this *YakuCalculator) tryTileTypes() {
	const (
		maskTerminal = 1 << iota
		maskMiddle
		maskHonor
		maskChi
	)
	tmp := 0
	same := func(m meld.Meld) {
		t := m.Base()
		switch {
		case t.IsHonor():
			tmp |= maskHonor
		case t.IsTerminal():
			tmp |= maskTerminal
		default:
			tmp |= maskMiddle
		}
	}
	for _, v := range this.melds {
		switch v.Type() {
		case meld.TypeSame, meld.TypePair:
			same(v)
		case meld.TypeSeq:
			t := this.finalMeld(v).Base()
			tmp |= maskChi
			if t.NumberInSequence() == 1 || t.NumberInSequence() == 7 {
				tmp |= maskTerminal
			} else {
				tmp |= maskMiddle
			}
		default:
			return
		}
	}

	switch tmp {
	case maskTerminal | maskChi:
		this.addYaku(YakuJunchan)
	case maskTerminal:
		this.addYakuman(YakumanChinrouto)
	case maskHonor:
		this.addYakuman(YakumanTsuiisou)
	case maskTerminal | maskHonor:
		this.addYaku(YakuHonrouto)
	case maskTerminal | maskHonor | maskChi:
		this.addYaku(YakuChanta)
	case maskMiddle, maskMiddle | maskChi:
		this.addYaku(YakuTanyao)
	}
}

func (this *YakuCalculator) trySeq() {
	forSanshoku := make(map[int]int)
	forItsu := make(map[tile.Type]int)

	for _, v := range this.seq {
		first := v.Base()
		typ := first.Type()
		number := first.NumberInSequence()
		switch typ {
		case tile.TypeMan:
			forSanshoku[number] |= 1
		case tile.TypePin:
			forSanshoku[number] |= 2
		case tile.TypeSou:
			forSanshoku[number] |= 4
		}
		switch number {
		case 1:
			forItsu[typ] |= 1
		case 4:
			forItsu[typ] |= 2
		case 7:
			forItsu[typ] |= 4
		}
	}

	for _, v := range forItsu {
		if v == 7 {
			this.addYaku(YakuItsuu)
			return
		}
	}

	for _, v := range forSanshoku {
		if v == 7 {
			this.addYaku(YakuSanshoku)
			return
		}
	}
}

func (this *YakuCalculator) allInstances() tile.Instances {
	x := make(tile.Instances, 0, 14)
	for _, m := range this.melds {
		x = append(x, m.Instances()...)
	}
	x = append(x, this.ctx.Tile)
	return x
}

func (this *YakuCalculator) tryGates() {
	// already checked for color
	tst := [9]int{-3, -1, -1, -1, -1, -1, -1, -1, -3}

	winIndex := this.ctx.Tile.Tile().NumberInSequence() - 1

	addTst := func(m meld.Meld, shift, val int) {
		tst[m.Base().NumberInSequence()-1+shift] += val
	}

	for _, m := range this.same {
		addTst(m, 0, 3)
	}
	for _, m := range this.seq {
		addTst(m, 0, 1)
		addTst(m, 1, 1)
		addTst(m, 2, 1)
	}
	addTst(this.pair.Meld(), 0, 2)

	index := -1
	for k, v := range tst {
		switch v {
		case 1:
			if index > -1 {
				return
			}
			index = k
		case 0:
		default:
			return
		}
	}

	if index == winIndex {
		this.addYakuman2(YakumanChuurenpooto9)
	} else {
		this.addYakuman(YakumanChuurenpooto)
	}
}

func (this *YakuCalculator) tryDora() {
	for _, i := range this.allInstances() {
		if this.ctx.Rules.CheckAka(i) {
			this.addBonus(YakuAkaDora)
		}
		for _, dora := range this.ctx.DoraTiles {
			if i.Tile() == dora {
				this.addBonus(YakuDora)
			}
		}
		if this.ctx.ShouldAddUras() {
			for _, dora := range this.ctx.UraTiles {
				if i.Tile() == dora {
					this.addBonus(YakuUraDora)
				}
			}
		}
	}
}

func (this *YakuCalculator) tryFu() {
	if this.ctx.IsTsumo {
		this.addFu(FuTsumo, 2)
	}
	if this.win.Interface().IsBadWait() {
		this.addFuMeld(this.win, FuBadWait, 2)
	}
	if this.pair != 0 {
		base := this.pair.Base()
		tpe := base.Type()

		if tpe == tile.TypeDragon {
			this.addFuMeld(this.pair.Meld(), FuPair, 2)
		} else if tpe == tile.TypeWind {
			// Do not merge this two conditions - wind could be counted twice
			var fu FuPoints
			if this.ctx.SelfWind.CheckTile(base) {
				fu += 2
			}
			if this.ctx.RoundWind.CheckTile(base) {
				fu += 2
			}
			if fu > 0 {
				this.addFuMeld(this.pair.Meld(), FuPair, fu)
			}
		}
	}
	for _, m := range this.same {
		var base FuPoints = 2
		t := m.Base()
		if t.IsTerminalOrHonor() {
			base = 4
		}
		if this.checkClosed(m) {
			base *= 2
		}
		if m.IsKan() {
			base *= 4
		}
		this.addFuMeld(m, FuOther, base)
	}
}

func (this *YakuCalculator) calculateResult() *YakuResult {
	result := this.result
	rule := this.ctx.Rules
	result.SetValues(YakuTsumo, 0, 1)

	if !rule.OpenTanyao {
		result.SetValues(YakuTanyao, 0, 1)
	}

	result.SetValues(YakuChiitoi, 0, 2)
	result.SetValues(YakuItsuu, 1, 2)
	result.SetValues(YakuSanankou, 1, 2)
	result.SetValues(YakuHonitsu, 2, 3)
	result.SetValues(YakuChinitsu, 5, 6)
	result.SetValues(YakuChanta, 1, 2)
	result.SetValues(YakuJunchan, 2, 3)
	result.SetValues(YakuDaburi, 0, 2)
	result.SetValues(YakuHonrouto, 2, 2)
	result.SetValues(YakuToitoi, 2, 2)
	result.SetValues(YakuSanankou, 2, 2)
	result.SetValues(YakuSankantsu, 2, 2)
	result.SetValues(YakuSanshoku, 1, 2)
	result.SetValues(YakuSanshokuDoukou, 2, 2)
	result.SetValues(YakuRyanpeikou, 0, 3)
	result.SetValues(YakuIppeiko, 0, 1)
	result.SetValues(YakuShousangen, 2, 2)

	if result.Yaku[YakuChiitoi] > 0 {
		result.Fus = Fus{&FuInfo{0, FuBase7, 25}}
	} else if result.Yaku[YakuPinfu] > 0 {
		if this.ctx.IsTsumo {
			result.Fus = Fus{&FuInfo{0, FuBase, 20}}
		} else {
			result.Fus = Fus{&FuInfo{0, FuBaseClosedRon, 30}}
		}
	} else if !result.IsClosed && len(result.Fus) == 1 {
		this.addFu(FuNoOpenFu, 2)
	}

	if result.Yaku[YakuRenhou] > 0 {
		switch rule.Renhou {
		case LimitNone:
			delete(result.Yaku, YakuRenhou)
		case LimitYakuman:
			this.addYakuman(YakumanRenhou)
			delete(result.Yaku, YakuRenhou)
		default:
			hans := rule.Renhou.BaseHans()
			delete(result.Yaku, YakuRenhou)
			sum := this.result.Yaku.Sum()
			if sum < hans {
				result.Yaku = YakuSet{YakuRenhou: hans}
			}
		}
	}
	return result
}

func (this *YakuCalculator) Calculate() *YakuResult {
	this.result = NewYakuResult(this.melds)

	if len(this.melds) > 11 {
		if this.win.IsTanki() {
			this.addYakuman2(YakumanKokushi13)
		} else {
			this.addYakuman(YakumanKokushi)
		}
		return this.calculateResult()
	}

	if this.ctx.IsRinshan {
		this.addYaku(YakuRinshan)
	}

	if (!this.ctx.IsRinshan || !this.ctx.Rules.HaiteiIsFromLiveOnly) && this.ctx.IsLastTile {
		if this.ctx.IsTsumo {
			this.addYaku(YakuHaitei)
		} else if !this.ctx.Rules.HaiteiIsFromLiveOnly || !this.ctx.IsRinshan {
			this.addYaku(YakuHoutei)
		}
	}

	if this.isClosed && this.ctx.IsTsumo {
		if this.ctx.IsFirstTake {
			if this.ctx.SelfWind == base.WindEast {
				// TODO: could be mangan
				this.addYakuman(YakumanTenhou)
			} else {
				this.addYakuman(YakumanChihou)
			}
		}

		this.addYaku(YakuTsumo)
	}

	if this.ctx.IsRiichi {
		if this.ctx.ShouldAddIpatsu() {
			this.addYaku(YakuIppatsu)
		}
		if this.ctx.IsDaburi {
			this.addYaku(YakuDaburi)
		} else {
			this.addYaku(YakuRiichi)
		}
	}
	if this.ctx.IsChankan {
		this.addYaku(YakuChankan)
	}
	isChitoi := (this.isClosed && len(this.melds) == 7)
	if isChitoi {
		this.addFu(FuBase7, 25)
		this.addYaku(YakuChiitoi)
	} else {
		if this.isClosed && this.ctx.IsRon() {
			this.addFu(FuBaseClosedRon, 30)
		} else {
			this.addFu(FuBase, 20)
		}
		if this.isPinfu() {
			this.addYaku(YakuPinfu)
		} else {
			this.tryFat()
			this.tryYakuhai()
		}
		this.trySeq()
		if this.isClosed {
			this.tryXpeko()
		}
	}
	// Not listed in normal form cause of Tsuiso possible in chitoi
	this.tryTileTypes()
	colorYaku := this.tryColor()
	if !isChitoi && (colorYaku == YakuChinitsu) && this.isClosed {
		this.tryGates()
	}

	this.tryDora()
	if this.ctx.IsFirstTake && this.ctx.IsRon() {
		this.addYaku(YakuRenhou)
	}

	this.result.IsClosed = this.isClosed
	this.tryFu()

	return this.calculateResult()
}
