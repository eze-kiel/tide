package editor

import (
	"strings"
	"time"

	"github.com/eze-kiel/tide/str"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) commandModeRoutine() {
	var start time.Time
	ev := e.Screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		start = time.Now()
		switch ev.Key() {
		case tcell.KeyEsc:
			e.exitCommandMode()
		case tcell.KeyEnter:
			if e.CommandBuffer != "" {
				e.executeCommand(e.CommandBuffer)
				e.CommandBuffer = ""
				e.CommandCursorPos = 0
			} else {
				e.exitCommandMode()
			}

		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if e.CommandCursorPos > 0 {
				e.CommandBuffer = e.CommandBuffer[:e.CommandCursorPos-1] + e.CommandBuffer[e.CommandCursorPos:]
				e.CommandCursorPos--
			}

		case tcell.KeyDelete:
			if e.CommandCursorPos < len(e.CommandBuffer) {
				e.CommandBuffer = e.CommandBuffer[:e.CommandCursorPos] + e.CommandBuffer[e.CommandCursorPos+1:]
			}

		case tcell.KeyLeft:
			if e.CommandCursorPos > 0 {
				e.CommandCursorPos--
			}

		case tcell.KeyRight:
			if e.CommandCursorPos < len(e.CommandBuffer) {
				e.CommandCursorPos++
			}

		case tcell.KeyRune:
			e.CommandBuffer = e.CommandBuffer[:e.CommandCursorPos] + string(ev.Rune()) + e.CommandBuffer[e.CommandCursorPos:]
			e.CommandCursorPos++

		case tcell.KeyHome:
			e.CommandCursorPos = 0

		case tcell.KeyEnd:
			e.CommandCursorPos = len(e.CommandBuffer)
		}
	}
	if e.PerfAnalysis {
		e.metadata.Elapsed = time.Since(start)
	}
}

func (e *Editor) exitCommandMode() {
	e.Mode = VisualMode
	e.CommandBuffer = ""

	e.updateRenderCursor()
}

func (e *Editor) executeCommand(cmd string) {
	cmd = strings.TrimSpace(cmd)

	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "q", "quit":
		e.Quit()

	case "w", "write":
		if len(parts) > 1 {
			e.Filename = parts[1]
		}
		e.SaveToFile(false)
		e.exitCommandMode()

	case "wq", "x":
		if len(parts) > 1 {
			e.Filename = parts[1]
		}
		e.SaveToFile(false)
		e.Quit()

	default:
		e.StatusMsg = str.UnknownCommandErr + parts[0]
		e.StatusTimeout = 5
	}
}
