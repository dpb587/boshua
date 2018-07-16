package datastore

import (
	"fmt"
	"path/filepath"

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
	// analysisIndex, err := index.GetAnalysisDatastore()
	// if err != nil {
	// 	return releaseversion.Artifact{}, analysis.Artifact{}, errors.Wrap(err, "getting analysis datastore")
	// }

	analysisRef := analysis.Reference{
		Subject:  subject,
		Analyzer: analyzer,
	}

	analysisSubject, analysisSubjectErr := analysisdatastore.FilterForOne(analysisIndex, analysisRef)
	if analysisSubjectErr == analysisdatastore.NoMatchErr {
		tt, err := analysistask.New(subject, analyzer)
		if err != nil {
			return releaseversion.Artifact{}, analysis.Artifact{}, errors.Wrap(err, "preparing task")
		}

		tt = append(tt, task.Step{
			Name: "storing",
			Args: []string{
				"release",
				fmt.Sprintf("--release=%s/%s", ref.Name, ref.Version),
				// TODO more options
				"datastore",
				"store-analysis",
				fmt.Sprintf("--analyzer=%s", analyzer),
				filepath.Join("input", "results.jsonl.gz"),
			},
		})

		scheduledTask, err := scheduler_.Schedule(tt)
		if err != nil {
			return releaseversion.Artifact{}, analysis.Artifact{}, errors.Wrap(err, "scheduling task")
		}

		status, err := scheduler.WaitForTask(scheduledTask, nil) // TODO status reporter
		if err != nil {
			return releaseversion.Artifact{}, analysis.Artifact{}, errors.Wrap(err, "checking task")
		} else if status != task.StatusSucceeded {
			return releaseversion.Artifact{}, analysis.Artifact{}, fmt.Errorf("task did not succeed: %s", status)
		}

		analysisSubject, analysisSubjectErr = analysisdatastore.FilterForOne(analysisIndex, analysisRef)
	}

	if analysisSubjectErr != nil {
		return releaseversion.Artifact{}, analysis.Artifact{}, errors.Wrap(analysisSubjectErr, "finding analysis")
	}

	return subject, analysisSubject, nil
}
