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
	c.AppOpts.ConfigureLogger("release/compilation/artifact")

	return c.ArtifactCmd.ExecuteArtifact(func() (artifact.Artifact, error) {
		artifact, err := c.CompiledReleaseOpts.Artifact()
		if err != nil {
			return nil, errors.Wrap(err, "finding compiled release")
		}

		return artifact, nil
	})
}
