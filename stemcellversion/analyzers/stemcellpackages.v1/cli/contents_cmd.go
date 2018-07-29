package cli

import (
	"os"

	"github.com/dpb587/boshua/stemcellversion/analyzers/stemcellpackages.v1/formatter"
)

type ContentsCmd struct {
	*CmdOpts `no-flag:"true"`
}

func (c *ContentsCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("analysis/formatter/stemcellpackages.v1/contents")

	f := formatter.Contents{}
	return f.Format(os.Stdout, os.Stdin)
}
