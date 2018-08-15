package cmd

import (
	"fmt"
	"time"

	"github.com/dpb587/boshua/cli"
	flags "github.com/jessevdk/go-flags"
)

type VersionCmd struct {
	Name   bool `long:"name" description:"Show only the application name"`
	Semver bool `long:"semver" description:"Show only the semver version value"`
	Commit bool `long:"commit" description:"Show only the versioning commit reference"`
	Built  bool `long:"built" description:"Show only the build date"`

	App cli.App
}

var _ flags.Commander = VersionCmd{}

func (c VersionCmd) Execute(_ []string) error {
	if c.Name {
		fmt.Println(c.App.Name)
	} else if c.Semver {
		fmt.Println(c.App.Semver)
	} else if c.Commit {
		fmt.Println(c.App.Commit)
	} else if c.Built {
		fmt.Println(c.App.Built.Format(time.RFC3339))
	} else {
		fmt.Println(c.App.String())
	}

	return nil
}
