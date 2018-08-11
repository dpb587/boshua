package cliutil

import (
	"fmt"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/cli/clicommon/opts"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	schedulerpkg "github.com/dpb587/boshua/task/scheduler"
	"github.com/pkg/errors"
)

func LoadAnalysis(
	analysisIndexLoader func(analysis.Reference) (analysisdatastore.Index, error),
	subjectLoader func() (analysis.Subject, error),
	analysisOpts *opts.Opts,
	schedulerLoader func() (schedulerpkg.Scheduler, error),
	callback schedulerpkg.StatusChangeCallback,
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

	result, err := analysisdatastore.GetAnalysisArtifact(analysisIndex, analysisRef)
	if err == analysisdatastore.NoMatchErr {
		if analysisOpts.NoWait {
			return analysis.Artifact{}, errors.New("no analysis found")
		}

		scheduler, err := schedulerLoader()
		if err != nil {
			return analysis.Artifact{}, errors.Wrap(err, "loading scheduler")
		}

		scheduledTask, err := scheduler.ScheduleAnalysis(analysisRef)
		if err != nil {
			return analysis.Artifact{}, errors.Wrap(err, "creating analysis")
		}

		status, err := schedulerpkg.WaitForScheduledTask(scheduledTask, callback)
		if err != nil {
			return analysis.Artifact{}, errors.Wrap(err, "checking task")
		} else if status != schedulerpkg.StatusSucceeded {
			return analysis.Artifact{}, fmt.Errorf("task did not succeed: %s", status)
		}

		err = analysisIndex.FlushAnalysisCache()
		if err != nil {
			return analysis.Artifact{}, errors.Wrap(err, "flushing cache")
		}

		result, err = analysisdatastore.GetAnalysisArtifact(analysisIndex, analysisRef)
		if err != nil {
			return analysis.Artifact{}, errors.Wrap(err, "finding finished analysis")
		}
	} else if err != nil {
		return analysis.Artifact{}, errors.Wrap(err, "finding analysis")
	}

	return result, nil
}
