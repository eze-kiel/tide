package main

import (
	"flag"

	"github.com/eze-kiel/tide/editor"
	"github.com/eze-kiel/tide/file"
)

func main() {
	var perf bool
	flag.BoolVar(&perf, "perf", false, "start with performance analysis parameters")
	flag.Parse()

	e, err := editor.New()
	if err != nil {
		e.Screen.Fini()
		panic(err)
	}
	e.PerfAnalysis = perf

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
