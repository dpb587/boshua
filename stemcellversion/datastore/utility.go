package datastore

import (
	"fmt"
	"path/filepath"

	"github.com/dpb587/boshua/analysis"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	analysistask "github.com/dpb587/boshua/analysis/task"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/task"
	"github.com/dpb587/boshua/task/scheduler"
	"github.com/pkg/errors"
)

func FilterForOne(index Index, ref stemcellversion.Reference) (stemcellversion.Artifact, error) {
	results, err := index.Filter(ref)
	if err != nil {
		return stemcellversion.Artifact{}, err
	} else if len(results) == 0 {
		return stemcellversion.Artifact{}, NoMatchErr
	} else if len(results) > 1 {
		return stemcellversion.Artifact{}, MultipleMatchErr
	}

	return results[0], nil
}

func FindOrCreateAnalysis(index Index, scheduler_ scheduler.Scheduler, ref stemcellversion.Reference, analyzer analysis.AnalyzerName) (stemcellversion.Artifact, analysis.Artifact, error) {
	subject, err := index.Find(ref)
	if err != nil {
		return stemcellversion.Artifact{}, analysis.Artifact{}, errors.Wrap(err, "finding stemcell")
	}

	analysisIndex, err := index.GetAnalysisDatastore(ref)
	if err != nil {
		return stemcellversion.Artifact{}, analysis.Artifact{}, errors.Wrap(err, "getting analysis datastore")
	}

	analysisRef := analysis.Reference{
		Subject:  subject,
		Analyzer: analyzer,
	}

	analysisSubject, analysisSubjectErr := analysisdatastore.FilterForOne(analysisIndex, analysisRef)
	if analysisSubjectErr == analysisdatastore.NoMatchErr {
		tt, err := analysistask.New(subject, analyzer)
		if err != nil {
			return stemcellversion.Artifact{}, analysis.Artifact{}, errors.Wrap(err, "preparing task")
		}

		tt = append(tt, task.Step{
			Name: "storing",
			Args: []string{
				"stemcell",
				fmt.Sprintf("--stemcell=%s/%s", ref.Name(), ref.Version),
				// TODO more options
				"datastore",
				"store-analysis",
				fmt.Sprintf("--analyzer=%s", analyzer),
				filepath.Join("input", "results.jsonl.gz"),
			},
		})

		scheduledTask, err := scheduler_.Schedule(tt)
		if err != nil {
			return stemcellversion.Artifact{}, analysis.Artifact{}, errors.Wrap(err, "scheduling task")
		}

		status, err := scheduler.WaitForTask(scheduledTask, nil) // TODO status reporter
		if err != nil {
			return stemcellversion.Artifact{}, analysis.Artifact{}, errors.Wrap(err, "checking task")
		} else if status != task.StatusSucceeded {
			return stemcellversion.Artifact{}, analysis.Artifact{}, fmt.Errorf("task did not succeed: %s", status)
		}

		analysisSubject, analysisSubjectErr = analysisdatastore.FilterForOne(analysisIndex, analysisRef)
	}

	if analysisSubjectErr != nil {
		return stemcellversion.Artifact{}, analysis.Artifact{}, errors.Wrap(analysisSubjectErr, "finding analysis")
	}

	return subject, analysisSubject, nil
}
