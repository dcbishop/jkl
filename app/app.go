package app

import "github.com/dcbishop/gim/cli"

// App is the main program
type App struct{}

// New constructs a new app from the given options.
func New(options cli.Options) App {
	app := App{}
	return app
}
