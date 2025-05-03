package editor

import (
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) editModeRoutine() {
	ev := e.Screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEsc:
			if e.autoSaveOnSwitch {
				e.SaveToFile()
			}
			e.SwitchMode()
		case tcell.KeyRight:
			e.moveInternalCursor(1, 0)
		case tcell.KeyLeft:
			e.moveInternalCursor(-1, 0)
		case tcell.KeyDown:
			e.moveInternalCursor(0, 1)
		case tcell.KeyUp:
			e.moveInternalCursor(0, -1)
		case tcell.KeyRune:
			e.insertRune(ev.Rune())
		case tcell.KeyEnter:
			e.insertNewlineAtCursor()
		case tcell.KeyBackspace2:
			e.deleteRuneBeforeCursor()
		case tcell.KeyTab:
			e.insertRune('\t')
		}
	}
}
