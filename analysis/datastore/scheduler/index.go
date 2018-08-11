package scheduler

import (
	"fmt"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/task/scheduler"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
)

type index struct {
	index             datastore.Index
	scheduler         scheduler.Scheduler
	schedulerCallback scheduler.StatusChangeCallback
}

var _ datastore.Index = &index{}

func New(idx datastore.Index, scheduler scheduler.Scheduler, schedulerCallback scheduler.StatusChangeCallback) datastore.Index {
	return &index{
		index:             idx,
		scheduler:         scheduler,
		schedulerCallback: schedulerCallback,
	}
}

func (i *index) GetAnalysisArtifacts(ref analysis.Reference) ([]analysis.Artifact, error) {
	results, err := i.index.GetAnalysisArtifacts(ref)
	if err != nil {
		return nil, err
	}

	if len(results) > 0 {
		return results, nil
	}

	scheduledTask, err := i.scheduler.ScheduleAnalysis(ref)
	if err != nil {
		return nil, errors.Wrap(err, "creating analysis")
	}

	status, err := scheduler.WaitForScheduledTask(scheduledTask, i.schedulerCallback)
	if err != nil {
		return nil, errors.Wrap(err, "checking task")
	} else if status != scheduler.StatusSucceeded {
		return nil, fmt.Errorf("task did not succeed: %s", status)
	}

	err = i.index.FlushAnalysisCache()
	if err != nil {
		return nil, errors.Wrap(err, "flushing cache")
	}

	return i.index.GetAnalysisArtifacts(ref)
}

func (i *index) StoreAnalysisResult(ref analysis.Reference, meta4 metalink.Metalink) error {
	return i.index.StoreAnalysisResult(ref, meta4)
}

func (i *index) FlushAnalysisCache() error {
	return i.index.FlushAnalysisCache()
}
