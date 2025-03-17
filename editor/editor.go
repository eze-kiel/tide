package editor

import (
	"os"

	"github.com/eze-kiel/tide/buffer"
	"github.com/eze-kiel/tide/str"
	"github.com/gdamore/tcell/v2"
)

const (
	VisualMode = iota
	EditMode
	CommandMode
)

var autoSaveOnSwitch = true

type Editor struct {
	Mode     int
	Screen   tcell.Screen
	Buffer   string
	Filename string

	CursorX, CursorY int
	OffsetX, OffsetY int

	StatusMsg     string
	StatusTimeout int

	/*
		stuff that will be configurable in the future starts here
	*/
	fastJumpLength   int  // how far you go when you hit D or U in VISU mode
	autoSaveOnSwitch bool // auto save when going from EDIT to VISU modes
}

func New() (*Editor, error) {
	e := &Editor{
		Mode:             VisualMode,
		fastJumpLength:   10,
		autoSaveOnSwitch: true,
	}

	var err error
	e.Screen, err = tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err = e.Screen.Init(); err != nil {
		return nil, err
	}

	return e, nil
}

func (e *Editor) Run() error {
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
		case EditMode:
			for i, r := range str.EditMode {
				e.Screen.SetContent(i, sheight-1, r, nil, tcell.StyleDefault.Reverse(true))
			}
		case VisualMode:
			for i, r := range str.VisualMode {
				e.Screen.SetContent(i, sheight-1, r, nil, tcell.StyleDefault.Reverse(true))
			}
		case CommandMode:
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
		case EditMode:
			e.EditModeRoutine(lines)
		case VisualMode:
			e.VisualModeRoutine(lines)
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

// if e.Mode is 1, then e.Mode ^ (EditMode | VisualMode) -> 1 ^ (1 | 2) -> 1 ^ 3 = 2
// if e.Mode is 2, then e.Mode ^ (EditMode | VisualMode) -> 2 ^ (1 | 2) -> 2 ^ 3 = 1
func (e *Editor) SwitchMode() {
	if e.Mode == EditMode || e.Mode == VisualMode {
		e.Mode ^= (EditMode | VisualMode)
	}
}

func (e Editor) Quit() {
	e.Screen.Fini()
	os.Exit(0)
}
