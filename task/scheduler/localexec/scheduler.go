package localexec

import (
	"os/exec"

	"github.com/dpb587/boshua/task"
	"github.com/dpb587/boshua/task/scheduler"
	"github.com/sirupsen/logrus"
)

type Scheduler struct {
	config Config
	logger logrus.FieldLogger
}

var _ scheduler.Scheduler = &Scheduler{}

func New(config Config, logger logrus.FieldLogger) scheduler.Scheduler {
	return Scheduler{
		config: config,
		logger: logger,
	}
}

func (s Scheduler) cmd(args ...string) *exec.Cmd {
	cmd := exec.Command(s.config.Exec, append(s.config.Args, args...)...)

	return cmd
}

func (s Scheduler) Schedule(tt task.Task) (scheduler.Task, error) {
	return NewTask(s.cmd, tt), nil
}
