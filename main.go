package main

import (
	"fmt"
	"os"

	"github.com/spf13/afero"
)

func main() {
	options := processArguments()
	driver := NewTermboxDriver()
	ui := NewTerminalUI(&driver)
	fs := afero.OsFs{}
	app := NewApp(&fs)

	options = append(options, SetUI(&ui))
	app.LoadOptions(options...)
	app.Run()
}

func processArguments() Options {
	options, err := ParseArgs(os.Args)

	if err != nil {
		fmt.Println("ERROR: Invalid arguments.", err)
		fmt.Println(Usage())
		os.Exit(1)
	}

	return options
}
