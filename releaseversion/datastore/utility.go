package datastore

import (
	"fmt"

	"github.com/dpb587/boshua/analysis"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	analysistask "github.com/dpb587/boshua/analysis/task"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/task"
	"github.com/dpb587/boshua/task/scheduler"
	"github.com/pkg/errors"
)

func FilterForOne(index Index, ref releaseversion.Reference) (releaseversion.Artifact, error) {
	results, err := index.Filter(ref)
	if err != nil {
		return releaseversion.Artifact{}, err
	} else if len(results) == 0 {
		return releaseversion.Artifact{}, NoMatchErr
	} else if len(results) > 1 {
		return releaseversion.Artifact{}, MultipleMatchErr
	}

	return results[0], nil
}

func FindOrCreateAnalysis(index Index, scheduler_ scheduler.Scheduler, ref releaseversion.Reference, analyzer analysis.AnalyzerName) (releaseversion.Artifact, analysis.Artifact, error) {
	subject, err := index.Find(ref)
	if err != nil {
		return releaseversion.Artifact{}, analysis.Artifact{}, errors.Wrap(err, "finding release")
	}

	analysisIndex := index.GetAnalysisDatastore()

	analysisRef := analysis.Reference{
		Subject:  subject,
		Analyzer: analyzer,
	}

	analysisSubject, err := analysisdatastore.FilterForOne(analysisIndex, analysisRef)
	if err == analysisdatastore.NoMatchErr {
		tt, err := analysistask.New(subject, analyzer)
		if err != nil {
			return releaseversion.Artifact{}, analysis.Artifact{}, errors.Wrap(err, "preparing task")
		}

		scheduledTask, err := scheduler_.Schedule(tt)
		if err != nil {
			return releaseversion.Artifact{}, analysis.Artifact{}, errors.Wrap(err, "scheduling task")
		}

		status, err := scheduledTask.Wait(nil) // TODO status reporter
		if err != nil {
			return releaseversion.Artifact{}, analysis.Artifact{}, errors.Wrap(err, "checking task")
		} else if status != task.StatusSucceeded {
			return releaseversion.Artifact{}, analysis.Artifact{}, fmt.Errorf("task did not succeed: %s", status)
		}

		analysisSubject, err = analysisdatastore.FilterForOne(analysisIndex, analysisRef)
	}

	if err != nil {
		return releaseversion.Artifact{}, analysis.Artifact{}, errors.Wrap(err, "finding analysis")
	}

	return subject, analysisSubject, nil
}
