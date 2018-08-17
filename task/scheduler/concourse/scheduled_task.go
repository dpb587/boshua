package concourse

import (
	"strings"

	"github.com/dpb587/boshua/task/scheduler"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type scheduledTask struct {
	fly          *Fly
	pipelineName string
	logger       logrus.FieldLogger
	subject      interface{}
}

var _ scheduler.ScheduledTask = &scheduledTask{}

func newScheduledTask(fly *Fly, pipelineName string, subject interface{}, logger logrus.FieldLogger) *scheduledTask {
	return &scheduledTask{
		fly:          fly,
		pipelineName: pipelineName,
		logger:       logger,
		subject:      subject,
	}
}

func (t *scheduledTask) Subject() interface{} {
	return t.subject
}

func (t *scheduledTask) Status() (scheduler.Status, error) {
	stdout, stderr, err := t.fly.Run("jobs", "--pipeline", t.pipelineName)
	if err != nil {
		if strings.Contains(string(stderr), "error: resource not found") {
			return scheduler.StatusUnknown, nil
		}

		return scheduler.StatusUnknown, errors.Wrap(err, "listing jobs")
	}

	t.logger.Debugf("checked status of pipeline %s", t.pipelineName)

	lines := strings.Split(string(stdout), "\n")
	if len(lines) < 1 {
		return scheduler.StatusUnknown, errors.New("listing jobs: lines missing")
	}

	fields := strings.Fields(string(stdout))
	if len(fields) != 4 {
		return scheduler.StatusUnknown, errors.New("listing jobs: columns incorrect")
	}

	if fields[3] == "started" {
		return scheduler.StatusRunning, nil
	} else if fields[2] == "succeeded" {
		return scheduler.StatusSucceeded, nil
	} else if fields[2] == "aborted" {
		return scheduler.StatusFailed, nil
	} else if fields[2] == "failed" {
		return scheduler.StatusFailed, nil
	} else if fields[2] == "errored" {
		return scheduler.StatusFailed, nil
	} else if fields[3] == "pending" {
		return scheduler.StatusPending, nil
	} else if fields[2] == "n/a" && fields[3] == "n/a" {
		return scheduler.StatusPending, nil
	}

	return scheduler.StatusUnknown, errors.New("unrecognized pipeline state")
}
