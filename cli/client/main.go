package main

import (
	"os"

	"github.com/dpb587/bosh-compiled-releases/cli/client/cmd"

	flags "github.com/jessevdk/go-flags"
)

func main() {
	c := struct {
		PatchManifest cmd.PatchManifest `command:"patch-manifest" description:"Patch a manifest for compiled releases"`
	}{
		PatchManifest: cmd.PatchManifest{},
	}

	var parser = flags.NewParser(&c, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}
