package cli

import (
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/artifact/cli/clicommon"
	"github.com/pkg/errors"
)

type UploadReleaseCmd struct {
	clicommon.UploadReleaseCmd

	*CmdOpts `no-flag:"true"`
}

func (c *UploadReleaseCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("release/upload-release")

	return c.UploadReleaseCmd.ExecuteArtifact(func() (artifact.Artifact, error) {
		artifact, err := c.CmdOpts.ReleaseOpts.Artifact()
		if err != nil {
			return nil, errors.Wrap(err, "finding release")
		}

		return artifact, nil
	})
}
