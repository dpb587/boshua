package compiledrelease

import (
	"log"

	"github.com/dpb587/boshua/cli/client/cmd/artifactutil"
	"github.com/dpb587/metalink"
)

type ArtifactCmd struct {
	artifactutil.ArtifactCmd

	*CmdOpts `no-flag:"true"`
}

func (c *ArtifactCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("release/artifact")

	return c.ArtifactCmd.ExecuteArtifact(func() (metalink.File, error) {
		resInfo, err := c.getCompiledRelease()
		if err != nil {
			log.Fatalf("requesting compiled version info: %v", err)
		} else if resInfo == nil {
			log.Fatalf("no compiled release available")
		}

		return resInfo.Data.Artifact, nil
	})
}