package tracing

import "time"

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
}
