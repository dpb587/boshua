package stemcell

import (
	"errors"
)

type ArtifactCmd struct {
	*CmdOpts `no-flag:"true"`

	Format string `long:"format" description:"Output format for the stemcell reference" default:"tsv"`
}

func (c *ArtifactCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("stemcell/metalink")

	return errors.New("TODO")
}
