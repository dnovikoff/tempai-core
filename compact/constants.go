package compact

import "github.com/dnovikoff/tempai-core/tile"

const (
	Man1 Tiles = 1 << iota
	Man2
	Man3
	Man4
	Man5
	Man6
	Man7
	Man8
	Man9

	Pin1
	Pin2
	Pin3
	Pin4
	Pin5
	Pin6
	Pin7
	Pin8
	Pin9

	Sou1
	Sou2
	Sou3
	Sou4
	Sou5
	Sou6
	Sou7
	Sou8
	Sou9

	East
	South
	West
	North

	White
	Green
	Red

	TileEnd
)

const (
	Terminal = Sou1 | Sou9 | Man1 | Man9 | Pin1 | Pin9
	Wind     = East | South | West | North
	Dragon   = White | Green | Red

	Man = Man1 | Man2 | Man3 | Man4 | Man5 | Man6 | Man7 | Man8 | Man9
	Pin = Pin1 | Pin2 | Pin3 | Pin4 | Pin5 | Pin6 | Pin7 | Pin8 | Pin9
	Sou = Sou1 | Sou2 | Sou3 | Sou4 | Sou5 | Sou6 | Sou7 | Sou8 | Sou9

	Sequence = Man | Pin | Sou

	Honor           = Wind | Dragon
	TerminalOrHonor = Terminal | Honor

	AllTiles = TileEnd - 1
	Middle   = AllTiles ^ TerminalOrHonor

	GreenYakuman = Sou2 | Sou3 | Sou4 | Sou6 | Sou8 | Green
)

// 53 is js limit for int64
const _ = uint(53 - tile.TileEnd)
