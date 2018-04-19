package main

import (
	"os"

	"github.com/dpb587/boshua/cli/client/cmd"

	flags "github.com/jessevdk/go-flags"
)

func main() {
	c := cmd.New()

	var parser = flags.NewParser(c, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}
