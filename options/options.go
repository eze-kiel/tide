package options

import (
	"fmt"
)

// Opts contains options that are provided using command-line flags
type Opts struct {
	AutoSaveOnSwitch bool
	Theme            string
}

func (o Opts) Verify() error {
	// check that the theme provided is supported
	// the themes are set up in editor/themes.go so any addition there should be
	// backported here to avoid errors
	switch o.Theme {
	case "dark", "light", "valensole":
		return nil
	default:
		return fmt.Errorf("theme '%s' is not supported", o.Theme)
	}
}
