package main

import (
	"flag"

	"github.com/eze-kiel/tide/editor"
	"github.com/eze-kiel/tide/file"
	"github.com/eze-kiel/tide/options"
)

func main() {
	var o options.Opts
	flag.BoolVar(&o.AutoSaveOnSwitch, "autosave-on-switch", false, "enable autosave when switching modes")
	flag.StringVar(&o.Theme, "color-theme", "dark", "set color theme (can be 'dark', 'light', 'valensole')")
	flag.Parse()

	if err := o.Verify(); err != nil {
		panic(err)
	}

	e, err := editor.New(o)
	if err != nil {
		e.Crash(err)
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
