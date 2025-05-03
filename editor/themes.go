package editor

import "github.com/gdamore/tcell/v2"

// setTheme sets the editor's theme for foreground and background colors
// any new theme added here should also be inserted in options/options.go to pass
// the pre-flight checks
func (e *Editor) setTheme() {
	switch e.theme {
	case "dark":
		e.backgroundColor = tcell.ColorBlack
		e.foregroundColor = tcell.ColorWhiteSmoke
		e.highlightColor = tcell.ColorRebeccaPurple
	case "light": // todo: fix the cursor color
		e.backgroundColor = tcell.ColorWhiteSmoke
		e.foregroundColor = tcell.ColorBlack
		e.highlightColor = tcell.ColorDarkOliveGreen
	case "valensole":
		e.backgroundColor = tcell.ColorRebeccaPurple
		e.foregroundColor = tcell.ColorWhiteSmoke
		e.highlightColor = tcell.ColorDarkOliveGreen
	}
}
