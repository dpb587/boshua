package cli

import (
	"os"

	"github.com/dpb587/boshua/releaseversion/analyzers/releaseartifactfiles.v1/formatter"
)

type LsCmd struct {
	*CmdOpts `no-flag:"true"`
}

func (c *LsCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("analysis/formatter/releaseartifactfiles.v1/ls")

	f := formatter.Ls{}
	return f.Format(os.Stdout, os.Stdin)
}
