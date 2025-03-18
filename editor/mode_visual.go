package editor

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

func (e *Editor) visualModeRoutine() {
	var start time.Time
	ev := e.Screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		start = time.Now()
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
			case 'd':
				e.deleteRuneAtCursor()
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
			}
		}
	}
	if e.PerfAnalysis {
		e.metadata.Elapsed = time.Since(start)
	}
}
