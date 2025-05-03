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
			e.cancelSelection()
			e.moveInternalCursor(1, 0)
		case tcell.KeyLeft:
			e.cancelSelection()
			e.moveInternalCursor(-1, 0)
		case tcell.KeyDown:
			e.cancelSelection()
			e.moveInternalCursor(0, 1)
		case tcell.KeyUp:
			e.cancelSelection()
			e.moveInternalCursor(0, -1)
		case tcell.KeyCtrlU:
			e.moveInternalCursor(0, -e.fastJumpLength)
		case tcell.KeyCtrlD:
			e.moveInternalCursor(0, e.fastJumpLength)
		case tcell.KeyRune:
			switch ev.Rune() {
			case 'd':
				if e.Selection.Content == "" {
					e.deleteRuneAtCursor()
				} else {
					e.deleteSelection()
				}
			case 'e':
				e.moveInternalCursor(0, len(e.InternalBuffer.Data)-1)
			case 't':
				e.moveInternalCursor(0, -len(e.InternalBuffer.Data)) // the overflow is handled in another place, but hacky
			case 'h':
				e.moveInternalCursor(-1000, 0) // hacky lol
			case 'l':
				e.moveInternalCursor(1000, 0) // hacky lol
			case 'o':
				e.insertNewlineUnder()
				e.SwitchMode()
			case 'O':
				e.insertNewlineAbove()
				e.SwitchMode()
			case ':':
				e.Mode = CommandMode
			case 'i':
				e.SwitchMode()
			case 'r':
				e.replaceRuneUnder()
			case 'x':
				e.selectLine()
			case 'a':
				e.cancelSelection()
			case 'y':
				e.copySelection()
			case 'p':
				e.pasteUnder()
			case 'u':
				e.undo()
			}
		case tcell.KeyCtrlC:
			e.toggleCommentLine()
		}
	}
}
