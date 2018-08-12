package cli

import (
	"os"

	"github.com/dpb587/boshua/releaseversion/analyzers/releasemanifests.v1/formatter"
)

type SpecCmd struct {
	*CmdOpts `no-flag:"true"`

	ReleaseOnly bool     `long:"release" description:"Show only release manifest"`
	Jobs        []string `long:"job" description:"Show spec files for a specific job"`
}

func (c *SpecCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("analysis/formatter/releasemanifests.v1/properties")

	f := formatter.Spec{
		ReleaseOnly: c.ReleaseOnly,
		Jobs:        c.Jobs,
	}
	return f.Format(os.Stdout, os.Stdin)
}
