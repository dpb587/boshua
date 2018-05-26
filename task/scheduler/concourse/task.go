package concourse

import (
	"strings"

	"github.com/dpb587/boshua/task"
	"github.com/dpb587/boshua/task/scheduler"
	"github.com/pkg/errors"
)

type Task struct {
	fly          *Fly
	pipelineName string
}

var _ scheduler.Task = &Task{}

func NewTask(fly *Fly, pipelineName string) *Task {
	return &Task{
		fly:          fly,
		pipelineName: pipelineName,
	}
}

func (t *Task) Status() (task.Status, error) {
	stdout, stderr, err := t.fly.Run("jobs", "--pipeline", t.pipelineName)
	if err != nil {
		if strings.Contains(string(stderr), "error: resource not found") {
			return task.StatusUnknown, nil
		}

		return task.StatusUnknown, errors.Wrap(err, "listing jobs")
	}

	lines := strings.Split(string(stdout), "\n")
	if len(lines) < 1 {
		return task.StatusUnknown, errors.New("listing jobs: lines missing")
	}

	fields := strings.Fields(string(stdout))
	if len(fields) != 4 {
		return task.StatusUnknown, errors.New("listing jobs: columns incorrect")
	}

	if fields[2] == "succeeded" {
		return task.StatusSucceeded, nil
	} else if fields[3] == "started" {
		return task.StatusRunning, nil
	} else if fields[2] == "aborted" {
		return task.StatusFailed, nil
	} else if fields[2] == "failed" {
		return task.StatusFailed, nil
	} else if fields[2] == "errored" {
		return task.StatusFailed, nil
	} else if fields[3] == "pending" {
		return task.StatusPending, nil
	} else if fields[2] == "n/a" && fields[3] == "n/a" {
		return task.StatusPending, nil
	}

	return task.StatusUnknown, errors.New("unrecognized pipeline state")
}
