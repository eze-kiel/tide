package main

import (
	"os"

	"github.com/eze-kiel/tide/editor"
	"github.com/eze-kiel/tide/file"
)

func main() {
	e, err := editor.New()
	if err != nil {
		e.Screen.Fini()
		panic(err)
	}

	if len(os.Args) > 1 {
		e.Filename = os.Args[1]
	}

	defer e.Screen.Fini()

	if file.Exists(e.Filename) {
		e.InternalBuffer.Data, err = file.Read(e.Filename)
		if err != nil {
			e.Screen.Fini()
			panic(err)
		}
	}

	if err := e.Run(); err != nil {
		e.Screen.Fini()
		panic(err)
	}
}
