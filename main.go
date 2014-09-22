package main

import (
	"fmt"
	"os"

	"github.com/dcbishop/gim/cli"
)

func main() {
	options, err := cli.ParseArgs(os.Args)

	if err != nil {
		fmt.Println(cli.Usage())
		os.Exit(1)
	}

	if options.Help {
		fmt.Println(cli.Usage())
		os.Exit(0)
	}
}
