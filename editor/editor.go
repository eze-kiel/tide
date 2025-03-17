package editor

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/eze-kiel/tide/buffer"
	"github.com/eze-kiel/tide/cursor"
	"github.com/eze-kiel/tide/file"
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
	sigs chan os.Signal

	Mode     int
	Screen   tcell.Screen
	Filename string

	Width, Height    int
	OffsetX, OffsetY int

	InternalBuffer, RenderBuffer buffer.Buffer
	InternalCursor, RenderCursor cursor.Cursor

	CommandBuffer    string
	CommandCursorPos int

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
		autoSaveOnSwitch: false,
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
	e.fastJumpLength = (e.Height / 2) + 2
	for {
		e.Screen.Clear()

		e.RenderBuffer = e.InternalBuffer.Render()
		lines := e.RenderBuffer.SplitLines()

		for i := e.OffsetY; i < len(lines) && i < e.OffsetY+e.Height-2; i++ {
			l := lines[i]
			for j := e.OffsetX; j < len(l) && j < e.OffsetX+e.Width; j++ {
				e.Screen.SetContent(j-e.OffsetX, i-e.OffsetY, rune(l[j]), nil, tcell.StyleDefault)
			}
		}

		// the switch below could be simplified?
		switch e.Mode {
		case EditMode:
			for i, r := range str.EditMode {
				e.Screen.SetContent(i, e.Height-1, r, nil, tcell.StyleDefault.Reverse(true))
			}
		case VisualMode:
			for i, r := range str.VisualMode {
				e.Screen.SetContent(i, e.Height-1, r, nil, tcell.StyleDefault.Reverse(true))
			}
		case CommandMode:
			for i := 0; i < e.Width; i++ {
				e.Screen.SetContent(i, e.Height-1, ' ', nil, tcell.StyleDefault)
			}

			e.Screen.SetContent(0, e.Height-1, ':', nil, tcell.StyleDefault)
			for i, r := range e.CommandBuffer {
				e.Screen.SetContent(i+1, e.Height-1, r, nil, tcell.StyleDefault)
			}
		}

		if e.StatusMsg != "" && e.StatusTimeout > 0 {
			e.StatusTimeout--
			for i, r := range e.StatusMsg {
				if i < e.Width {
					e.Screen.SetContent(e.Width-len(e.StatusMsg)+i, e.Height-1, r, nil, tcell.StyleDefault)
				}
			}
		}

		if e.Mode == CommandMode {
			e.Screen.ShowCursor(len(e.CommandBuffer)+1, e.Height-1)
		} else {
			e.Screen.ShowCursor(e.RenderCursor.X-e.OffsetX, e.RenderCursor.Y-e.OffsetY)
		}
		e.Screen.Show()

		switch e.Mode {
		case EditMode:
			e.editModeRoutine()
		case VisualMode:
			e.visualModeRoutine()
		case CommandMode:
			e.commandModeRoutine()
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

// properly quit the editor
func (e Editor) Quit() {
	e.Screen.Fini()
	os.Exit(0)
}

// map internal buffer position to render buffer position
func (e *Editor) internalToRenderPos(x, y int) (rx, ry int) {
	// y is usually the same in both buffers unless there is folding but it's
	// not implemented yet
	ry = y

	// for x we need to count expanded characters
	lines := e.InternalBuffer.SplitLines()
	if y < 0 || y >= len(lines) {
		return 0, y
	}

	// limit x to the boundaries
	if x < 0 {
		x = 0
	} else if x > len(lines[y]) {
		x = len(lines[y])
	}

	rx = 0
	for i := 0; i < x; i++ {
		if lines[y][i] == '\t' {
			// align to the next tab stop
			rx += buffer.TAB_SIZE - (rx % buffer.TAB_SIZE)
		} else {
			rx++
		}
	}

	return rx, ry
}

func (e *Editor) moveInternalCursor(dx, dy int) {
	lines := e.InternalBuffer.SplitLines()

	// empty buffer
	if len(lines) == 0 {
		e.InternalCursor.X = 0
		e.InternalCursor.Y = 0
		e.updateRenderCursor()
		return
	}

	// compute new Y without going offlimits
	newY := e.InternalCursor.Y + dy
	if newY < 0 {
		newY = 0
	} else if newY >= len(lines) {
		newY = len(lines) - 1
	}

	// compute new X without going offlimits
	newX := e.InternalCursor.X + dx
	if newX < 0 {
		newX = 0
	} else if newX > len(lines[newY]) {
		newX = len(lines[newY])
	}

	// when moving to a short line, do not end up in the abyss but go to the
	// end of the next line
	if dy != 0 && newX > len(lines[newY]) {
		newX = len(lines[newY])
	}

	e.InternalCursor.X = newX
	e.InternalCursor.Y = newY

	e.updateRenderCursor()
}

// update the editor's render cursor based on the position of the internal cursor
func (e *Editor) updateRenderCursor() {
	e.RenderCursor.X, e.RenderCursor.Y = e.internalToRenderPos(e.InternalCursor.X, e.InternalCursor.Y)
	e.handleScrolling()
}

func (e *Editor) handleScrolling() {
	if e.RenderCursor.Y < e.OffsetY {
		e.OffsetY = e.RenderCursor.Y
	}
	if e.RenderCursor.Y >= e.OffsetY+e.Height-2 {
		e.OffsetY = e.RenderCursor.Y - (e.Height - 3)
	}
	if e.RenderCursor.X >= e.OffsetX+e.Width {
		e.OffsetX = e.RenderCursor.X - e.Width + 1
	}
	if e.RenderCursor.X < e.OffsetX {
		e.OffsetX = e.RenderCursor.X
	}
}

// insert a character at the current cursor position
func (e *Editor) insertRune(ch rune) {
	lines := e.InternalBuffer.SplitLines()

	if len(lines) == 0 {
		e.InternalBuffer.Data = "" + string(ch)
		e.InternalCursor.X = 1
		e.InternalCursor.Y = 0
		e.updateRenderCursor()
		return
	}

	y := e.InternalCursor.Y
	x := e.InternalCursor.X
	currentLine := lines[y]

	newLine := ""
	if x > 0 {
		newLine = currentLine[:x]
	}
	newLine += string(ch)
	if x < len(currentLine) {
		newLine += currentLine[x:]
	}

	lines[y] = newLine
	e.updateBufferFromLines(lines)

	e.InternalCursor.X++
	e.updateRenderCursor()
}

// insert a newline at the current cursor position
func (e *Editor) insertNewline() {
	lines := e.InternalBuffer.SplitLines()

	if len(lines) == 0 {
		e.InternalBuffer.Data = "\n"
		e.InternalCursor.X = 0
		e.InternalCursor.Y = 1
		e.updateRenderCursor()
		return
	}

	y := e.InternalCursor.Y
	x := e.InternalCursor.X
	currentLine := lines[y]

	firstPart := ""
	if x > 0 {
		firstPart = currentLine[:x]
	}
	secondPart := ""
	if x < len(currentLine) {
		secondPart = currentLine[x:]
	}

	newLines := make([]string, 0, len(lines)+1)
	newLines = append(newLines, lines[:y]...)
	newLines = append(newLines, firstPart)
	newLines = append(newLines, secondPart)
	if y+1 < len(lines) {
		newLines = append(newLines, lines[y+1:]...)
	}

	e.updateBufferFromLines(newLines)

	e.InternalCursor.X = 0
	e.InternalCursor.Y = y + 1
	e.updateRenderCursor()
}

// delete the character before the cursor (backspace)
func (e *Editor) deleteRuneBeforeCursor() {
	lines := e.InternalBuffer.SplitLines()

	y := e.InternalCursor.Y
	x := e.InternalCursor.X

	// nothing to delete in empty buffer or at the beginning
	if len(lines) == 0 || (y == 0 && x == 0) {
		return
	}

	// if at beginning of line but not first line, join with previous line
	if x == 0 && y > 0 {
		prevLine := lines[y-1]
		currentLine := lines[y]

		newLine := prevLine + currentLine

		newLines := make([]string, 0, len(lines)-1)
		newLines = append(newLines, lines[:y-1]...)
		newLines = append(newLines, newLine)
		if y+1 < len(lines) {
			newLines = append(newLines, lines[y+1:]...)
		}

		e.updateBufferFromLines(newLines)

		e.InternalCursor.X = len(prevLine)
		e.InternalCursor.Y = y - 1

		e.updateRenderCursor()
		return
	}

	currentLine := lines[y]
	newLine := currentLine[:x-1] + currentLine[x:]
	lines[y] = newLine

	e.updateBufferFromLines(lines)

	e.InternalCursor.X--

	e.updateRenderCursor()
}

// delete the character at the cursor position (delete key)
func (e *Editor) deleteRuneAtCursor() {
	lines := e.InternalBuffer.SplitLines()

	if len(lines) == 0 {
		return
	}

	y := e.InternalCursor.Y
	x := e.InternalCursor.X
	currentLine := lines[y]

	// if at end of line but not last line, join with next line
	if x == len(currentLine) && y < len(lines)-1 {
		nextLine := lines[y+1]

		newLine := currentLine + nextLine

		newLines := make([]string, 0, len(lines)-1)
		newLines = append(newLines, lines[:y]...)
		newLines = append(newLines, newLine)
		if y+2 < len(lines) {
			newLines = append(newLines, lines[y+2:]...)
		}

		e.updateBufferFromLines(newLines)
	}

	if x < len(currentLine) {
		newLine := currentLine[:x] + currentLine[x+1:]
		lines[y] = newLine

		e.updateBufferFromLines(lines)
	}

	e.updateRenderCursor()
}

// helper function to update the internal buffer from an array of lines
func (e *Editor) updateBufferFromLines(lines []string) {
	e.InternalBuffer.Data = strings.Join(lines, "\n")
}

// save internal buffer to file
func (e *Editor) SaveToFile(autosave bool) {
	if err := file.Write(e.Filename, e.InternalBuffer.Data); err != nil {
		e.StatusMsg = "Error: " + err.Error()
	} else {
		if autosave {
			e.StatusMsg = str.AutoSavedMsg + e.Filename
		} else {
			e.StatusMsg = str.SavedMsg + e.Filename
		}
	}
	e.StatusTimeout = 5
}
