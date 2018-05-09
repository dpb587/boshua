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
		datastore, err := c.AppOpts.GetReleaseIndex("default")
		if err != nil {
			return metalink.File{}, fmt.Errorf("loading datastore: %v", err)
		}

		res, err := datastore.Find(c.ReleaseOpts.Reference())
		if err != nil {
			return metalink.File{}, fmt.Errorf("fetching: %v", err)
		}

		return res.ArtifactMetalinkFile(), nil
	})
}
