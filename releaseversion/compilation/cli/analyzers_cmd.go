package cli

import (
	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/cli/clicommon"
	"github.com/dpb587/boshua/config/provider/setter"
	"github.com/sirupsen/logrus"
)

type AnalyzersCmd struct {
	setter.AppConfig `no-flag:"true"`
	*CmdOpts         `no-flag:"true"`

	clicommon.AnalyzersCmd
}

func (c *AnalyzersCmd) Execute(_ []string) error {
	c.Config.AppendLoggerFields(logrus.Fields{"cli.command": "release/compilation/analyzers"})

	return c.AnalyzersCmd.Execute(func() (analysis.Subject, error) {
		return c.CompiledReleaseOpts.Artifact(c.AppConfig.Config)
	})
}
