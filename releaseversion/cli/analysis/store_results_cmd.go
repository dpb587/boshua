package analysis

import (
	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/cli/clicommon"
	"github.com/dpb587/boshua/config/provider/setter"
	"github.com/sirupsen/logrus"
)

type StoreResultsCmd struct {
	setter.AppConfig `no-flag:"true"`
	*CmdOpts         `no-flag:"true"`

	clicommon.StoreResultsCmd
}

func (c *StoreResultsCmd) Execute(_ []string) error {
	c.AppConfig.AppendLoggerFields(logrus.Fields{"cli.command": "release/analysis/store-results"})

	return c.StoreResultsCmd.ExecuteStore(
		c.Config.GetAnalysisIndex,
		func() (analysis.Subject, error) {
			return c.ReleaseOpts.Artifact(c.AppConfig.Config)
		},
		c.Analyzer,
	)
}
