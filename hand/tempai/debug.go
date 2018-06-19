package tempai

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/dnovikoff/tempai-core/hand/calc"
)

func DebugTempai(in *TempaiResult) string {
	buf := &bytes.Buffer{}
	// fmt.Fprintf(buf, "%v", in.Type)
	// fmt.Fprintf(buf, " Tiles: %v", in.Closed.Instances())
	// if in.Declared != nil {
	// 	fmt.Fprintf(buf, " [%v]", meld.DebugMelds(in.Declared))
	// }
	if in.Pair != nil {
		fmt.Fprintf(buf, " %v", calc.DebugMeld(in.Pair))
	}
	if in.Melds != nil {
		fmt.Fprintf(buf, " %v", calc.DebugMeldsString(in.Melds))
	}
	fmt.Fprintf(buf, " %v", calc.DebugMeld(in.Last))
	return strings.TrimSpace(buf.String())
}
