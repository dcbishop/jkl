package main

import (
	"fmt"
	"os"

	"github.com/dcbishop/jkl/globals"
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
	if len(args) < 2 {
		return Options{}, nil
	}

	options, err := parseWithDocopt(args[1:])
	return options, err
}

// parseWithDocopt takes a slice of args (without the program name) and returns an Options
func parseWithDocopt(args []string) (Options, error) {
	// Docopt.go doesn't seem to have a way to stop it spamming the console.
	disableStdout()
	defer restoreStdout()

	version := nameVersion()
	arguments, err := docopt.Parse(Usage(), args, false, version, false, false)
	if err != nil {
		return Options{}, err
	}

	options := docoptArgsToOptions(arguments)

	return options, nil
}

// Returns the name and version (ie "Jkl 0.1")
func nameVersion() string {
	return globals.Name() + " " + globals.VersionString()
}

// Converts the result of docopt.Parse into an Options
func docoptArgsToOptions(arguments map[string]interface{}) Options {
	options := Options{}
	if arguments["--help"].(bool) {
		options.Help = true
		return options
	}

	options.FilesToOpen = arguments["<file>"].([]string)

	return options
}

// Usage returns the usage message
func Usage() string {
	return fmt.Sprintf(usageMessage, globals.Name(), globals.Executable())
}

var initialStdout = os.Stdout

// Replaces the stdout with a dummy pipe to stop console output.
func disableStdout() {
	_, w, _ := os.Pipe()
	os.Stdout = w
}

// Restores the stdout to the writer it has at program start.
func restoreStdout() {
	os.Stdout = initialStdout
}
