package cursor

import "github.com/eze-kiel/tide/buffer"

type Cursor struct {
	X, Y int
}

func (c *Cursor) Move(dx, dy int, b buffer.Buffer) {
	lines := b.SplitLines()
	if dx > 0 && c.X < len(lines[c.Y]) {
		c.X += dx
	} else if dx < 0 && c.X > 0 {
		c.X -= dx
	}

	if dy > 0 && c.Y < len(lines)-1 {
		c.Y += dy
	} else if dy < 0 && c.Y > 0 {
		c.Y -= dy
	}

	// adjust X to stay within line bounds after Y movement
	if c.X > len(lines[c.Y]) {
		c.X = len(lines[c.Y])
	}
}
