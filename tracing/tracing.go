package tracing

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Tracing struct {
	Elapsed time.Duration
	Traces  []Action
}

type Action struct {
	Kind     string
	Value    string
	Duration time.Duration
}

const (
	// insertions
	InsertRune = iota
	InsertNL
	InsertNLAbove
	InsertNLUnder

	// deletions
	DeleteRuneBefore
	DeleteRunUnder

	// cursor movements
	MoveCL
	MoveCR
	MoveCU
	MoveCD

	// misc
	SwitchMode
)

var Actions = []string{
	// insertions
	InsertRune:    "InsertRune",
	InsertNL:      "InsertNewLine",
	InsertNLAbove: "InsertNewLineAbove",
	InsertNLUnder: "InsertNewLineUnder",

	// deletions
	DeleteRuneBefore: "DeleteRuneBefore",
	DeleteRunUnder:   "DeleteRuneUnder",

	// cursor movements
	MoveCL: "MoveCursorLeft",
	MoveCR: "MoveCursorRight",
	MoveCU: "MoveCursorUp",
	MoveCD: "MoveCursorDown",

	// misc
	SwitchMode: "SwitchMode",
}

func (t Tracing) Dump() error {
	var out []string
	for _, t := range t.Traces {
		out = append(out, fmt.Sprintf("%s,%s,%d", t.Kind, t.Value, t.Duration))
	}
	return os.WriteFile("/tmp/tide.trace", []byte(strings.Join(out, "\n")), 0644)
}
