package cli

import (
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/artifact/cli/clicommon"
	"github.com/pkg/errors"
)

type DownloadCmd struct {
	clicommon.DownloadCmd

	*CmdOpts `no-flag:"true"`
}

func (c *DownloadCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("release/download")

	return c.DownloadCmd.ExecuteArtifact(func() (artifact.Artifact, error) {
		artifact, err := c.CmdOpts.ReleaseOpts.Artifact()
		if err != nil {
			return nil, errors.Wrap(err, "finding release")
		}

		return artifact, nil
	})
}
