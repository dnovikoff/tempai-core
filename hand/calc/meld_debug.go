package calc

import (
	"strings"
)

func DebugMeld(m Meld) string {
	if m == nil {
		return "nil"
	}
	x := m.Tiles().String()
	w := m.Waits()
	if len(w) > 0 {
		x += " (" + w.String() + ")"
	}
	return x
}

func DebugMeldsString(m Melds) string {
	return strings.Join(DebugMelds(m), " ")
}

func DebugMelds(m Melds) []string {
	data := make([]string, len(m))
	for k, v := range m {
		data[k] = DebugMeld(v)
	}
	return data
}
