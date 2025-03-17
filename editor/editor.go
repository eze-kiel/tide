package editor

import (
	"github.com/gdamore/tcell/v2"
)

const (
	EditMode = iota
	VisualMode
	CommandMode
)

type Editor struct {
	Mode   int
	Screen tcell.Screen
	Buffer string
}

func New() (*Editor, error) {
	e := &Editor{
		Mode: EditMode,
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
