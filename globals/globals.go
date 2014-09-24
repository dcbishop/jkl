package globals

import (
	"os"
	"path/filepath"
)

// Name returns the name of the app
func Name() string {
	return "Gim"
}

// Executable returns the name of the executable that invoked the app
func Executable() string {
	return filepath.Base(os.Args[0])
}

// VersionString returns the version of the app as a string.
func VersionString() string {
	return "0.1"
}
