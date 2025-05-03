package actions

type Action struct {
	Kind  Kind
	Value string
	Pos   struct {
		X, Y int
	}
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

type Kind string

var Kinds = []Kind{
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
