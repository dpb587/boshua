package cli

import (
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/artifact/cli/clicommon"
	"github.com/dpb587/boshua/config/provider/setter"
	"github.com/sirupsen/logrus"
)

type UploadStemcellCmd struct {
	setter.AppConfig `no-flag:"true"`
	*CmdOpts         `no-flag:"true"`

	clicommon.UploadStemcellCmd
}

func (c *UploadStemcellCmd) Execute(extra []string) error {
	c.Config.AppendLoggerFields(logrus.Fields{"cli.command": "stemcell/upload-stemcell"})

	return c.UploadStemcellCmd.ExecuteArtifact(
		c.Config.GetDownloader,
		func() (artifact.Artifact, error) {
			return c.StemcellOpts.Artifact(c.AppConfig.Config)
		},
		clicommon.UploadStemcellOpts{},
	)
}
