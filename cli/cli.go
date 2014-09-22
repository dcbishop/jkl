package cli

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/docopt/docopt-go"
)

// UsageMessage is the message displayed in the console giving the programs usage
var UsageMessage = `gim

Usage:
  gim [<file>...]
  gim -h | --help

Options:
  -h --help     Show this screen.
`

// Options stores options parsed from the command line
type Options struct {
	FilesToOpen []string
	Help        bool
}

// ParseArgs takes arguments and returns a cli.Options. Will return error if parsing failed.
func ParseArgs(args []string) (Options, error) {
	options := Options{}

	if len(args) < 2 {
		return options, nil
	}

	disableStdout()
	defer restoreStdout()

	arguments, err := docopt.Parse(Usage(), args[1:], false, "gim 0.1", false, false)
	if err != nil {
		return options, err
	}

	if arguments["--help"].(bool) {
		options.Help = true
		return options, nil
	}

	options.FilesToOpen = arguments["<file>"].([]string)

	app := cli.NewApp()
	app.Name = "greet"
	app.Usage = "test"

	return options, nil
}

// Usage returns the usage message
func Usage() string {
	return UsageMessage
}

var initialStdout = os.Stdout

func disableStdout() {
	_, w, _ := os.Pipe()
	os.Stdout = w
}

func restoreStdout() {
	os.Stdout = initialStdout
}
