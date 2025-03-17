package editor

import (
	"github.com/eze-kiel/tide/file"
	"github.com/eze-kiel/tide/str"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) VisualModeRoutine(lines []string) {
	ev := e.Screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyRune:
			switch ev.Rune() {
			case 'q':
				e.Quit()
			case 'i':
				e.SwitchMode()
			case 'l':
				e.CursorX = len(lines[e.CursorY])
			case 'h':
				e.CursorX = 0
			case 'd':
				if e.CursorY+e.fastJumpLength <= len(lines) {
					e.CursorY += e.fastJumpLength
				} else {
					e.CursorY = len(lines)
				}
			case 'u':
				if e.CursorY-e.fastJumpLength > 0 {
					e.CursorY -= e.fastJumpLength
				} else {
					e.CursorY = 0
				}
			case 'w':
				if err := file.Write(e.Filename, e.Buffer); err != nil {
					e.StatusMsg = "Error: " + err.Error()
				} else {
					e.StatusMsg = str.AutoSavedMsg + e.Filename
				}
				e.StatusTimeout = 5
			}
		case tcell.KeyUp:
			e.CursorY--
		case tcell.KeyDown:
			e.CursorY++
		case tcell.KeyLeft:
			e.CursorX--
		case tcell.KeyRight:
			e.CursorX++
		}
	}
}
