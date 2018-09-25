package cli

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
	c.Config.AppendLoggerFields(logrus.Fields{"cli.command": "release/compilation/download"})

	return c.DownloadCmd.ExecuteArtifact(c.Config.GetDownloader, func() (artifact.Artifact, error) {
		return c.CompiledReleaseOpts.Artifact(c.AppConfig.Config)
	})
}
