package cli

import (
	"os"

	"github.com/dpb587/boshua/releaseversion/analyzers/releasemanifests.v1/formatter"
)

type PropertiesCmd struct {
	Jobs []string `long:"job" description:"Show properties for a specific job"`
}

func (c *PropertiesCmd) Execute(_ []string) error {
	f := formatter.Properties{
		Jobs: c.Jobs,
	}
	return f.Format(os.Stdout, os.Stdin)
}
