package cli

import (
	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/cli/clicommon"
)

type AnalyzersCmd struct {
	clicommon.AnalyzersCmd

	*CmdOpts `no-flag:"true"`
}

func (c *AnalyzersCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("compiledrelease/analyzers")

	return c.AnalyzersCmd.Execute(func() (analysis.Subject, error) {
		return c.CompiledReleaseOpts.Artifact()
	})
}
