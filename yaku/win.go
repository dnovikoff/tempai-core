package yaku

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/calc"
	"github.com/dnovikoff/tempai-core/hand/tempai"
)

func Win(results *tempai.TempaiResults, ctx *Context, declaredTiles compact.Instances) *Result {
	isRon := ctx.isRon()
	top := 0
	var ret *Result
	for _, v := range results.Results {
		win := v.Last.Complete(ctx.Tile.Tile())
		if win == nil {
			continue
		}
		if isRon {
			win = calc.Open(win)
		}
		args := &args{
			ctx:    ctx,
			result: v,
			hand:   results.Hand,
			declared: declared{
				melds: results.Declared,
				tiles: declaredTiles,
			},
			win: win,
		}
		res := calculate(args)
		if res == nil {
			continue
		}
		if len(res.Yakuman) > 0 {
			top = 14
			return res
		}
		sum := int(res.Fus.Sum()) + int(res.Yaku.Sum()*1000)
		if sum > top && sum > 1000 {
			ret = res
			top = sum
		}
	}
	return ret
}
