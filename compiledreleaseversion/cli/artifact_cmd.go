package cli

import (
	"fmt"

	"github.com/dpb587/boshua/artifact/cli/clicommon"
	"github.com/dpb587/metalink"
)

type ArtifactCmd struct {
	clicommon.ArtifactCmd

	*CmdOpts `no-flag:"true"`
}

func (c *ArtifactCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("release/artifact")

	return c.ArtifactCmd.ExecuteArtifact(func() (metalink.File, error) {
		artifact, err := c.getCompiledRelease()
		if err != nil {
			return metalink.File{}, fmt.Errorf("finding compiled release: %v", err)
		}

		return artifact.ArtifactMetalinkFile(), nil
	})
}
