package main

import (
	"fmt"
	"io"
	"os"

	"github.com/dcbishop/jkl/globals"
	"github.com/docopt/docopt-go"
	"github.com/spf13/afero"
)

var usageMessage = `%[1]s

Usage:
  %[2]s [<file>...]
  %[2]s -h | --help

Options:
  -h --help     Show this screen.
`

// Option is a command line option.
type Option func(*App) error

// Options a slice of Options
type Options []Option

// OpenFile opens the given file.
func OpenFile(filename string) func(*App) error {
	return func(a *App) error {
		a.Editor().OpenFile(filename)
		return nil
	}
}

// DisplayHelp displays program help and quits.
func DisplayHelp() func(*App) error {
	return func(a *App) error {
		// [TODO]: Store stdout in app for testing. Don't use terminalui when it's just help. - 2015-01-25 10:27pm
		fmt.Fprintln(a.Out, Usage())
		os.Exit(0)
		return nil
	}
}

// SetUI sets the Apps UI.
func SetUI(ui UI) func(*App) error {
	return func(a *App) error {
		a.SetUI(ui)
		return nil
	}
}

// SetOut sets the Apps output stream.
func SetOut(out io.Writer) func(*App) error {
	return func(a *App) error {
		a.SetOut(out)
		return nil
	}
}

// SetErrOut sets the Apps error output stream.
func SetErrOut(eout io.Writer) func(*App) error {
	return func(a *App) error {
		a.SetErrOut(eout)
		return nil
	}
}

// SetFS sets the Apps filesystem handler.
func SetFS(fs afero.Fs) func(*App) error {
	return func(a *App) error {
		a.SetFS(fs)
		return nil
	}
}

// ParseArgs takes arguments and returns a cli.Options. Will return error if parsing failed.
func ParseArgs(args []string) (Options, error) {
	if len(args) < 2 {
		return []Option{}, nil
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
		return []Option{}, err
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
		return []Option{DisplayHelp()}
	}

	files := arguments["<file>"].([]string)
	for _, f := range files {
		options = append(options, OpenFile(f))
	}

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
