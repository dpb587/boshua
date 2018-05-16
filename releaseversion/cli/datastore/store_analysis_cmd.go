package datastore

import (
	"fmt"
	"path/filepath"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
)

type StoreAnalysisCmd struct {
	*CmdOpts `no-flag:"true"`

	Analyzer string `long:"analyzer" description:"The analyzer used for the results"`

	Args StoreAnalysisArgs `positional-args:"true"`
}

type StoreAnalysisArgs struct {
	Artifact string `positional-arg-name:"ARTIFACT-PATH" description:"Artifact to store"`
}

func (c *StoreAnalysisCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("release/datastore/store-analysis")

	index, err := c.getDatastore()
	if err != nil {
		return errors.Wrap(err, "loading datastore")
	}

	subject, err := index.Find(c.ReleaseOpts.Reference())
	if err != nil {
		return errors.Wrap(err, "filtering")
	}

	analysisIndex := index.GetAnalysisDatastore()

	ref := analysis.Reference{
		Subject:  subject,
		Analyzer: analysis.AnalyzerName(c.Analyzer),
	}

	path, err := filepath.Abs(c.Args.Artifact)
	if err != nil {
		return errors.Wrap(err, "expanding artifact path")
	}

	err = analysisIndex.Store(ref, metalink.Metalink{
		Files: []metalink.File{
			{
				URLs: []metalink.URL{
					{
						URL: fmt.Sprintf("file://%s", path),
					},
				},
			},
		}})
	if err != nil {
		return errors.Wrap(err, "storing artifact")
	}

	return nil
}
