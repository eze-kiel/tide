package buffer

import (
	"strings"
)

const (
	TAB_SIZE   = 4
	TAB_SYMBOL = " "
)

type Buffer struct {
	Data string
}

func (b Buffer) SplitLines() []string {
	return strings.Split(b.Data, "\n")
}

// this function translates the internal buffer to something human-readable (\t
// are multiple spaces, ...)
func (b Buffer) Render() Buffer {
	tmp := strings.ReplaceAll(b.Data, "\t", strings.Repeat(TAB_SYMBOL, TAB_SIZE))
	return Buffer{Data: tmp}
}
