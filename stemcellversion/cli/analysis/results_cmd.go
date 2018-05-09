package analysis

import (
	"github.com/dpb587/boshua/analysis/cli/clicommon"
)

type ResultsCmd struct {
	clicommon.ResultsCmd

	*CmdOpts `no-flag:"true"`
}

func (c *ResultsCmd) Execute(extra []string) error {
	c.AppOpts.ConfigureLogger("stemcell/analysis/results")

	return c.ResultsCmd.ExecuteAnalysis(c.CmdOpts.AnalysisOpts.Analyzer, c.CmdOpts.getAnalysis, extra)
}
