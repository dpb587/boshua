package localexec

import (
	"github.com/dpb587/boshua/task"
	"github.com/dpb587/boshua/task/scheduler"
	"github.com/sirupsen/logrus"
)

type Scheduler struct {
	config     Config
	cmdFactory CmdFactory
	logger     logrus.FieldLogger
}

var _ scheduler.Scheduler = &Scheduler{}

func New(config Config, cmdFactory CmdFactory, logger logrus.FieldLogger) scheduler.Scheduler {
	return Scheduler{
		config:     config,
		cmdFactory: cmdFactory,
		logger:     logger,
	}
}

func (s Scheduler) Schedule(tt task.Task) (scheduler.Task, error) {
	return NewTask(s.cmdFactory, tt), nil
}
