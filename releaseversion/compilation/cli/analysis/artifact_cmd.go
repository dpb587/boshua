package analysis

import (
	"github.com/dpb587/boshua/analysis/cli/clicommon"
)

type ArtifactCmd struct {
	clicommon.ArtifactCmd

	*CmdOpts `no-flag:"true"`
}

func (c *ArtifactCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("release/analysis/artifact")

	return c.ArtifactCmd.ExecuteAnalysis(c.CmdOpts.getAnalysis)
}
