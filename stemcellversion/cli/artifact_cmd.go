package cli

import (
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/artifact/cli/clicommon"
	"github.com/dpb587/boshua/config/provider/setter"
	"github.com/sirupsen/logrus"
)

type ArtifactCmd struct {
	setter.AppConfig `no-flag:"true"`
	*CmdOpts         `no-flag:"true"`

	clicommon.ArtifactCmd
}

func (c *ArtifactCmd) Execute(_ []string) error {
	c.Config.AppendLoggerFields(logrus.Fields{"cli.command": "stemcell/artifact"})

	return c.ArtifactCmd.ExecuteArtifact(func() (artifact.Artifact, error) {
		return c.StemcellOpts.Artifact(c.AppConfig.Config)
	})
}
