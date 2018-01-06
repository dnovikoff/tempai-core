package meld

func DebugDescribe(m Meld) string {
	ret := m.Instances().Tiles().String()
	if ret == "" {
		ret = "_"
	}
	if !m.IsComplete() {
		ret += "(" + m.Waits().Tiles().String() + ")"
	} else if m.Interface().IsOpened() {
		ret += "+"
	}
	return ret
}
