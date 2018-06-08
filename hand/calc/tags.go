package calc

type Tags int

const (
	TagChi Tags = 1 << iota
	TagPon
	TagKan
	TagPair
	TagTanki
	TagKanchan
	TagPenchan
	TagRyanman
	TagComplete
	TagKokushi
	TagKoksuhi13
	TagOpened

	TagHonor
	TagTerminal
	TagMiddle
)

func (t Tags) CheckAny(x Tags) bool {
	return (t & x) != 0
}

func (t Tags) CheckAll(x Tags) bool {
	return (t & x) == x
}
