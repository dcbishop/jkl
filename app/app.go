package app

import "github.com/dcbishop/gim/cli"

// App is the main program
type App struct{}

// New constructs a new app from the given options.
func New() App {
	app := App{}
	return app
}

// LoadOptions loads the given options
func (app *App) LoadOptions(options cli.Options) {

}
