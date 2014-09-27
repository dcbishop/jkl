package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dcbishop/fileaccessor"
	"github.com/dcbishop/jkl/app"
	"github.com/dcbishop/jkl/cli"
)

func main() {
	options := processArguments()
	app := app.New(fileaccessor.LocalStorage{})
	app.LoadOptions(options)
	time.AfterFunc(time.Second, app.Stop)
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
