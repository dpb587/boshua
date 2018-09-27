package cli

import (
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/artifact/cli/clicommon"
	"github.com/dpb587/boshua/config/provider/setter"
	"github.com/sirupsen/logrus"
)

type UploadReleaseCmd struct {
	setter.AppConfig `no-flag:"true"`
	*CmdOpts         `no-flag:"true"`

	clicommon.UploadReleaseCmd
}

func (c *UploadReleaseCmd) Execute(extra []string) error {
	c.Config.AppendLoggerFields(logrus.Fields{"cli.command": "release/compilation/upload-release"})

	return c.UploadReleaseCmd.ExecuteArtifact(
		c.Config.GetDownloader,
		func() (artifact.Artifact, error) {
			return c.CompiledReleaseOpts.Artifact(c.AppConfig.Config)
		},
		clicommon.UploadReleaseOpts{
			ExtraArgs: extra,
		},
	)
}
