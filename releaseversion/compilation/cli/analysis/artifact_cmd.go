package analysis

import (
	"github.com/dpb587/boshua/analysis/cli/clicommon"
	"github.com/dpb587/boshua/config/provider/setter"
	"github.com/sirupsen/logrus"
)

type ArtifactCmd struct {
	setter.AppConfig `no-flag:"true"`
	*CmdOpts         `no-flag:"true"`

	clicommon.ArtifactCmd
}

func (c *ArtifactCmd) Execute(_ []string) error {
	c.AppConfig.AppendLoggerFields(logrus.Fields{"cli.command": "release/compilation/analysis/artifact"})

	return c.ArtifactCmd.ExecuteAnalysis(c.Config.GetDownloader, c.CmdOpts.getAnalysis)
}
