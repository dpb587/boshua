package analysis

import (
	"github.com/dpb587/boshua/cli/client/cmd/analysisutil"
)

type ResultsCmd struct {
	analysisutil.ResultsCmd

	*CmdOpts `no-flag:"true"`
}

func (c *ResultsCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("stemcell/analysis/results")

	return c.ResultsCmd.ExecuteAnalysis(c.CmdOpts.getAnalysis)
}
