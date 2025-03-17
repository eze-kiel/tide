package editor

import (
	"os"

	"github.com/gdamore/tcell/v2"
)

const (
	VisualMode = iota
	EditMode
	CommandMode
)

var autoSaveOnSwitch = true

type Editor struct {
	Mode   int
	Screen tcell.Screen
	Buffer string

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
