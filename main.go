package main

import (
	"fmt"
	"os"

	"github.com/dcbishop/jkl/app"
	"github.com/dcbishop/jkl/cli"
	"github.com/dcbishop/jkl/ui"
	"github.com/spf13/afero"
)

func main() {
	options := processArguments()
	driver := ui.NewTermboxDriver()
	ui := ui.NewTerminalUI(&driver)
	fs := afero.OsFs{}
	app := app.New(&fs, &ui)
	app.LoadOptions(options)
	app.Run()
}

func processArguments() cli.Options {
	options, err := cli.ParseArgs(os.Args)

	if err != nil {
		fmt.Println("ERROR: Invalid arguments.", err)
		fmt.Println(cli.Usage())
		os.Exit(1)
	}

	if options.Help {
		fmt.Println(cli.Usage())
		os.Exit(0)
	}

	return options
}
