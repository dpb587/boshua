package cli

import (
	"os"

	"github.com/dpb587/boshua/stemcellversion/analyzers/stemcellpackages.v1/formatter"
)

type PackagesCmd struct {
	*CmdOpts `no-flag:"true"`
}

func (c *PackagesCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("analysis/formatter/stemcellpackages.v1/packages")

	f := formatter.Packages{}
	return f.Format(os.Stdout, os.Stdin)
}