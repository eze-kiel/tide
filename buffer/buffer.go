package buffer

import (
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
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
func (b Buffer) Translate() Buffer {
	tmp := strings.ReplaceAll(b.Data, "\t", strings.Repeat(TAB_SYMBOL, TAB_SIZE))
	return Buffer{Data: tmp}
}

// RuneLength returns the number of runes in a string
func RuneLength(s string) int {
	return utf8.RuneCountInString(s)
}

// DisplayWidth returns the display width of a string, accounting for wide characters
func DisplayWidth(s string) int {
	return runewidth.StringWidth(s)
}

// TruncateToRunes truncates a string to n runes
func TruncateToRunes(s string, n int) string {
	if n <= 0 {
		return ""
	}
	
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n])
}

// SubstringRunes returns a substring based on rune positions
func SubstringRunes(s string, start, end int) string {
	runes := []rune(s)
	if start < 0 {
		start = 0
	}
	if end > len(runes) {
		end = len(runes)
	}
	if start >= end {
		return ""
	}
	return string(runes[start:end])
}
