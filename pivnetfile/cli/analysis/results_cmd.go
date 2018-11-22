package analysis

import (
	"github.com/dpb587/boshua/analysis/cli/clicommon"
	"github.com/dpb587/boshua/config/provider/setter"
	"github.com/sirupsen/logrus"
)

type ResultsCmd struct {
	setter.AppConfig `no-flag:"true"`
	*CmdOpts         `no-flag:"true"`

	clicommon.ResultsCmd
}

func (c *ResultsCmd) Execute(args []string) error {
	c.AppConfig.AppendLoggerFields(logrus.Fields{"cli.command": "pivnetfile/analysis/results"})

	return c.ResultsCmd.ExecuteAnalysis(
		c.Config.GetDownloader,
		c.CmdOpts.AnalysisOpts.Analyzer,
		c.CmdOpts.getAnalysis,
		args,
	)
}
