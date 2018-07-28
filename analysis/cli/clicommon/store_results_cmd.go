package clicommon

import (
	"fmt"
	"path/filepath"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/pkg/errors"
)

type StoreResultsCmd struct {
	Analyzer analysis.AnalyzerName `long:"analyzer" description:"The analyzer used for the results"`

	Args StoreResultsArgs `positional-args:"true"`
}

type StoreResultsArgs struct {
	Artifact string `positional-arg-name:"ARTIFACT-PATH" description:"Artifact to store"`
}

func (c *StoreResultsCmd) ExecuteStore(
	analysisIndexLoader func(analysis.Reference) (datastore.Index, error),
	subjectLoader func() (analysis.Subject, error),
	analyzer analysis.AnalyzerName,
) error {
	subject, err := subjectLoader()
	if err != nil {
		return errors.Wrap(err, "finding release")
	}

	ref := analysis.Reference{
		Subject:  subject,
		Analyzer: analyzer,
	}

	index, err := analysisIndexLoader(ref)
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

	err = index.StoreAnalysisResult(ref, *meta4)
	if err != nil {
		return errors.Wrap(err, "storing artifact")
	}

	return nil
}
