package schedulerutil

import (
	"fmt"
	"path/filepath"

	"github.com/dpb587/boshua/releaseversion"
	compilationtask "github.com/dpb587/boshua/releaseversion/compilation/task"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/task"
	schedulerpkg "github.com/dpb587/boshua/task/scheduler"
	"github.com/pkg/errors"
)

func CreateCompilation(scheduler schedulerpkg.Scheduler, release releaseversion.Artifact, stemcell stemcellversion.Artifact) error {
	tt, err := compilationtask.New(release, stemcell)
	if err != nil {
		return errors.Wrap(err, "preparing task")
	}

	tt.Steps = append(tt.Steps, task.Step{
		Name: "storing",
		Args: []string{
			"release",
			fmt.Sprintf("--release-name=%s", release.Name),
			fmt.Sprintf("--release-version=%s", release.Version),
			"compilation",
			fmt.Sprintf("--os=%s/%s", stemcell.OS, stemcell.Version),
			"datastore",
			"store",
			filepath.Join("input", fmt.Sprintf("%s-%s-on-%s-stemcell-%s.tgz", release.Name, release.Version, stemcell.OS, stemcell.Version)),
		},
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
