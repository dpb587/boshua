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
	c.AppOpts.ConfigureLogger("release/compilation/download")

	return c.DownloadCmd.ExecuteArtifact(func() (artifact.Artifact, error) {
		artifact, err := c.CompiledReleaseOpts.Artifact()
		if err != nil {
			return nil, errors.Wrap(err, "finding compiled release")
		}

		return artifact, nil
	})
}
