package localexec

import (
	"github.com/dpb587/boshua/task"
	"github.com/dpb587/boshua/task/scheduler"
)

type Scheduler struct {
	cmdFactory CmdFactory
}

var _ scheduler.Scheduler = &Scheduler{}

func NewScheduler(cmdFactory CmdFactory) scheduler.Scheduler {
	return Scheduler{
		cmdFactory: cmdFactory,
	}
}

func (s Scheduler) Schedule(tt task.Task) (scheduler.Task, error) {
	return NewTask(s.cmdFactory, tt), nil
}
