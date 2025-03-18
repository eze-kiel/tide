package main

import (
	"errors"
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
		e.Crash(err)
	}

	switch trace {
	case "":
	case "exec":
		e.TraceExec = true
	case "all":
		e.TraceAll = true
	default:
		e.Crash(errors.New("trace " + trace + " is not supported"))
	}

	if len(flag.Args()) > 0 {
		e.Filename = flag.Arg(0)
	}

	defer e.Screen.Fini()

	if file.Exists(e.Filename) {
		e.InternalBuffer.Data, err = file.Read(e.Filename)
		if err != nil {
			e.Crash(err)
		}
	}

	if err := e.Run(); err != nil {
		e.Crash(err)
	}
}
