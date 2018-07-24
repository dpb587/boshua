package cliutil

import (
	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/cli/clicommon/opts"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/task/scheduler"
	"github.com/pkg/errors"
)

func LoadAnalysis(
	analysisIndexLoader func(analysis.Reference) (analysisdatastore.Index, error),
	subjectLoader func() (analysis.Subject, error),
	analysisOpts *opts.Opts,
	schedulerLoader func() (scheduler.Scheduler, error),
	contextArgs []string,
) (analysis.Artifact, error) {
	subject, err := subjectLoader()
	if err != nil {
		return analysis.Artifact{}, errors.Wrap(err, "loading release")
	}

	analysisRef := analysis.Reference{
		Subject:  subject,
		Analyzer: analysisOpts.Analyzer,
	}

	analysisIndex, err := analysisIndexLoader(analysisRef)
	if err != nil {
		return analysis.Artifact{}, errors.Wrap(err, "loading analysis index")
	}

	results, err := analysisIndex.Filter(analysisRef)
	if err != nil {
		return analysis.Artifact{}, errors.Wrap(err, "finding analysis")
	}

	if len(results) == 0 {
		if analysisOpts.NoWait {
			return analysis.Artifact{}, errors.New("no analysis found")
		}

		scheduler, err := schedulerLoader()
		if err != nil {
			return analysis.Artifact{}, errors.Wrap(err, "loading scheduler")
		}

		err = analysisdatastore.CreateAnalysis(scheduler, analysisRef, contextArgs)
		if err != nil {
			return analysis.Artifact{}, errors.Wrap(err, "creating analysis")
		}

		results, err = analysisIndex.Filter(analysisRef)
		if err != nil {
			return analysis.Artifact{}, errors.Wrap(err, "finding finished analysis")
		}
	}

	result, err := analysisdatastore.RequireSingleResult(results)
	if err != nil {
		return analysis.Artifact{}, errors.Wrap(err, "finding analysis")
	}

	return result, nil
}