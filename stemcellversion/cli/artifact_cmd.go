package cli

import (
	"github.com/dpb587/boshua/artifact/cli/clicommon"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
)

type ArtifactCmd struct {
	clicommon.ArtifactCmd

	*CmdOpts `no-flag:"true"`
}

func (c *ArtifactCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("stemcell/artifact")

	return c.ArtifactCmd.ExecuteArtifact(func() (metalink.File, error) {
		artifact, err := c.StemcellOpts.Artifact()
		if err != nil {
			return metalink.File{}, errors.Wrap(err, "finding stemcell")
		}

		return artifact.MetalinkFile(), nil
	})
}
