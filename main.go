package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
)

const (
	editModeStr = "EDIT"
	savedMsgStr = "Saved to "
)

func main() {
	if len(os.Args) < 2 {
		panic(errors.New("missing file name"))
	}
	f := os.Args[1]
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer screen.Fini()

	var buffer string
	if fileExists(f) {
		buffer, err = readFile(f)
		if err != nil {
			panic(err)
		}
	}

	cx, cy := 0, 0
	swidth, sheight := screen.Size()
	offsetX, offsetY := 0, 0

	var statusMsg string
	statusTimeout := 0
	for {
		screen.Clear()

		lines := splitLines(buffer)
		for i := offsetY; i < len(lines) && i < offsetY+sheight-1; i++ {
			l := lines[i]
			for j := offsetX; j < len(l) && j < offsetX+swidth; j++ {
				screen.SetContent(j-offsetX, i-offsetY, rune(l[j]), nil, tcell.StyleDefault)
			}
		}

		for i, r := range editModeStr {
			screen.SetContent(i, sheight-1, r, nil, tcell.StyleDefault.Reverse(true))
		}
		if statusMsg != "" && statusTimeout > 0 {
			statusTimeout--
			for i, r := range statusMsg {
				if i < swidth {
					screen.SetContent(swidth-len(statusMsg)+i, sheight-1, r, nil, tcell.StyleDefault)
				}
			}
		}

		screen.ShowCursor(cx-offsetX, cy-offsetY)
		screen.Show()

		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyRune:
				buffer = insertRune(buffer, cx, cy, ev.Rune())
				cx++
			case tcell.KeyEnter:
				buffer = insertNewline(buffer, cx, cy)
				cx = 0
				cy++
			case tcell.KeyBackspace2:
				if cx > 0 {
					buffer = removeRune(buffer, cx, cy)
					cx--
				} else if cy > 0 {
					lines := splitLines(buffer)
					prevLineLen := len(lines[cy-1])
					buffer = removeNewline(buffer, cy)

					cy--
					cx = prevLineLen
				}
			case tcell.KeyTab:
				buffer = insertRune(buffer, cx, cy, '\t')
				cx++
			/*
				Composite keys (CTRL + stuff)
			*/
			case tcell.KeyCtrlW:
				if writeFile(f, buffer) != nil {
					statusMsg = "Error: " + err.Error()
				} else {
					statusMsg = savedMsgStr + f
				}
				statusTimeout = 5
			case tcell.KeyCtrlQ, tcell.KeyEscape:
				return
			case tcell.KeyCtrlX:
				buffer = removeLine(buffer, cy)
			case tcell.KeyCtrlL:
				cx = len(lines[cy])
			case tcell.KeyCtrlH:
				cx = 0
			/*
				Navigation keys
			*/
			case tcell.KeyUp:
				cy--
			case tcell.KeyDown:
				cy++
			case tcell.KeyLeft:
				cx--
			case tcell.KeyRight:
				cx++
			}
		}

		// keep cursor in bounds
		lines = splitLines(buffer)
		if cx < 0 {
			cx = 0
		}
		if cy < 0 {
			cy = 0
		}

		if cy >= len(lines) {
			cy = len(lines) - 1
		}

		if cx > len(lines[cy]) {
			cx = len(lines[cy])
		}

		if cx < offsetX {
			offsetX = cx
		} else if cx >= offsetX+swidth {
			offsetX = cx - swidth + 1
		}

		if cy < offsetY {
			offsetY = cy
		} else if cy >= offsetY+sheight-1 {
			offsetY = cy - sheight + 2
		}
	}
}

func fileExists(fname string) bool {
	info, err := os.Stat(fname)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func readFile(fname string) (string, error) {
	data, err := os.ReadFile(fname)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func writeFile(fname string, content string) error {
	return os.WriteFile(fname, []byte(content), 0644)
}

func splitLines(text string) []string {
	return strings.Split(text, "\n")
}

func getIndexFromPosition(text string, x, y int) int {
	lines := splitLines(text)
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

func insertRune(text string, x, y int, r rune) string {
	index := getIndexFromPosition(text, x, y)
	return text[:index] + string(r) + text[index:]
}

func removeRune(text string, x, y int) string {
	index := getIndexFromPosition(text, x, y)
	if index > 0 {
		return text[:index-1] + text[index:]
	}
	return text
}

func insertNewline(text string, x, y int) string {
	index := getIndexFromPosition(text, x, y)
	return text[:index] + "\n" + text[index:]
}

func removeNewline(text string, lineIndex int) string {
	lines := splitLines(text)
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

func removeLine(text string, lineIndex int) string {
	lines := splitLines(text)
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
