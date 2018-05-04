package release

import (
	"fmt"

	"github.com/dpb587/boshua/checksum"
	"github.com/dpb587/boshua/cli/client/cmd/artifactutil"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/metalink"
)

type ArtifactCmd struct {
	artifactutil.ArtifactCmd

	*CmdOpts `no-flag:"true"`
}

func (c *ArtifactCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("release/artifact")

	return c.ArtifactCmd.ExecuteArtifact(func() (metalink.File, error) {
		client := c.AppOpts.GetClient()

		res, err := client.GetReleaseVersion(releaseversion.Reference{
			Name:      c.ReleaseOpts.Release.Name,
			Version:   c.ReleaseOpts.Release.Version,
			Checksums: checksum.ImmutableChecksums{c.ReleaseOpts.ReleaseChecksum.ImmutableChecksum},
		})
		if err != nil {
			return metalink.File{}, fmt.Errorf("fetching: %v", err)
		}

		return res.Data.Artifact, nil
	})
}
