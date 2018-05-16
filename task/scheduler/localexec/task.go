package localexec

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dpb587/boshua/task"
	"github.com/dpb587/boshua/task/scheduler"
	"github.com/pkg/errors"
)

type Task struct {
	tt     task.Task
	status task.Status
}

var _ scheduler.Task = &Task{}

func NewTask(tt task.Task) *Task {
	return &Task{
		tt:     tt,
		status: task.StatusPending,
	}
}

func (t *Task) Status() (task.Status, error) {
	return t.status, nil
}

// TODO synchronous
func (t *Task) Wait(_ func(task.Status)) (task.Status, error) {
	inputDir, err := ioutil.TempDir("", "boshua-localexec-input")
	if err != nil {
		return task.StatusFailed, errors.Wrap(err, "creating input directory")
	}

	defer os.RemoveAll(inputDir)

	outputDir, err := ioutil.TempDir("", "boshua-localexec-output")
	if err != nil {
		return task.StatusFailed, errors.Wrap(err, "creating output directory")
	}

	defer os.RemoveAll(outputDir)

	for stepIdx, step := range t.tt {
		tmpdir, err := ioutil.TempDir("", "boshua-localexec-")
		if err != nil {
			return task.StatusFailed, errors.Wrap(err, "creating tmpdir")
		}

		defer os.RemoveAll(tmpdir)

		if len(step.Input) > 0 {
			for fileName, fileBytes := range step.Input {
				err = ioutil.WriteFile(filepath.Join(inputDir, fileName), fileBytes, 0600)
				if err != nil {
					return task.StatusFailed, errors.Wrapf(err, "writing %s", fileName)
				}
			}
		}

		err = os.Symlink(inputDir, filepath.Join(tmpdir, "input"))
		if err != nil {
			return task.StatusFailed, errors.Wrapf(err, "linking input")
		}

		err = os.Symlink(outputDir, filepath.Join(tmpdir, "output"))
		if err != nil {
			return task.StatusFailed, errors.Wrapf(err, "linking output")
		}

		cmd := exec.Command("boshua", step.Args...)
		cmd.Dir = tmpdir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// fmt.Printf("%s\n", step.Args)

		err = cmd.Run()
		if err != nil {
			return task.StatusFailed, errors.Wrapf(err, "running step %d", stepIdx)
		}

		{ // outputs -> inputs
			oldInputs, err := filepath.Glob(filepath.Join(inputDir, "*"))
			if err != nil {
				return task.StatusFailed, errors.Wrapf(err, "globbing old inputs")
			}

			for _, path := range oldInputs {
				err = os.RemoveAll(path)
				if err != nil {
					return task.StatusFailed, errors.Wrapf(err, "removing old input")
				}
			}

			newInputs, err := filepath.Glob(filepath.Join(outputDir, "*"))
			if err != nil {
				return task.StatusFailed, errors.Wrapf(err, "globbing new inputs")
			}

			for _, path := range newInputs {
				err = os.Rename(path, filepath.Join(inputDir, strings.TrimPrefix(strings.TrimPrefix(path, outputDir), "/")))
				if err != nil {
					return task.StatusFailed, errors.Wrapf(err, "moving output to next input")
				}
			}
		}
	}

	return task.StatusSucceeded, nil
}
