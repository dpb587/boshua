package analyzer

import (
	"os"

	"github.com/dpb587/boshua/analysis/analyzer/stemcellmanifest.v1/formatter"
)

type ContentsCmd struct {
	*CmdOpts `no-flag:"true"`
}

func (c *ContentsCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("analysis/stemcellmanifest.v1/contents")

	f := formatter.Contents{}
	return f.Format(os.Stdout, os.Stdin)
}
