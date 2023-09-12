package yaku

import (
	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/calc"
	"github.com/dnovikoff/tempai-core/hand/tempai"
	"github.com/dnovikoff/tempai-core/tile"
)

type calculator struct {
	args    *args
	helpers helpers
	result  *Result
}

type declared struct {
	melds calc.Melds
	tiles compact.Instances
}

type args struct {
	ctx      *Context
	result   *tempai.TempaiResult
	hand     compact.Instances
	declared declared
	win      calc.Meld
}

type melds struct {
	all calc.Melds
	chi calc.Melds
	pon calc.Melds
}

type helpers struct {
	melds  melds
	all    compact.Instances
	closed bool
}

func calculate(args *args) *Result {
	return newYakuCalculator(args).run()
}

func (c *calculator) initTiles() {
	c.helpers.all = c.args.hand.Clone().Merge(c.args.declared.tiles)
	c.helpers.all.Set(c.args.ctx.Tile)
}

func (c *calculator) initMelds() bool {
	all := make(calc.Melds, 0, 5)

	all = append(all, c.args.result.Melds...)
	all = append(all, c.args.declared.melds...)
	all = append(all, c.args.win)
	pair := c.args.result.Pair
	if pair != nil {
		all = append(all, pair)
	}

	chi := make(calc.Melds, 0, 5)
	pon := make(calc.Melds, 0, 5)
	for _, v := range all {
		t := v.Tags()
		if t.CheckAny(calc.TagPon) {
			pon = append(pon, v)
		} else if t.CheckAny(calc.TagChi) {
			chi = append(chi, v)
		}
	}
	c.helpers.melds.all = all
	c.helpers.melds.pon = pon
	c.helpers.melds.chi = chi
	return true
}

func (c *calculator) init(args *args) {
	c.args = args
	c.initTiles()
	c.initMelds()
	for _, v := range args.declared.melds {
		if v.Tags().CheckAny(calc.TagOpened) {
			return
		}
	}
	c.helpers.closed = true
}

func newYakuCalculator(args *args) *calculator {
	c := &calculator{}
	c.init(args)
	return c
}

func (c *calculator) run() *Result {
	ctx := c.args.ctx
	c.result = newResult()

	c.tryTsumo()
	if c.tryKokushi() {
		return c.calculateResult()
	}

	if ctx.IsRinshan {
		c.addYaku(YakuRinshan)
	}

	c.tryLastTile()
	c.tryRiichi()

	if ctx.IsChankan {
		c.addYaku(YakuChankan)
	}
	if c.isChitoi() {
		c.addFu(FuBase7, 25)
		c.addYaku(YakuChiitoi)
	} else {
		if c.helpers.closed && ctx.isRon() {
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
		if c.helpers.closed {
			c.tryXpeko()
		}
	}
	// Not listed in normal form cause of Tsuiso possible in chitoi
	c.tryTileTypes()
	c.tryColor()
	c.tryDora()
	if ctx.IsFirstTake && ctx.isRon() {
		c.addYaku(YakuRenhou)
	}

	c.result.IsClosed = c.helpers.closed
	c.tryFu()

	return c.calculateResult()
}

func (c *calculator) addYaku(yaku Yaku) {
	c.result.Yaku[yaku]++
}

func (c *calculator) addBonus(yaku Yaku) {
	c.result.Bonuses[yaku]++
}

func (c *calculator) addFu(fu Fu, count FuPoints) {
	c.addFuMeld(fu, count, nil)
}

func (c *calculator) addFuMeld(fu Fu, count FuPoints, meld calc.Meld) {
	c.result.Fus = append(c.result.Fus, &FuInfo{fu, count, meld})
}
func (c *calculator) addYakuman(y Yakuman) {
	c.result.Yakumans = append(c.result.Yakumans, y)
}

// 1. A closed hand
// 2. All chi
// 3. A non-yakuhai pair
// 4. Ryanman wait
func (c *calculator) isPinfu() bool {
	// 1. A closed hand
	if !c.helpers.closed {
		return false
	}
	// Closed kan is also declared, but no fit for pinfu anyway
	if len(c.args.declared.melds) > 0 {
		return false
	}
	res := c.args.result
	// 4. Ryanman wait
	if !res.Last.Tags().CheckAny(calc.TagRyanman) {
		return false
	}
	// 2. All chi
	for _, v := range res.Melds {
		if !v.Tags().CheckAll(calc.TagComplete | calc.TagChi) {
			return false
		}
	}
	// 3. A non-yakuhai pair
	t := res.Pair.Tile()
	ctx := c.args.ctx
	if ctx.SelfWind.CheckTile(t) ||
		ctx.RoundWind.CheckTile(t) ||
		t.Type() == tile.TypeDragon {
		return false
	}
	return true
}

func (c *calculator) tryFat() {
	count := 0
	closed := 0
	kan := 0

	numbers := make(map[int]int, 4)
	for _, v := range c.helpers.melds.pon {
		t := v.Tile()
		tags := v.Tags()
		if compact.Sequence.Check(t) {
			numbers[t.Number()]++
		}
		count++
		if !tags.CheckAny(calc.TagOpened) {
			closed++
		}
		if tags.CheckAny(calc.TagKan) {
			kan++
		}
	}
	if kan == 4 {
		c.addYakuman(YakumanSuukantsu)
	}
	if closed == 4 {
		if c.args.win.Tags().CheckAny(calc.TagPair) {
			c.addYakuman(YakumanSuuankouTanki)
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
	if !c.helpers.closed {
		return
	}
	tst := make(map[tile.Tile]int, 4)
	for _, v := range c.helpers.melds.chi {
		tst[v.Tile()]++
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
	pair := c.args.result.Pair
	if pair == nil && c.args.win.Tags().CheckAny(calc.TagPair) {
		pair = c.args.win
	}
	if pair != nil {
		switch pair.Tile().Type() {
		case tile.TypeDragon:
			dragons += pairValue
		case tile.TypeWind:
			winds += pairValue
		}
	}

	for _, v := range c.helpers.melds.pon {
		t := v.Tile()
		shift := Yaku(t.Number() - 1)
		ctx := c.args.ctx
		switch t.Type() {
		case tile.TypeDragon:
			c.addYaku(YakuHaku + shift)
			dragons += ponValue
		case tile.TypeWind:
			winds += ponValue
			if ctx.SelfWind.CheckTile(t) {
				c.addYaku(YakuTonSelf + shift)
			}
			if ctx.RoundWind.CheckTile(t) {
				c.addYaku(YakuTonRound + shift)
			}
		}
	}

	switch {
	case dragons == ponValue*3:
		c.addYakuman(YakumanDaisangen)
	case dragons == ponValue*2+pairValue:
		c.addYaku(YakuShousangen)
	case winds == ponValue*4:
		c.addYakuman(YakumanDaisuushi)
	case winds == ponValue*3+pairValue:
		c.addYakuman(YakumanShousuushi)
	}
}

func (c *calculator) tryColor() {
	var checker compact.Tiles
	var tags calc.Tags
	for _, v := range c.helpers.melds.all {
		tags |= v.Tags()
		checker = checker.Set(v.Tile().Type().Tile(1))
	}
	if (compact.Sequence & checker).Count() != 1 {
		return
	}
	clean := false
	if (compact.Honor & checker).IsEmpty() {
		c.addYaku(YakuChinitsu)
		clean = true
	} else {
		c.addYaku(YakuHonitsu)
	}
	c.tryColorYakumans(clean)
	return
}

func (c *calculator) tryColorYakumans(isClean bool) {
	if c.isChitoi() {
		return
	}
	if isClean && c.helpers.closed {
		c.tryGates()
	}
	if !isClean || !c.args.ctx.Rules.GreenRequired() {
		c.tryGreenYakuman()
	}
}

func (c *calculator) tryTileTypes() {
	const mask = calc.TagChi |
		calc.TagPon |
		calc.TagMiddle |
		calc.TagTerminal |
		calc.TagHonor

	var tags calc.Tags
	for _, v := range c.helpers.melds.all {
		tags |= v.Tags()
	}
	tags &= mask

	switch tags {
	case calc.TagChi | calc.TagPon | calc.TagTerminal:
		c.addYaku(YakuJunchan)
	case calc.TagTerminal | calc.TagPon:
		c.addYakuman(YakumanChinrouto)
	case calc.TagPon | calc.TagHonor:
		c.addYakuman(YakumanTsuiisou)
	case calc.TagPon | calc.TagHonor | calc.TagTerminal:
		c.addYaku(YakuHonrouto)
	case calc.TagPon | calc.TagChi | calc.TagHonor | calc.TagTerminal:
		c.addYaku(YakuChanta)
	default:
		if tags&(calc.TagMiddle|calc.TagTerminal|calc.TagHonor) == calc.TagMiddle {
			c.addYaku(YakuTanyao)
		}
	}
}

func (c *calculator) trySeq() {
	forSanshoku := make(map[int]int)
	forItsu := make(map[tile.Type]int)

	for _, v := range c.helpers.melds.chi {
		first := v.Tile()
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

func (c *calculator) tryGreenYakuman() {
	// TODO: optimize
	for _, v := range c.helpers.all.Instances() {
		if !compact.GreenYakuman.Check(v.Tile()) {
			return
		}
	}
	c.addYakuman(YakumanRyuuiisou)
}

func (c *calculator) tryGates() {
	if len(c.args.declared.melds) != 0 {
		return
	}
	// Already checked for color
	if c.helpers.all.UniqueTiles().Count() != 9 {
		return
	}
	tst := [9]int{3, 1, 1, 1, 1, 1, 1, 1, 3}
	var is9 bool
	if !c.helpers.all.Each(func(m compact.Mask) bool {
		index := m.Tile().Number() - 1
		diff := m.Count() - tst[index]
		if diff < 0 {
			return false
		}
		if diff == 1 {
			is9 = (m.Tile() == c.args.ctx.Tile.Tile())
			return true
		}
		return diff == 0
	}) {
		return
	}
	if is9 {
		c.addYakuman(YakumanChuurenpooto9)
	} else {
		c.addYakuman(YakumanChuurenpooto)
	}
}

func (c *calculator) tryDora() {
	// TODO: optimize
	for _, i := range c.helpers.all.Instances() {
		c.tryDoraInstance(i)
	}
}

func (c *calculator) tryDoraInstance(i tile.Instance) {
	ctx := c.args.ctx
	if ctx.Rules.CheckAka(i) {
		c.addBonus(YakuAkaDora)
	}
	t := i.Tile()
	for _, dora := range ctx.DoraTiles {
		if t == dora {
			c.addBonus(YakuDora)
		}
	}
	if !ctx.shouldAddUras() {
		return
	}
	for _, dora := range ctx.UraTiles {
		if t == dora {
			c.addBonus(YakuUraDora)
		}
	}
}

func (c *calculator) tryFu() {
	ctx := c.args.ctx
	if ctx.IsTsumo {
		if ctx.Rules.RinshanFu() || c.result.Yaku[YakuRinshan] == 0 {
			c.addFu(FuTsumo, 2)
		}
	}
	if c.args.result.Last.Tags().CheckAny(calc.TagTanki | calc.TagPenchan | calc.TagKanchan) {
		c.addFuMeld(FuBadWait, 2, c.args.result.Last)
	}
	pair := c.args.result.Pair
	if pair != nil {
		t := pair.Tile()
		tpe := t.Type()

		if tpe == tile.TypeDragon {
			c.addFuMeld(FuPair, 2, pair)
		} else if tpe == tile.TypeWind {
			// Do not merge this two conditions - wind could be counted twice
			var fu FuPoints
			if ctx.SelfWind.CheckTile(t) {
				fu += 2
			}
			if ctx.RoundWind.CheckTile(t) {
				fu += 2
			}
			if fu > 0 {
				c.addFuMeld(FuPair, fu, pair)
			}
		}
	}
	for _, m := range c.helpers.melds.pon {
		var base FuPoints = 2
		t := m.Tile()
		if compact.TerminalOrHonor.Check(t) {
			base = 4
		}
		if !m.Tags().CheckAny(calc.TagOpened) {
			base *= 2
		}
		if m.Tags().CheckAny(calc.TagKan) {
			base *= 4
		}
		c.addFuMeld(FuMeld, base, m)
	}
}

func (c *calculator) calculateResult() *Result {
	result := c.result
	rule := c.args.ctx.Rules
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
		result.Fus = Fus{&FuInfo{FuBase7, 25, nil}}
	} else if result.Yaku[YakuPinfu] > 0 {
		if c.args.ctx.IsTsumo {
			result.Fus = Fus{&FuInfo{FuBase, 20, nil}}
		} else {
			result.Fus = Fus{&FuInfo{FuBaseClosedRon, 30, nil}}
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
	// TODO: check by result type
	tags := c.args.result.Last.Tags()
	if tags.CheckAny(calc.TagKoksuhi13) {
		c.addYakuman(YakumanKokushi13)
	} else if tags.CheckAny(calc.TagKokushi) {
		c.addYakuman(YakumanKokushi)
	} else {
		return false
	}
	return true
}

func (c *calculator) tryLastTile() {
	ctx := c.args.ctx
	if !ctx.IsLastTile {
		return
	}
	if ctx.IsRinshan && ctx.Rules.HaiteiFromLiveOnly() {
		return
	}

	if ctx.IsTsumo {
		c.addYaku(YakuHaitei)
	} else {
		c.addYaku(YakuHoutei)
	}
}

func (c *calculator) tryTsumo() {
	ctx := c.args.ctx
	if !c.helpers.closed || !ctx.IsTsumo {
		return
	}
	if ctx.IsFirstTake {
		if ctx.SelfWind == base.WindEast {
			// TODO: could be mangan
			c.addYakuman(YakumanTenhou)
		} else {
			c.addYakuman(YakumanChihou)
		}
	}
	c.addYaku(YakuTsumo)
}

func (c *calculator) tryRiichi() {
	ctx := c.args.ctx
	if !ctx.IsRiichi {
		return
	}
	if ctx.shouldAddIpatsu() {
		c.addYaku(YakuIppatsu)
	}
	if ctx.IsDaburi {
		c.addYaku(YakuDaburi)
	} else {
		c.addYaku(YakuRiichi)
	}
}

func (c *calculator) isChitoi() bool {
	return c.args.result.Type == tempai.TypePairs
}
