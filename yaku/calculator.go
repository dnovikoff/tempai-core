package yaku

import (
	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/meld"
	"github.com/dnovikoff/tempai-core/tile"
)

type calculator struct {
	ctx       *Context
	melds     meld.Melds
	seq       meld.Melds
	same      meld.Melds
	result    *Result
	pair      meld.Pair
	win       meld.Meld
	winOpened meld.Meld
	base      tile.Tiles
	isClosed  bool
}

func calculate(ctx *Context, melds meld.Melds) *Result {
	return newYakuCalculator(ctx, melds).run()
}

func newYakuCalculator(ctx *Context, melds meld.Melds) *calculator {
	ret := &calculator{}
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

func (c *calculator) run() *Result {
	c.result = newResult(c.melds)

	if c.tryKokushi() {
		return c.calculateResult()
	}

	if c.ctx.IsRinshan {
		c.addYaku(YakuRinshan)
	}

	c.tryLastTile()
	c.tryTsumo()
	c.tryRiichi()

	if c.ctx.IsChankan {
		c.addYaku(YakuChankan)
	}
	if c.isChitoi() {
		c.addFu(FuBase7, 25)
		c.addYaku(YakuChiitoi)
	} else {
		if c.isClosed && c.ctx.isRon() {
			c.addFu(FuBaseClosedRon, 30)
		} else {
			c.addFu(FuBase, 20)
		}
		if c.isPinfu() {
			c.addYaku(YakuPinfu)
		} else {
			c.tryFat()
			c.tryYakuhai()
		}
		c.trySeq()
		if c.isClosed {
			c.tryXpeko()
		}
	}
	// Not listed in normal form cause of Tsuiso possible in chitoi
	c.tryTileTypes()
	c.tryColor()
	c.tryDora()
	if c.ctx.IsFirstTake && c.ctx.isRon() {
		c.addYaku(YakuRenhou)
	}

	c.result.IsClosed = c.isClosed
	c.tryFu()

	return c.calculateResult()
}

func (c *calculator) findWin() meld.Meld {
	t := c.ctx.Tile.Tile()
	for _, v := range c.melds {
		if !v.IsComplete() && v.Waits().Check(t) {
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

func (c *calculator) finalMeld(m meld.Meld) meld.Meld {
	if m == c.win {
		return c.winOpened
	}
	return m
}

func (c *calculator) checkClosed(m meld.Meld) bool {
	if m == c.win || m == c.winOpened {
		return c.ctx.IsTsumo
	}
	return !m.Interface().IsOpened()
}

func (c *calculator) addYaku(yaku Yaku) {
	c.result.Yaku[yaku]++
}

func (c *calculator) addBonus(yaku Yaku) {
	c.result.Bonuses[yaku]++
}

func (c *calculator) addFu(fu Fu, count FuPoints) {
	c.addFuMeld(0, fu, count)
}

func (c *calculator) addFuMeld(meld meld.Meld, fu Fu, count FuPoints) {
	c.result.Fus = append(c.result.Fus, &FuInfo{meld, fu, count})
}

func (c *calculator) addYakuman2(yaku Yakuman) {
	c.result.Yakuman[yaku] = 2
}

func (c *calculator) addYakuman(yaku Yakuman) {
	c.result.Yakuman[yaku] = 1
}

func (c *calculator) isPinfu() bool {
	for _, v := range c.melds {
		i := v.Interface()
		if i.IsBadWait() ||
			i.IsOpened() ||
			v.Type() == meld.TypeSame {
			return false
		}
		base := i.Base()
		if c.ctx.SelfWind.CheckTile(base) ||
			c.ctx.RoundWind.CheckTile(base) ||
			v.Base().Type() == tile.TypeDragon {
			return false
		}
	}
	return true
}

func (c *calculator) tryFat() {
	count := 0
	closed := 0
	kan := 0

	numbers := make(map[int]int)
	for _, v := range c.same {
		i := v.Interface()
		start := i.Base()
		if compact.Sequence.Check(start) {
			numbers[start.Number()]++
		}
		count++
		if c.checkClosed(v) {
			closed++
		}
		if v.IsKan() {
			kan++
		}
	}

	if kan == 4 {
		c.addYakuman(YakumanSuukantsu)
	}
	if closed == 4 {
		if c.win.IsTanki() {
			c.addYakuman2(YakumanSuuankouTanki)
		} else {
			c.addYakuman(YakumanSuuankou)
		}
	}
	if kan > 2 {
		c.addYaku(YakuSankantsu)
	}
	if closed > 2 {
		c.addYaku(YakuSanankou)
	}
	if count == 4 {
		c.addYaku(YakuToitoi)
	}
	for _, n := range numbers {
		if n == 3 {
			c.addYaku(YakuSanshokuDoukou)
			// Just one possible
			break
		}
	}
}

func (c *calculator) tryXpeko() {
	if !c.isClosed {
		return
	}
	tst := make(map[tile.Tile]int)
	for _, v := range c.seq {
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
		c.addYaku(YakuIppeiko)
	case 2:
		c.addYaku(YakuRyanpeikou)
	}
}

func (c *calculator) tryYakuhai() {
	dragons := 0
	winds := 0
	const ponValue = 10
	const pairValue = 1

	if c.pair != 0 {
		switch c.pair.Base().Type() {
		case tile.TypeDragon:
			dragons += pairValue
		case tile.TypeWind:
			winds += pairValue
		}
	}

	for _, v := range c.same {
		start := v.Interface().Base()

		switch start.Type() {
		case tile.TypeDragon:
			c.addYaku(YakuHaku + Yaku(start.Number()-1))
			dragons += ponValue
		case tile.TypeWind:
			winds += ponValue
			if c.ctx.SelfWind.CheckTile(start) {
				c.addYaku(YakuTonSelf + Yaku(start.Number()-1))
			}
			if c.ctx.RoundWind.CheckTile(start) {
				c.addYaku(YakuTonRound + Yaku(start.Number()-1))
			}
		}
	}

	switch {
	case dragons == ponValue*3:
		c.addYakuman(YakumanDaisangen)
	case dragons == ponValue*2+pairValue:
		c.addYaku(YakuShousangen)
	case winds == ponValue*4:
		c.addYakuman2(YakumanDaisuushi)
	case winds == ponValue*3+pairValue:
		c.addYakuman(YakumanShousuushi)
	}
}

func (c *calculator) tryColor() {
	haveHonor := false
	color := tile.TypeWind
	for _, v := range c.base {
		if compact.Honor.Check(v) {
			haveHonor = true
			continue
		}
		if color != tile.TypeWind && color != v.Type() {
			return
		}
		color = v.Type()
	}
	y := YakuChinitsu
	if haveHonor {
		y = YakuHonitsu
	}
	c.addYaku(y)
	c.tryColorYakumans(y == YakuChinitsu)
	return
}

func (c *calculator) tryColorYakumans(isClean bool) {
	if c.isChitoi() {
		return
	}
	if isClean && c.isClosed {
		c.tryGates()
	}
	if !isClean || !c.ctx.Rules.GreenRequired() {
		c.tryGreenYakuman()
	}
}

const (
	maskTerminal = 1 << iota
	maskMiddle
	maskHonor
	maskChi
)

func maskForChi(t tile.Tile) int {
	tmp := maskChi
	if t.Number() == 1 || t.Number() == 7 {
		tmp |= maskTerminal
	} else {
		tmp |= maskMiddle
	}
	return tmp
}

func maskForSame(t tile.Tile) int {
	switch {
	case compact.Honor.Check(t):
		return maskHonor
	case compact.Terminal.Check(t):
		return maskTerminal
	}
	return maskMiddle
}

func maskForMeld(m meld.Meld) int {
	t := m.Base()
	switch m.Type() {
	case meld.TypeSame, meld.TypePair:
		return maskForSame(t)
	case meld.TypeSeq:
		return maskForChi(t)
	}
	return 0
}

func (c *calculator) tryTileTypes() {
	tmp := 0
	for _, v := range c.melds {
		tmp |= maskForMeld(c.finalMeld(v))
	}

	switch tmp {
	case maskTerminal | maskChi:
		c.addYaku(YakuJunchan)
	case maskTerminal:
		c.addYakuman(YakumanChinrouto)
	case maskHonor:
		c.addYakuman(YakumanTsuiisou)
	case maskTerminal | maskHonor:
		c.addYaku(YakuHonrouto)
	case maskTerminal | maskHonor | maskChi:
		c.addYaku(YakuChanta)
	case maskMiddle, maskMiddle | maskChi:
		c.addYaku(YakuTanyao)
	}
}

func (c *calculator) trySeq() {
	forSanshoku := make(map[int]int)
	forItsu := make(map[tile.Type]int)

	for _, v := range c.seq {
		first := v.Base()
		typ := first.Type()
		number := first.Number()
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
			c.addYaku(YakuItsuu)
			return
		}
	}

	for _, v := range forSanshoku {
		if v == 7 {
			c.addYaku(YakuSanshoku)
			return
		}
	}
}

func (c *calculator) allInstances() tile.Instances {
	x := make(tile.Instances, 0, 14)
	for _, m := range c.melds {
		x = append(x, m.Instances()...)
	}
	x = append(x, c.ctx.Tile)
	return x
}

func (c *calculator) tryGreenYakuman() {
	for _, v := range c.allInstances() {
		if !compact.GreenYakuman.Check(v.Tile()) {
			return
		}
	}
	c.addYakuman(YakumanRyuuiisou)
}

func (c *calculator) tryGates() {
	// already checked for color
	tst := [9]int{-3, -1, -1, -1, -1, -1, -1, -1, -3}

	winIndex := c.ctx.Tile.Tile().Number() - 1

	addTst := func(m meld.Meld, shift, val int) {
		tst[m.Base().Number()-1+shift] += val
	}

	for _, m := range c.same {
		addTst(m, 0, 3)
	}
	for _, m := range c.seq {
		addTst(m, 0, 1)
		addTst(m, 1, 1)
		addTst(m, 2, 1)
	}
	addTst(c.pair.Meld(), 0, 2)

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
		c.addYakuman2(YakumanChuurenpooto9)
	} else {
		c.addYakuman(YakumanChuurenpooto)
	}
}

func (c *calculator) tryDora() {
	for _, i := range c.allInstances() {
		c.tryDoraInstance(i)
	}
}

func (c *calculator) tryDoraInstance(i tile.Instance) {
	if c.ctx.Rules.CheckAka(i) {
		c.addBonus(YakuAkaDora)
	}
	t := i.Tile()
	for _, dora := range c.ctx.DoraTiles {
		if t == dora {
			c.addBonus(YakuDora)
		}
	}
	if !c.ctx.shouldAddUras() {
		return
	}
	for _, dora := range c.ctx.UraTiles {
		if t == dora {
			c.addBonus(YakuUraDora)
		}
	}
}

func (c *calculator) tryFu() {
	if c.ctx.IsTsumo {
		if c.ctx.Rules.RinshanFu() || c.result.Yaku[YakuRinshan] == 0 {
			c.addFu(FuTsumo, 2)
		}
	}
	if c.win.Interface().IsBadWait() {
		c.addFuMeld(c.win, FuBadWait, 2)
	}
	if c.pair != 0 {
		base := c.pair.Base()
		tpe := base.Type()

		if tpe == tile.TypeDragon {
			c.addFuMeld(c.pair.Meld(), FuPair, 2)
		} else if tpe == tile.TypeWind {
			// Do not merge this two conditions - wind could be counted twice
			var fu FuPoints
			if c.ctx.SelfWind.CheckTile(base) {
				fu += 2
			}
			if c.ctx.RoundWind.CheckTile(base) {
				fu += 2
			}
			if fu > 0 {
				c.addFuMeld(c.pair.Meld(), FuPair, fu)
			}
		}
	}
	for _, m := range c.same {
		var base FuPoints = 2
		t := m.Base()
		if compact.TerminalOrHonor.Check(t) {
			base = 4
		}
		if c.checkClosed(m) {
			base *= 2
		}
		if m.IsKan() {
			base *= 4
		}
		c.addFuMeld(m, FuOther, base)
	}
}

func (c *calculator) calculateResult() *Result {
	result := c.result
	rule := c.ctx.Rules
	result.setValues(YakuTsumo, 0, 1)

	if !rule.OpenTanyao() {
		result.setValues(YakuTanyao, 0, 1)
	}

	result.setValues(YakuChiitoi, 0, 2)
	result.setValues(YakuItsuu, 1, 2)
	result.setValues(YakuSanankou, 1, 2)
	result.setValues(YakuHonitsu, 2, 3)
	result.setValues(YakuChinitsu, 5, 6)
	result.setValues(YakuChanta, 1, 2)
	result.setValues(YakuJunchan, 2, 3)
	result.setValues(YakuDaburi, 0, 2)
	result.setValues(YakuHonrouto, 2, 2)
	result.setValues(YakuToitoi, 2, 2)
	result.setValues(YakuSanankou, 2, 2)
	result.setValues(YakuSankantsu, 2, 2)
	result.setValues(YakuSanshoku, 1, 2)
	result.setValues(YakuSanshokuDoukou, 2, 2)
	result.setValues(YakuRyanpeikou, 0, 3)
	result.setValues(YakuIppeiko, 0, 1)
	result.setValues(YakuShousangen, 2, 2)

	if result.Yaku[YakuChiitoi] > 0 {
		result.Fus = Fus{&FuInfo{0, FuBase7, 25}}
	} else if result.Yaku[YakuPinfu] > 0 {
		if c.ctx.IsTsumo {
			result.Fus = Fus{&FuInfo{0, FuBase, 20}}
		} else {
			result.Fus = Fus{&FuInfo{0, FuBaseClosedRon, 30}}
		}
	} else if !result.IsClosed && len(result.Fus) == 1 {
		c.addFu(FuNoOpenFu, 2)
	}

	if result.Yaku[YakuRenhou] > 0 {
		switch rule.Renhou() {
		case LimitNone:
			delete(result.Yaku, YakuRenhou)
		case LimitYakuman:
			c.addYakuman(YakumanRenhou)
			delete(result.Yaku, YakuRenhou)
		default:
			hans := rule.Renhou().BaseHans()
			delete(result.Yaku, YakuRenhou)
			sum := c.result.Yaku.Sum()
			if sum < hans {
				result.Yaku = YakuSet{YakuRenhou: hans}
			}
		}
	}
	return result
}

func (c *calculator) tryKokushi() bool {
	if len(c.melds) < 12 {
		return false
	}
	if c.win.IsTanki() {
		c.addYakuman2(YakumanKokushi13)
	} else {
		c.addYakuman(YakumanKokushi)
	}
	return true
}

func (c *calculator) tryLastTile() {
	if !c.ctx.IsLastTile {
		return
	}
	if c.ctx.IsRinshan && c.ctx.Rules.HaiteiFromLiveOnly() {
		return
	}

	if c.ctx.IsTsumo {
		c.addYaku(YakuHaitei)
	} else {
		c.addYaku(YakuHoutei)
	}
}

func (c *calculator) tryTsumo() {
	if !c.isClosed || !c.ctx.IsTsumo {
		return
	}
	if c.ctx.IsFirstTake {
		if c.ctx.SelfWind == base.WindEast {
			// TODO: could be mangan
			c.addYakuman(YakumanTenhou)
		} else {
			c.addYakuman(YakumanChihou)
		}
	}
	c.addYaku(YakuTsumo)
}

func (c *calculator) tryRiichi() {
	if !c.ctx.IsRiichi {
		return
	}
	if c.ctx.shouldAddIpatsu() {
		c.addYaku(YakuIppatsu)
	}
	if c.ctx.IsDaburi {
		c.addYaku(YakuDaburi)
	} else {
		c.addYaku(YakuRiichi)
	}
}

func (c *calculator) isChitoi() bool {
	return c.isClosed && len(c.melds) == 7
}
