package analysis

import (
	"github.com/dpb587/boshua/cli/client/cmd/analysisutil"
)

type MetalinkCmd struct {
	analysisutil.MetalinkCmd

	*CmdOpts `no-flag:"true"`
}

func (c *MetalinkCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("release/analysis/metalink")

	return c.MetalinkCmd.ExecuteAnalysis(c.CmdOpts.getAnalysis)
}
