package release

import (
	"errors"
)

type MetalinkCmd struct {
	*CmdOpts `no-flag:"true"`

	Format string `long:"format" description:"Output format for metalink"`
}

func (c *MetalinkCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("release/metalink")

	return errors.New("TODO")
}
