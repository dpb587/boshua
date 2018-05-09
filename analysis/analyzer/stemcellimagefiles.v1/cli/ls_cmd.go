package cli

import (
	"os"

	"github.com/dpb587/boshua/analysis/analyzer/stemcellimagefiles.v1/formatter"
)

type LsCmd struct {
	*CmdOpts `no-flag:"true"`
}

func (c *LsCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("analysis/stemcellimagefiles.v1/ls")

	f := formatter.Ls{}
	return f.Format(os.Stdout, os.Stdin)
}
