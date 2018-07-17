package datastore

import (
	"fmt"
	"path/filepath"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/stemcellversion/datastore"
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
	c.AppOpts.ConfigureLogger("stemcell/datastore/store-analysis")

	index, err := c.getDatastore()
	if err != nil {
		return errors.Wrap(err, "loading datastore")
	}

	subjects, err := index.Filter(c.StemcellOpts.FilterParams())
	if err != nil {
		return errors.Wrap(err, "filtering")
	}

	subject, err := datastore.RequireSingleResult(subjects)
	if err != nil {
		return errors.Wrap(err, "finding stemcell")
	}

	ref := analysis.Reference{
		Subject:  subject,
		Analyzer: analysis.AnalyzerName(c.Analyzer),
	}

	analysisIndex, err := c.AppOpts.GetAnalysisIndex(ref)
	if err != nil {
		return errors.Wrap(err, "loading analysis datastore")
	}

	path, err := filepath.Abs(c.Args.Artifact)
	if err != nil {
		return errors.Wrap(err, "expanding artifact path")
	}

	meta4, err := metalinkutil.CreateFromFiles(fmt.Sprintf("file://%s", path))
	if err != nil {
		return errors.Wrap(err, "creating in-memory metalink")
	}

	err = analysisIndex.Store(ref, *meta4)
	if err != nil {
		return errors.Wrap(err, "storing artifact")
	}

	return nil
}
