package stemcell

import (
	"fmt"

	"github.com/dpb587/boshua/cli/client/cmd/artifactutil"
	"github.com/dpb587/metalink"
)

type ArtifactCmd struct {
	artifactutil.ArtifactCmd

	*CmdOpts `no-flag:"true"`
}

func (c *ArtifactCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("stemcell/artifact")

	return c.ArtifactCmd.ExecuteArtifact(func() (metalink.File, error) {
		client := c.AppOpts.GetClient()

		res, err := client.GetStemcellVersion(c.StemcellOpts.Reference())
		if err != nil {
			return metalink.File{}, fmt.Errorf("fetching: %v", err)
		}

		return res.Data.Artifact, nil
	})
}
