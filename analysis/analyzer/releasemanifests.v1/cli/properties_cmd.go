package cli

import (
	"os"

	"github.com/dpb587/boshua/analysis/analyzer/releasemanifests.v1/formatter"
)

type PropertiesCmd struct {
	*CmdOpts `no-flag:"true"`

	Jobs []string `long:"job" description:"Show properties for a specific job"`
}

func (c *PropertiesCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("analysis/formatter/releasemanifests.v1/properties")

	f := formatter.Properties{
		Jobs: c.Jobs,
	}
	return f.Format(os.Stdout, os.Stdin)
}
