package stemcell

import (
	"errors"
)

type MetalinkCmd struct {
	*CmdOpts `no-flag:"true"`

	Format string `long:"format" description:"Output format for metalink"`
}

func (c *MetalinkCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("stemcell/metalink")

	return errors.New("TODO")
}
