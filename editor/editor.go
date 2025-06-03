package editor

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/eze-kiel/tide/actions"
	"github.com/eze-kiel/tide/buffer"
	"github.com/eze-kiel/tide/cursor"
	"github.com/eze-kiel/tide/file"
	"github.com/eze-kiel/tide/options"
	"github.com/eze-kiel/tide/str"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

const (
	VisualMode = iota
	EditMode
	CommandMode

	LineNumberWidth = 5
)

var DefaultMsgTimeout = 5

type Editor struct {
	sigs chan os.Signal

	Mode     int
	Screen   tcell.Screen
	Filename string

	Width, Height    int
	OffsetX, OffsetY int

	InternalBuffer, RenderBuffer buffer.Buffer
	InternalCursor, RenderCursor cursor.Cursor

	PreviousActions []actions.Action // hold the previous internalBuffer changes, for the undo mechanism

	CommandBuffer    string
	CommandCursorPos int

	Selection struct {
		Line         int
		StartX, EndX int
		Content      string
	}

	Clipboard string

	StatusMsg     string
	StatusTimeout int
	fileChanged   bool

	theme           string
	backgroundColor tcell.Color
	foregroundColor tcell.Color
	highlightColor  tcell.Color

	fastJumpLength   int  // how far you go when you hit D or U in VISU mode
	autoSaveOnSwitch bool // auto save when going from EDIT to VISU modes
}

func New(o options.Opts) (*Editor, error) {
	e := &Editor{
		sigs:             make(chan os.Signal, 1),
		Mode:             VisualMode,
		autoSaveOnSwitch: o.AutoSaveOnSwitch,
		fileChanged:      false,
		theme:            o.Theme,
	}

	e.setTheme()

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
	e.fastJumpLength = (e.Height / 3)

	for {
		e.Screen.Clear()
		e.RenderBuffer = e.InternalBuffer.Translate()
		lines := e.RenderBuffer.SplitLines()

		// set the background color of the whole editor screen
		for y := range e.Height {
			for x := range e.Width {
				e.Screen.SetContent(x, y, rune(0), nil, tcell.StyleDefault.
					Background(e.backgroundColor))
			}
		}

		for i := e.OffsetY; i < len(lines) && i < e.OffsetY+e.Height-2; i++ {
			lineNumStr := fmt.Sprintf("%*d ", LineNumberWidth-1, i+1)
			style := tcell.StyleDefault.
				Background(e.backgroundColor).
				Foreground(e.foregroundColor)

			if i == e.InternalCursor.Y {
				style = style.
					Background(e.highlightColor).
					Foreground(e.foregroundColor).
					Bold(true)
			}
			for j, r := range lineNumStr {
				if j < LineNumberWidth {
					e.Screen.SetContent(j, i-e.OffsetY, r, nil, tcell.StyleDefault.
						Background(e.backgroundColor).
						Foreground(e.foregroundColor))
				}
			}

			for j, r := range lineNumStr {
				if j < LineNumberWidth {
					e.Screen.SetContent(j, i-e.OffsetY, r, nil, style)
				}
			}

			l := lines[i]
			lineRunes := []rune(l)
			renderX := 0
			for runeIdx := 0; runeIdx < len(lineRunes) && renderX < e.Width-LineNumberWidth; runeIdx++ {
				r := lineRunes[runeIdx]
				charWidth := 1
				if r == '\t' {
					// Handle tab expansion
					tabPos := renderX % buffer.TAB_SIZE
					charWidth = buffer.TAB_SIZE - tabPos
					for k := 0; k < charWidth && renderX+k < e.Width-LineNumberWidth; k++ {
						if renderX+k >= e.OffsetX {
							e.Screen.SetContent(LineNumberWidth+(renderX+k-e.OffsetX), i-e.OffsetY, ' ', nil, tcell.StyleDefault.
								Background(e.backgroundColor).
								Foreground(e.foregroundColor))
						}
					}
				} else {
					// Handle regular characters including wide ones
					charWidth = runewidth.RuneWidth(r)
					if charWidth == 0 {
						charWidth = 1 // control characters
					}
					if renderX >= e.OffsetX && renderX < e.OffsetX+e.Width-LineNumberWidth {
						e.Screen.SetContent(LineNumberWidth+(renderX-e.OffsetX), i-e.OffsetY, r, nil, tcell.StyleDefault.
							Background(e.backgroundColor).
							Foreground(e.foregroundColor))
					}
					// For wide characters, fill additional columns with spaces
					for k := 1; k < charWidth && renderX+k < e.Width-LineNumberWidth; k++ {
						if renderX+k >= e.OffsetX {
							e.Screen.SetContent(LineNumberWidth+(renderX+k-e.OffsetX), i-e.OffsetY, ' ', nil, tcell.StyleDefault.
								Background(e.backgroundColor).
								Foreground(e.foregroundColor))
						}
					}
				}
				renderX += charWidth
			}
		}

		if e.Selection.Content != "" {
			selectionRunes := []rune(e.Selection.Content)
			renderX := e.Selection.StartX
			for runeIdx := 0; runeIdx < len(selectionRunes) && renderX < e.Selection.EndX; runeIdx++ {
				r := selectionRunes[runeIdx]
				charWidth := runewidth.RuneWidth(r)
				if charWidth == 0 {
					charWidth = 1
				}
				if renderX >= e.OffsetX && renderX < e.OffsetX+e.Width-LineNumberWidth {
					e.Screen.SetContent(LineNumberWidth+(renderX-e.OffsetX), e.Selection.Line-e.OffsetY, r, nil, tcell.StyleDefault.
						Background(e.highlightColor).
						Foreground(e.foregroundColor))
				}
				// For wide characters, fill additional columns
				for k := 1; k < charWidth && renderX+k < e.Selection.EndX; k++ {
					if renderX+k >= e.OffsetX && renderX+k < e.OffsetX+e.Width-LineNumberWidth {
						e.Screen.SetContent(LineNumberWidth+(renderX+k-e.OffsetX), e.Selection.Line-e.OffsetY, ' ', nil, tcell.StyleDefault.
							Background(e.highlightColor).
							Foreground(e.foregroundColor))
					}
				}
				renderX += charWidth
			}
		}

		switch e.Mode {
		case EditMode:
			for i, r := range str.EditMode {
				e.Screen.SetContent(i, e.Height-1, r, nil, tcell.StyleDefault.
					Background(e.highlightColor).
					Foreground(e.foregroundColor))
			}
		case VisualMode:
			for i, r := range str.VisualMode {
				e.Screen.SetContent(i, e.Height-1, r, nil, tcell.StyleDefault.
					Background(e.highlightColor).
					Foreground(e.foregroundColor))
			}
		case CommandMode:
			for i := range e.Width {
				e.Screen.SetContent(i, e.Height-1, rune(0), nil, tcell.StyleDefault.
					Background(e.backgroundColor).
					Foreground(e.foregroundColor))
			}
			e.Screen.SetContent(0, e.Height-1, ':', nil, tcell.StyleDefault.
				Background(e.backgroundColor).
				Foreground(e.foregroundColor))
			for i, r := range e.CommandBuffer {
				e.Screen.SetContent(i+1, e.Height-1, r, nil, tcell.StyleDefault.
					Background(e.backgroundColor).
					Foreground(e.foregroundColor))
			}
		}

		if e.StatusMsg != "" && e.StatusTimeout > 0 {
			e.StatusTimeout--
			for i, r := range e.StatusMsg {
				if i < e.Width {
					e.Screen.SetContent(e.Width-len(e.StatusMsg)+i, e.Height-1, r, nil, tcell.StyleDefault.
						Background(e.backgroundColor).
						Foreground(e.foregroundColor))
				}
			}
		}

		if e.Mode == CommandMode {
			e.Screen.ShowCursor(len(e.CommandBuffer)+1, e.Height-1)
		} else {
			e.Screen.ShowCursor(LineNumberWidth+(e.RenderCursor.X-e.OffsetX), e.RenderCursor.Y-e.OffsetY)
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

// big brain time
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

// crash properly when possible
func (e Editor) Crash(err error) {
	e.Screen.Fini()
	panic(err)
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

	// limit x to the boundaries based on rune count
	lineRunes := []rune(lines[y])
	if x < 0 {
		x = 0
	} else if x > len(lineRunes) {
		x = len(lineRunes)
	}

	rx = 0
	for i := 0; i < x && i < len(lineRunes); i++ {
		if lineRunes[i] == '\t' {
			// align to the next tab stop
			rx += buffer.TAB_SIZE - (rx % buffer.TAB_SIZE)
		} else {
			// account for wide characters
			rx += runewidth.RuneWidth(lineRunes[i])
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

	// compute new X without going offlimits (using rune count)
	lineLength := buffer.RuneLength(lines[newY])
	newX := e.InternalCursor.X + dx
	if newX < 0 {
		newX = 0
	} else if newX > lineLength {
		newX = lineLength
	}

	// when moving to a short line, do not end up in the abyss but go to the
	// end of the next line
	if dy != 0 && newX > lineLength {
		newX = lineLength
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

	// adjust horizontal scrolling to account for line number width
	if e.RenderCursor.X >= e.OffsetX+e.Width-LineNumberWidth {
		e.OffsetX = e.RenderCursor.X - (e.Width - LineNumberWidth) + 1
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
	currentLineRunes := []rune(lines[y])

	// Insert the rune at the correct position
	newRunes := make([]rune, 0, len(currentLineRunes)+1)
	newRunes = append(newRunes, currentLineRunes[:x]...)
	newRunes = append(newRunes, ch)
	newRunes = append(newRunes, currentLineRunes[x:]...)

	e.PreviousActions = append(e.PreviousActions, actions.Action{
		Kind:  actions.Kinds[actions.InsertRune],
		Value: string(ch),
		Pos: struct {
			X int
			Y int
		}{x, y},
	})

	lines[y] = string(newRunes)
	e.updateBufferFromLines(lines)

	e.InternalCursor.X++
	e.updateRenderCursor()
}

// insert a newline at the current cursor position
func (e *Editor) insertNewlineAtCursor() {
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
	currentLineRunes := []rune(lines[y])

	firstPart := ""
	if x > 0 && x <= len(currentLineRunes) {
		firstPart = string(currentLineRunes[:x])
	}
	secondPart := ""
	if x < len(currentLineRunes) {
		secondPart = string(currentLineRunes[x:])
	}

	newLines := make([]string, 0, len(lines)+1)
	newLines = append(newLines, lines[:y]...)
	newLines = append(newLines, firstPart)
	newLines = append(newLines, secondPart)
	if y+1 < len(lines) {
		newLines = append(newLines, lines[y+1:]...)
	}

	e.PreviousActions = append(e.PreviousActions, actions.Action{
		Kind:  actions.Kinds[actions.InsertNL],
		Value: "\n",
		Pos: struct {
			X int
			Y int
		}{x, y},
	})

	e.updateBufferFromLines(newLines)

	e.InternalCursor.X = 0
	e.InternalCursor.Y = y + 1
	e.updateRenderCursor()
}

// insert a newline above the current cursor position
func (e *Editor) insertNewlineAbove() {
	lines := e.InternalBuffer.SplitLines()

	if len(lines) == 0 {
		e.InternalBuffer.Data = "\n"
		e.InternalCursor.X = 0
		e.InternalCursor.Y = 1
		e.updateRenderCursor()
		return
	}

	y := e.InternalCursor.Y

	newLines := make([]string, 0, len(lines)+1)
	newLines = append(newLines, lines[:y]...)
	newLines = append(newLines, "")
	newLines = append(newLines, lines[y:]...)

	e.updateBufferFromLines(newLines)

	e.InternalCursor.X = 0
	e.InternalCursor.Y = y
	e.updateRenderCursor()
}

// insert a newline under the current cursor position
func (e *Editor) insertNewlineUnder() {
	lines := e.InternalBuffer.SplitLines()

	if len(lines) == 0 {
		e.InternalBuffer.Data = "\n"
		e.InternalCursor.X = 0
		e.InternalCursor.Y = 1
		e.updateRenderCursor()
		return
	}

	y := e.InternalCursor.Y

	newLines := make([]string, 0, len(lines)+1)
	newLines = append(newLines, lines[:y+1]...)
	newLines = append(newLines, "")
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

		e.InternalCursor.X = buffer.RuneLength(prevLine)
		e.InternalCursor.Y = y - 1

		e.updateRenderCursor()
		return
	}

	currentLineRunes := []rune(lines[y])
	if x > 0 && x <= len(currentLineRunes) {
		newRunes := make([]rune, 0, len(currentLineRunes)-1)
		newRunes = append(newRunes, currentLineRunes[:x-1]...)
		newRunes = append(newRunes, currentLineRunes[x:]...)
		lines[y] = string(newRunes)

		e.updateBufferFromLines(lines)

		e.InternalCursor.X--
	}

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
	currentLineRunes := []rune(lines[y])

	// if at end of line but not last line, join with next line
	if x == len(currentLineRunes) && y < len(lines)-1 {
		nextLine := lines[y+1]

		newLine := lines[y] + nextLine

		newLines := make([]string, 0, len(lines)-1)
		newLines = append(newLines, lines[:y]...)
		newLines = append(newLines, newLine)
		if y+2 < len(lines) {
			newLines = append(newLines, lines[y+2:]...)
		}

		e.updateBufferFromLines(newLines)
	}

	if x < len(currentLineRunes) {
		newRunes := make([]rune, 0, len(currentLineRunes)-1)
		newRunes = append(newRunes, currentLineRunes[:x]...)
		newRunes = append(newRunes, currentLineRunes[x+1:]...)
		lines[y] = string(newRunes)

		e.updateBufferFromLines(lines)
	}

	e.updateRenderCursor()
}

// delete the character at a given position
func (e *Editor) deleteRuneAt(x, y int) {
	lines := e.InternalBuffer.SplitLines()

	if len(lines) == 0 || y < 0 || y >= len(lines) {
		return
	}

	e.InternalCursor.X = x
	e.InternalCursor.Y = y
	currentLineRunes := []rune(lines[y])

	// if at end of line but not last line, join with next line
	if x == len(currentLineRunes) && y < len(lines)-1 {
		nextLine := lines[y+1]

		newLine := lines[y] + nextLine

		newLines := make([]string, 0, len(lines)-1)
		newLines = append(newLines, lines[:y]...)
		newLines = append(newLines, newLine)
		if y+2 < len(lines) {
			newLines = append(newLines, lines[y+2:]...)
		}

		e.updateBufferFromLines(newLines)
	}

	if x >= 0 && x < len(currentLineRunes) {
		newRunes := make([]rune, 0, len(currentLineRunes)-1)
		newRunes = append(newRunes, currentLineRunes[:x]...)
		newRunes = append(newRunes, currentLineRunes[x+1:]...)
		lines[y] = string(newRunes)

		e.updateBufferFromLines(lines)
	}

	e.updateRenderCursor()
}

// helper function to update the internal buffer from an array of lines
func (e *Editor) updateBufferFromLines(lines []string) {
	e.fileChanged = true
	e.InternalBuffer.Data = strings.Join(lines, "\n")
}

// save internal buffer to file
func (e *Editor) SaveToFile() {
	if err := file.Write(e.Filename, e.InternalBuffer.Data); err != nil {
		e.StatusMsg = "Error: " + err.Error()
	} else {
		if e.autoSaveOnSwitch {
			e.StatusMsg = str.AutoSavedMsg + e.Filename
		} else {
			e.StatusMsg = str.SavedMsg + e.Filename
		}
	}
	e.StatusTimeout = DefaultMsgTimeout
}

func (e *Editor) replaceRuneUnder() {
	ev := e.Screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyRune:
			e.deleteRuneAtCursor()
			e.insertRune(ev.Rune())
		}
	}
}

func (e *Editor) selectLine() {
	lines := e.InternalBuffer.SplitLines()

	if len(lines) == 0 {
		return
	}

	y := e.InternalCursor.Y
	// Use the rendered version for display, but track the actual line content
	renderLine := strings.ReplaceAll(lines[y], "\t", strings.Repeat(buffer.TAB_SYMBOL, buffer.TAB_SIZE))
	e.Selection.Content = renderLine
	e.Selection.StartX = 0
	e.Selection.EndX = buffer.DisplayWidth(renderLine)
	e.Selection.Line = y
}

func (e *Editor) cancelSelection() {
	lines := e.InternalBuffer.SplitLines()

	if len(lines) == 0 {
		return
	}

	e.Selection = struct {
		Line    int
		StartX  int
		EndX    int
		Content string
	}{}
}

func (e *Editor) copySelection() {
	e.Clipboard = e.Selection.Content
}

func (e *Editor) deleteSelection() {
	if e.Selection.Content == "" {
		return
	}

	lines := e.InternalBuffer.SplitLines()
	y := e.Selection.Line

	if y < 0 || y >= len(lines) {
		return
	}

	startX := e.renderToInternalX(e.Selection.StartX, y)
	endX := e.renderToInternalX(e.Selection.EndX, y)

	currentLine := lines[y]

	if startX < 0 {
		startX = 0
	}
	if startX > len(currentLine) {
		startX = len(currentLine)
	}

	if endX < 0 {
		endX = 0
	}
	if endX > len(currentLine) {
		endX = len(currentLine)
	}

	if startX >= endX {
		return
	}

	newLine := currentLine[:startX] + currentLine[endX:]
	lines[y] = newLine
	e.updateBufferFromLines(lines)

	e.InternalCursor.X = startX
	e.InternalCursor.Y = y
	e.updateRenderCursor()

	e.Selection.Content = ""
}

func (e *Editor) renderToInternalX(renderX, y int) int {
	lines := e.InternalBuffer.SplitLines()
	if y < 0 || y >= len(lines) {
		return -1
	}
	lineRunes := []rune(lines[y])

	internalX := 0
	renderCol := 0

	for i := 0; i < len(lineRunes); i++ {
		if lineRunes[i] == '\t' {
			tabStop := (renderCol/buffer.TAB_SIZE + 1) * buffer.TAB_SIZE
			if renderX < tabStop {
				return internalX
			}
			renderCol = tabStop
		} else {
			// account for wide characters
			renderCol += runewidth.RuneWidth(lineRunes[i])
		}

		if renderCol > renderX {
			return internalX
		}
		internalX++
	}
	return internalX
}

func (e *Editor) pasteUnder() {
	if e.Clipboard == "" {
		return
	}

	lines := e.InternalBuffer.SplitLines()

	if len(lines) == 0 {
		e.InternalBuffer.Data = e.Clipboard
		e.InternalCursor.X = 0
		e.InternalCursor.Y = 1
		e.updateRenderCursor()
		return
	}

	y := e.InternalCursor.Y

	newLines := make([]string, 0, len(lines)+1)
	newLines = append(newLines, lines[:y+1]...)
	newLines = append(newLines, e.Clipboard)
	if y+1 < len(lines) {
		newLines = append(newLines, lines[y+1:]...)
	}
	e.cancelSelection()
	e.updateBufferFromLines(newLines)

	e.InternalCursor.X = 0
	e.InternalCursor.Y = y + 1
	e.updateRenderCursor()
}

func (e *Editor) toggleCommentLine() {
	lines := e.InternalBuffer.SplitLines()

	if len(lines) == 0 {
		e.InternalCursor.X = 0
		e.InternalCursor.Y = 1
		e.updateRenderCursor()
		return
	}

	y := e.InternalCursor.Y

	// todo: use specific comment based on the file extension if known,
	// otherwise go to default

	newLines := make([]string, 0, len(lines)+1)
	parts := strings.Fields(lines[y])
	if len(parts) == 0 || !strings.Contains(parts[0], str.Comment) {
		// comment the line
		newLines = append(newLines, lines[:y]...)
		newLines = append(newLines, strings.Join([]string{str.Comment, lines[y]}, " "))
		if y+1 < len(lines) {
			newLines = append(newLines, lines[y+1:]...)
		}
	} else {
		// uncomment the line
		newLines = append(newLines, lines[:y]...)
		newLines = append(newLines, strings.Replace(lines[y], str.Comment+" ", "", 1))
		if y+1 < len(lines) {
			newLines = append(newLines, lines[y+1:]...)
		}
	}

	e.updateBufferFromLines(newLines)

	e.InternalCursor.X = len(lines[y]) + len(str.Comment)
	e.InternalCursor.Y = y
	e.updateRenderCursor()
}

func (e *Editor) undo() {
	if len(e.PreviousActions) == 0 {
		e.StatusMsg = str.NoMoreUndoMsg
		e.StatusTimeout = DefaultMsgTimeout
		return
	}

	// save the latest action, then drop it
	prev := e.PreviousActions[len(e.PreviousActions)-1]
	e.PreviousActions = e.PreviousActions[:len(e.PreviousActions)-1]

	switch prev.Kind {
	case actions.Kinds[actions.InsertRune], actions.Kinds[actions.InsertNL]:
		e.deleteRuneAt(prev.Pos.X, prev.Pos.Y)
	}
}
