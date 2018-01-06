package meld

type Type int

//go:generate stringer -type=Type
const (
	TypeSeq Type = iota + 1
	TypeSame
	TypePair
)
