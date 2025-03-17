package editor

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/eze-kiel/tide/buffer"
	"github.com/gdamore/tcell/v2"
)

const (
	VisualMode = iota
	EditMode
	CommandMode
)

var autoSaveOnSwitch = true

type Editor struct {
	sigs chan os.Signal

	Mode     int
	Screen   tcell.Screen
	Filename string

	Width, Height int

	InternalBuffer buffer.Buffer
	RenderBuffer   buffer.Buffer
	InternalCursor struct{ X, Y int }
	RenderCursor   struct{ X, Y int }

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
		sigs:             make(chan os.Signal, 1),
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
	signal.Notify(e.sigs, syscall.SIGWINCH)
	// stuff done by this goroutine should be protected by a mutex so we
	// do not update the size and read the values at the same time
	go func() {
		for range e.sigs {
			e.Width, e.Height = e.Screen.Size()
		}
	}()

	e.Width, e.Height = e.Screen.Size()
	for {
		e.Screen.Clear()

		e.RenderBuffer = e.InternalBuffer.Render()
		lines := e.RenderBuffer.SplitLines()

		for i := 0; i < len(lines) && i < e.Height-1; i++ {
			l := lines[i]
			for j := 0; j < len(l) && j < e.Width; j++ {
				e.Screen.SetContent(j, i, rune(l[j]), nil, tcell.StyleDefault)
			}
		}

		e.Screen.ShowCursor(e.RenderCursor.X, e.RenderCursor.Y)
		e.Screen.Show()

		ev := e.Screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEsc:
				return nil
			case tcell.KeyRight:
				e.moveInternalCursol(1, 0)
				e.updateRenderCursor()
			}
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

// move the editor's internal cursor the desired offset
func (e *Editor) moveInternalCursol(dx, dy int) {
	e.InternalCursor.X += dx
	e.InternalCursor.Y += dy
}

// update the editor's render cursor based on the position of the internal cursor
func (e *Editor) updateRenderCursor() {
	lines := e.R
	if e.RenderCursor.X <= len(e.RenderBuffer.Data) {
		e.RenderCursor.X = e.InternalCursor.X
	}

	if e.RenderCursor.Y <= 
	e.RenderCursor.Y = e.InternalCursor.Y
}
