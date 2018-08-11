package scheduler

import (
	"fmt"

	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
	"github.com/dpb587/boshua/task/scheduler"
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

func (i *index) GetCompilationArtifacts(f datastore.FilterParams) ([]compilation.Artifact, error) {
	results, err := i.index.GetCompilationArtifacts(f)
	if err != nil {
		return nil, err
	}

	if len(results) > 0 {
		return results, nil
	}

	scheduledTask, err := i.scheduler.ScheduleCompilation(f)
	if err != nil {
		return nil, errors.Wrap(err, "creating analysis")
	}

	status, err := scheduler.WaitForScheduledTask(scheduledTask, i.schedulerCallback)
	if err != nil {
		return nil, errors.Wrap(err, "checking task")
	} else if status != scheduler.StatusSucceeded {
		return nil, fmt.Errorf("task did not succeed: %s", status)
	}

	err = i.index.FlushCompilationCache()
	if err != nil {
		return nil, errors.Wrap(err, "flushing cache")
	}

	return i.index.GetCompilationArtifacts(f)
}

func (i *index) StoreCompilationArtifact(artifact compilation.Artifact) error {
	return i.index.StoreCompilationArtifact(artifact)
}

func (i *index) FlushCompilationCache() error {
	return i.index.FlushCompilationCache()
}
