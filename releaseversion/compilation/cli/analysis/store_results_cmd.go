package analysis

import (
	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/cli/clicommon"
)

type StoreResultsCmd struct {
	clicommon.StoreResultsCmd

	*CmdOpts `no-flag:"true"`
}

func (c *StoreResultsCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("release/compilation/analysis/store-results")

	return c.StoreResultsCmd.ExecuteStore(
		c.AppOpts.GetAnalysisIndex,
		func() (analysis.Subject, error) {
			return c.CompiledReleaseOpts.Artifact()
		},
		c.Analyzer,
	)
}
