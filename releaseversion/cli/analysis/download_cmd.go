package analysis

import (
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/artifact/cli/clicommon"
	"github.com/dpb587/boshua/config/provider/setter"
	"github.com/sirupsen/logrus"
)

type DownloadCmd struct {
	setter.AppConfig `no-flag:"true"`
	*CmdOpts         `no-flag:"true"`

	clicommon.DownloadCmd
}

func (c *DownloadCmd) Execute(_ []string) error {
	c.Config.AppendLoggerFields(logrus.Fields{"cli.command": "release/download"})

	return c.DownloadCmd.ExecuteArtifact(func() (artifact.Artifact, error) {
		return c.CmdOpts.getAnalysis()
	})
}
