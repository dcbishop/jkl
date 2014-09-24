package main

import (
	"fmt"
	"os"

	"github.com/dcbishop/gim/app"
	"github.com/dcbishop/gim/cli"
)

func main() {
	options := processArguments()
	app := app.New()
	app.LoadOptions(options)
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
