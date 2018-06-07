package yaku

import (
	"github.com/dnovikoff/tempai-core/hand/tempai"
	"github.com/dnovikoff/tempai-core/meld"
)

func Win(tempai tempai.IndexedResult, ctx *Context) *Result {
	current := tempai[ctx.Tile.Tile()]

	if len(current) == 0 {
		return nil
	}
	top := 0
	var res *Result
	for _, v := range current {
		waiting := append(meld.Melds{}, v...)
		winMeld := waiting.Win(ctx.Tile.Tile())
		if winMeld.IsNull() {
			return nil
		}
		result := calculate(ctx, waiting)
		if result == nil {
			return nil
		}

		if len(result.Yakuman) > 0 {
			top = 14
			return result
		}
		sum := int(result.Fus.Sum()) + int(result.Yaku.Sum()*1000)
		if sum > top && sum > 1000 {
			res = result
			top = sum
		}
	}
	return res
}
