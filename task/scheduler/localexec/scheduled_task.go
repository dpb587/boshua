package localexec

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dpb587/boshua/task"
	schedulerpkg "github.com/dpb587/boshua/task/scheduler"
	"github.com/pkg/errors"
)

type scheduledTask struct {
	cmdFactory cmdFactory
	tt         *task.Task

	subject interface{}
	status  *schedulerpkg.Status
}

var _ schedulerpkg.ScheduledTask = &scheduledTask{}

func newScheduledTask(cmdFactory cmdFactory, tt *task.Task, subject interface{}) schedulerpkg.ScheduledTask {
	return &scheduledTask{
		cmdFactory: cmdFactory,
		tt:         tt,
		subject:    subject,
	}
}

func (t *scheduledTask) Subject() interface{} {
	return t.subject
}

// TODO synchronous
func (t *scheduledTask) Status() (schedulerpkg.Status, error) {
	var err error

	if t.status == nil {
		var status schedulerpkg.Status

		status, err = t.run()
		t.status = &status
	}

	return *t.status, err
}

func (t *scheduledTask) run() (schedulerpkg.Status, error) {
	inputDir, err := ioutil.TempDir("", "boshua-localexec-input")
	if err != nil {
		return schedulerpkg.StatusFailed, errors.Wrap(err, "creating input directory")
	}

	defer os.RemoveAll(inputDir)

	outputDir, err := ioutil.TempDir("", "boshua-localexec-output")
	if err != nil {
		return schedulerpkg.StatusFailed, errors.Wrap(err, "creating output directory")
	}

	defer os.RemoveAll(outputDir)

	for stepIdx, step := range t.tt.Steps {
		tmpdir, err := ioutil.TempDir("", "boshua-localexec-")
		if err != nil {
			return schedulerpkg.StatusFailed, errors.Wrap(err, "creating tmpdir")
		}

		defer os.RemoveAll(tmpdir)

		if len(step.Input) > 0 {
			for fileName, fileBytes := range step.Input {
				err = ioutil.WriteFile(filepath.Join(inputDir, fileName), fileBytes, 0600)
				if err != nil {
					return schedulerpkg.StatusFailed, errors.Wrapf(err, "writing %s", fileName)
				}
			}
		}

		err = os.Symlink(inputDir, filepath.Join(tmpdir, "input"))
		if err != nil {
			return schedulerpkg.StatusFailed, errors.Wrapf(err, "linking input")
		}

		err = os.Symlink(outputDir, filepath.Join(tmpdir, "output"))
		if err != nil {
			return schedulerpkg.StatusFailed, errors.Wrapf(err, "linking output")
		}

		cmd := t.cmdFactory(step.Args...)
		cmd.Dir = tmpdir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Fprintf(os.Stderr, "%s\n", step.Args)

		err = cmd.Run()
		if err != nil {
			return schedulerpkg.StatusFailed, errors.Wrapf(err, "running step %d", stepIdx)
		}

		{ // outputs -> inputs
			oldInputs, err := filepath.Glob(filepath.Join(inputDir, "*"))
			if err != nil {
				return schedulerpkg.StatusFailed, errors.Wrapf(err, "globbing old inputs")
			}

			for _, path := range oldInputs {
				err = os.RemoveAll(path)
				if err != nil {
					return schedulerpkg.StatusFailed, errors.Wrapf(err, "removing old input")
				}
			}

			newInputs, err := filepath.Glob(filepath.Join(outputDir, "*"))
			if err != nil {
				return schedulerpkg.StatusFailed, errors.Wrapf(err, "globbing new inputs")
			}

			for _, path := range newInputs {
				err = os.Rename(path, filepath.Join(inputDir, strings.TrimPrefix(strings.TrimPrefix(path, outputDir), "/")))
				if err != nil {
					return schedulerpkg.StatusFailed, errors.Wrapf(err, "moving output to next input")
				}
			}
		}
	}

	return schedulerpkg.StatusSucceeded, nil
}
