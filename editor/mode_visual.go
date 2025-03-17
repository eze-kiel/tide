package editor

import (
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) visualModeRoutine() {
	ev := e.Screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyRight:
			e.moveInternalCursor(1, 0)
		case tcell.KeyLeft:
			e.moveInternalCursor(-1, 0)
		case tcell.KeyDown:
			e.moveInternalCursor(0, 1)
		case tcell.KeyUp:
			e.moveInternalCursor(0, -1)
		case tcell.KeyCtrlU:
			e.moveInternalCursor(0, -e.fastJumpLength)
		case tcell.KeyCtrlD:
			e.moveInternalCursor(0, e.fastJumpLength)
		case tcell.KeyRune:
			switch ev.Rune() {
			case ':':
				e.Mode = CommandMode
			case 'i':
				e.SwitchMode()
			}
		}
	}
}
