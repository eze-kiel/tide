package cursor

import (
	"github.com/eze-kiel/tide/buffer"
)

type Cursor struct {
	X, Y int
}

func (c *Cursor) Move(dx, dy int, b buffer.Buffer) {
	lines := b.SplitLines()
	if len(lines) == 0 {
		c.X = 0
		c.Y = 0
		return
	}

	// Handle vertical movement
	newY := c.Y + dy
	if newY < 0 {
		newY = 0
	} else if newY >= len(lines) {
		newY = len(lines) - 1
	}

	// Handle horizontal movement using rune count
	lineLength := buffer.RuneLength(lines[newY])
	newX := c.X + dx
	if newX < 0 {
		newX = 0
	} else if newX > lineLength {
		newX = lineLength
	}

	c.X = newX
	c.Y = newY
}
