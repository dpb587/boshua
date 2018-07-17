package analysis

import (
	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/cli/clicommon"
	"github.com/pkg/errors"
)

type StoreResultsCmd struct {
	clicommon.StoreResultsCmd

	*CmdOpts `no-flag:"true"`
}

func (c *StoreResultsCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("release/analysis/store-results")

	subject, err := c.ReleaseOpts.Artifact()
	if err != nil {
		return errors.Wrap(err, "finding release")
	}

	ref := analysis.Reference{
		Subject:  subject,
		Analyzer: analysis.AnalyzerName(c.Analyzer),
	}

	index, err := c.AppOpts.GetAnalysisIndex(ref)
	if err != nil {
		return errors.Wrap(err, "loading analysis datastore")
	}

	return c.StoreResultsCmd.ExecuteStore(index, ref)
}
