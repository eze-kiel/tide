package main

import (
	"flag"

	"github.com/eze-kiel/tide/editor"
	"github.com/eze-kiel/tide/file"
)

func main() {
	var trace string
	flag.StringVar(&trace, "trace", "", "trace program execution (can be exec, all)")
	flag.Parse()

	e, err := editor.New()
	if err != nil {
		e.Screen.Fini()
		panic(err)
	}

	switch trace {
	case "exec":
		e.TraceExec = true
	case "all":
		e.TraceAll = true
	default:
		e.Screen.Fini()
		panic("trace " + trace + " is not supported")
	}

	if len(flag.Args()) > 0 {
		e.Filename = flag.Arg(0)
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
