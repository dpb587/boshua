package schedulerutil

import (
	"fmt"
	"path/filepath"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/analyzer/factory"
	"github.com/dpb587/boshua/releaseversion"
	releaseopts "github.com/dpb587/boshua/releaseversion/cli/opts"
	"github.com/dpb587/boshua/releaseversion/compilation"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/stemcellversion"
	stemcellopts "github.com/dpb587/boshua/stemcellversion/cli/opts"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/dpb587/boshua/task"
	schedulerpkg "github.com/dpb587/boshua/task/scheduler"
	"github.com/pkg/errors"
)

func CreateAnalysis(scheduler schedulerpkg.Scheduler, analysisRef analysis.Reference, callback task.StatusChangeCallback) error {
	var storeArgs []string

	switch analysisSubject := analysisRef.Subject.(type) {
	case stemcellversion.Artifact:
		storeArgs = append(
			[]string{"stemcell"},
			stemcellopts.ArgsFromFilterParams(stemcellversiondatastore.FilterParamsFromArtifact(analysisSubject))...,
		)
	case releaseversion.Artifact:
		storeArgs = append(
			[]string{"release"},
			releaseopts.ArgsFromFilterParams(releaseversiondatastore.FilterParamsFromArtifact(analysisSubject))...,
		)
	case compilation.Artifact:
		analysisSubjectRef := analysisSubject.Reference().(compilation.Reference)
		storeArgs = append(
			append(
				append(
					[]string{"release"},
					releaseopts.ArgsFromFilterParams(releaseversiondatastore.FilterParamsFromReference(analysisSubjectRef.ReleaseVersion))...,
				),
				"compilation",
			),
			fmt.Sprintf("--os=%s/%s", analysisSubjectRef.OSVersion.Name, analysisSubjectRef.OSVersion.Version),
		)
	default:
		return errors.New("unsupported analysis subject")
	}

	tt, err := factory.SoonToBeDeprecatedFactory.BuildTask(analysisRef.Analyzer, analysisRef.Subject)
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

	status, err := schedulerpkg.WaitForTask(scheduledTask, callback)
	if err != nil {
		return errors.Wrap(err, "checking task")
	} else if status != task.StatusSucceeded {
		return fmt.Errorf("task did not succeed: %s", status)
	}

	return nil
}
