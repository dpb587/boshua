package cliutil

import (
	"fmt"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/cli/clicommon/opts"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/task"
	schedulerpkg "github.com/dpb587/boshua/task/scheduler"
	"github.com/pkg/errors"
)

func LoadAnalysis(
	analysisIndexLoader func(analysis.Reference) (analysisdatastore.Index, error),
	subjectLoader func() (analysis.Subject, error),
	analysisOpts *opts.Opts,
	schedulerLoader func() (schedulerpkg.Scheduler, error),
	callback task.StatusChangeCallback,
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

	results, err := analysisIndex.GetAnalysisArtifacts(analysisRef)
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

		scheduledTask, err := scheduler.ScheduleAnalysis(analysisRef)
		if err != nil {
			return analysis.Artifact{}, errors.Wrap(err, "creating analysis")
		}

		status, err := task.WaitForScheduledTask(scheduledTask, callback)
		if err != nil {
			return analysis.Artifact{}, errors.Wrap(err, "checking task")
		} else if status != task.StatusSucceeded {
			return analysis.Artifact{}, fmt.Errorf("task did not succeed: %s", status)
		}

		err = analysisIndex.FlushCache()
		if err != nil {
			return analysis.Artifact{}, errors.Wrap(err, "flushing cache")
		}

		results, err = analysisIndex.GetAnalysisArtifacts(analysisRef)
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
