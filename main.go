package main

import (
	"errors"
	"os"

	"github.com/eze-kiel/tide/editor"
	"github.com/eze-kiel/tide/file"
)

func main() {
	e, err := editor.New()
	if err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		panic(errors.New("missing file name"))
	}
	e.Filename = os.Args[1]

	defer e.Screen.Fini()

	if file.Exists(e.Filename) {
		e.Buffer, err = file.Read(e.Filename)
		if err != nil {
			panic(err)
		}
	}

	if err := e.Run(); err != nil {
		panic(err)
	}
}
