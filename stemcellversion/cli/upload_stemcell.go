package cli

import (
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/artifact/cli/clicommon"
	"github.com/pkg/errors"
)

type UploadStemcellCmd struct {
	clicommon.UploadStemcellCmd

	*CmdOpts `no-flag:"true"`
}

func (c *UploadStemcellCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("stemcell/artifact")

	return c.UploadStemcellCmd.ExecuteArtifact(func() (artifact.Artifact, error) {
		artifact, err := c.StemcellOpts.Artifact()
		if err != nil {
			return nil, errors.Wrap(err, "finding stemcell")
		}

		return artifact, nil
	})
}
