package cli

// Options stores options parsed from the command line
type Options struct {
	filesToOpen []string
}

// ParseArgs takes arguments and returns a cli.Options. Will return error if parsing failed.
func ParseArgs(args []string) (Options, error) {
	return Options{}, nil
}
