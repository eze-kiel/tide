package editor

import (
	"github.com/eze-kiel/tide/buffer"
	"github.com/eze-kiel/tide/file"
	"github.com/eze-kiel/tide/str"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) EditModeRoutine(lines []string) {
	ev := e.Screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyRune:
			e.Buffer = buffer.InsertRune(e.Buffer, e.CursorX, e.CursorY, ev.Rune())
			e.CursorX++
		case tcell.KeyEsc:
			if e.autoSaveOnSwitch {
				if err := file.Write(e.Filename, e.Buffer); err != nil {
					e.StatusMsg = "Error: " + err.Error()
				} else {
					e.StatusMsg = str.AutoSavedMsg + e.Filename
				}
				e.StatusTimeout = 5
			}
			e.SwitchMode()
		case tcell.KeyEnter:
			e.Buffer = buffer.InsertNewline(e.Buffer, e.CursorX, e.CursorY)
			e.CursorX = 0
			e.CursorY++
		case tcell.KeyBackspace2:
			if e.CursorX > 0 {
				e.Buffer = buffer.RemoveRune(e.Buffer, e.CursorX, e.CursorY)
				e.CursorX--
			} else if e.CursorY > 0 {
				lines := buffer.SplitLines(e.Buffer)
				prevLineLen := len(lines[e.CursorY-1])
				e.Buffer = buffer.RemoveNewline(e.Buffer, e.CursorY)

				e.CursorY--
				e.CursorX = prevLineLen
			}
		case tcell.KeyTab:
			e.Buffer = buffer.InsertRune(e.Buffer, e.CursorX, e.CursorY, '\t')
			e.CursorX++
		/*
			Composite keys (CTRL + stuff)
		*/
		case tcell.KeyCtrlW:
			if err := file.Write(e.Filename, e.Buffer); err != nil {
				e.StatusMsg = "Error: " + err.Error()
			} else {
				e.StatusMsg = str.SavedMsg + e.Filename
			}
			e.StatusTimeout = 5
		case tcell.KeyCtrlQ:
			e.Quit()
		case tcell.KeyCtrlX:
			e.Buffer = buffer.RemoveLine(e.Buffer, e.CursorY)
		case tcell.KeyCtrlL:
			e.CursorX = len(lines[e.CursorY])
		case tcell.KeyCtrlH:
			e.CursorX = 0
		/*
			Navigation keys
		*/
		case tcell.KeyUp:
			e.CursorY--
		case tcell.KeyDown:
			e.CursorY++
		case tcell.KeyLeft:
			e.CursorX--
		case tcell.KeyRight:
			e.CursorX++
		case tcell.KeyCtrlD:
			if e.CursorY+e.fastJumpLength <= len(lines) {
				e.CursorY += e.fastJumpLength
			} else {
				e.CursorY = len(lines)
			}
		case tcell.KeyCtrlU:
			if e.CursorY-e.fastJumpLength > 0 {
				e.CursorY -= e.fastJumpLength
			} else {
				e.CursorY = 0
			}
		}
	}
}
