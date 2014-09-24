package cli

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/dcbishop/gim/globals"
	"github.com/docopt/docopt-go"
)

var usageMessage = `%[1]s

Usage:
  %[2]s [<file>...]
  %[2]s -h | --help

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

	// Docopt.go doesn't seem to have a way to stop it spamming the console.
	disableStdout()
	defer restoreStdout()

	version := globals.Name() + " " + globals.VersionString()

	arguments, err := docopt.Parse(Usage(), args[1:], false, version, false, false)
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
	return fmt.Sprintf(usageMessage, globals.Name(), globals.Executable())
}

var initialStdout = os.Stdout

func disableStdout() {
	_, w, _ := os.Pipe()
	os.Stdout = w
}

func restoreStdout() {
	os.Stdout = initialStdout
}
