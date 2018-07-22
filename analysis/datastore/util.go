package datastore

import (
	"fmt"
	"path/filepath"

	"github.com/dpb587/boshua/analysis"
	analysistask "github.com/dpb587/boshua/analysis/task"
	"github.com/dpb587/boshua/task"
	schedulerpkg "github.com/dpb587/boshua/task/scheduler"
	"github.com/pkg/errors"
)

func RequireSingleResult(results []analysis.Artifact) (analysis.Artifact, error) {
	l := len(results)

	if l == 0 {
		return analysis.Artifact{}, errors.New("expected 1 analysis, found 0")
	} else if l > 1 {
		return analysis.Artifact{}, fmt.Errorf("expected 1 analysis, found %d", l)
	}

	return results[0], nil
}

func CreateAnalysis(scheduler schedulerpkg.Scheduler, analysisRef analysis.Reference, storeArgs []string) error {
	tt, err := analysistask.New(analysisRef.Subject, analysisRef.Analyzer)
	if err != nil {
		return errors.Wrap(err, "preparing task")
	}

	tt.Steps = append(tt.Steps, task.Step{
		Name: "storing",
		Args: append(
			storeArgs,
			"analysis",
			"store-results",
			fmt.Sprintf("--analyzer=%s", analysisRef.Analyzer),
			filepath.Join("input", "results.jsonl.gz"),
		),
	})

	scheduledTask, err := scheduler.Schedule(tt)
	if err != nil {
		return errors.Wrap(err, "scheduling task")
	}

	status, err := schedulerpkg.WaitForTask(scheduledTask, nil) // TODO status reporter
	if err != nil {
		return errors.Wrap(err, "checking task")
	} else if status != task.StatusSucceeded {
		return fmt.Errorf("task did not succeed: %s", status)
	}

	return nil
}
