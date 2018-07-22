package cli

import (
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/artifact/cli/clicommon"
	"github.com/pkg/errors"
)

type ArtifactCmd struct {
	clicommon.ArtifactCmd

	*CmdOpts `no-flag:"true"`
}

func (c *ArtifactCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("release/artifact")

	return c.ArtifactCmd.ExecuteArtifact(func() (artifact.Artifact, error) {
		artifact, err := c.CmdOpts.ReleaseOpts.Artifact()
		if err != nil {
			return nil, errors.Wrap(err, "finding release")
		}

		return artifact, nil
	})
}
