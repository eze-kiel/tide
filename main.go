package main

import (
	"errors"
	"os"

	"github.com/eze-kiel/tide/buffer"
	"github.com/eze-kiel/tide/editor"
	"github.com/eze-kiel/tide/file"
	"github.com/eze-kiel/tide/str"
	"github.com/gdamore/tcell/v2"
)

func main() {
	if len(os.Args) < 2 {
		panic(errors.New("missing file name"))
	}
	f := os.Args[1]

	e, err := editor.New()
	if err != nil {
		panic(err)
	}
	defer e.Screen.Fini()

	if file.Exists(f) {
		e.Buffer, err = file.Read(f)
		if err != nil {
			panic(err)
		}
	}

	swidth, sheight := e.Screen.Size()

	for {
		e.Screen.Clear()
		lines := buffer.SplitLines(e.Buffer)
		for i := e.OffsetY; i < len(lines) && i < e.OffsetY+sheight-1; i++ {
			l := lines[i]
			for j := e.OffsetX; j < len(l) && j < e.OffsetX+swidth; j++ {
				e.Screen.SetContent(j-e.OffsetX, i-e.OffsetY, rune(l[j]), nil, tcell.StyleDefault)
			}
		}

		// the switch below could be simplified?
		switch e.Mode {
		case editor.EditMode:
			for i, r := range str.EditMode {
				e.Screen.SetContent(i, sheight-1, r, nil, tcell.StyleDefault.Reverse(true))
			}
		case editor.VisualMode:
			for i, r := range str.VisualMode {
				e.Screen.SetContent(i, sheight-1, r, nil, tcell.StyleDefault.Reverse(true))
			}
		case editor.CommandMode:
			// TODO
		}

		if e.StatusMsg != "" && e.StatusTimeout > 0 {
			e.StatusTimeout--
			for i, r := range e.StatusMsg {
				if i < swidth {
					e.Screen.SetContent(swidth-len(e.StatusMsg)+i, sheight-1, r, nil, tcell.StyleDefault)
				}
			}
		}

		e.Screen.ShowCursor(e.CursorX-e.OffsetX, e.CursorY-e.OffsetY)
		e.Screen.Show()

		switch e.Mode {
		case editor.EditMode:
			e.EditModeRoutine(f, lines)
		case editor.VisualMode:
			e.VisualModeRoutine(f, lines)
		default:
			continue
		}

		// keep the cursor in bounds
		lines = buffer.SplitLines(e.Buffer)
		if e.CursorX < 0 {
			e.CursorX = 0
		}
		if e.CursorY < 0 {
			e.CursorY = 0
		}

		if e.CursorY >= len(lines) {
			e.CursorY = len(lines) - 1
		}

		if e.CursorX > len(lines[e.CursorY]) {
			e.CursorX = len(lines[e.CursorY])
		}

		if e.CursorX < e.OffsetX {
			e.OffsetX = e.CursorX
		} else if e.CursorX >= e.OffsetX+swidth {
			e.OffsetX = e.CursorX - swidth + 1
		}

		if e.CursorY < e.OffsetY {
			e.OffsetY = e.CursorY
		} else if e.CursorY >= e.OffsetY+sheight-1 {
			e.OffsetY = e.CursorY - sheight + 2
		}
	}
}
