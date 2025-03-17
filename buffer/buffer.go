package buffer

import "strings"

func SplitLines(text string) []string {
	return strings.Split(text, "\n")
}

func GetIndexFromPosition(text string, x, y int) int {
	lines := SplitLines(text)
	if y >= len(lines) {
		return len(text)
	}

	index := 0
	for i := 0; i < y; i++ {
		index += len(lines[i]) + 1 // +1 for the newline character
	}

	if y > 0 && index > len(text) {
		index = len(text)
	}

	if x > len(lines[y]) {
		index += len(lines[y])
	} else {
		index += x
	}

	return index
}

func InsertRune(text string, x, y int, r rune) string {
	index := GetIndexFromPosition(text, x, y)
	return text[:index] + string(r) + text[index:]
}

func RemoveRune(text string, x, y int) string {
	index := GetIndexFromPosition(text, x, y)
	if index > 0 {
		return text[:index-1] + text[index:]
	}
	return text
}

func InsertNewline(text string, x, y int) string {
	index := GetIndexFromPosition(text, x, y)
	return text[:index] + "\n" + text[index:]
}

func RemoveNewline(text string, lineIndex int) string {
	lines := SplitLines(text)
	if lineIndex <= 0 || lineIndex >= len(lines) {
		return text
	}

	result := ""
	for i, line := range lines {
		result += line
		if i < len(lines)-1 && i != lineIndex-1 {
			result += "\n"
		}
	}
	return result
}

func RemoveLine(text string, lineIndex int) string {
	lines := SplitLines(text)
	if lineIndex < 0 || lineIndex >= len(lines) {
		return text
	}

	result := ""
	for i, line := range lines {
		if lineIndex == i {
			continue
		}
		result += line
		result += "\n"
	}
	return result
}
