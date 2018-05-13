package localexec

import (
	"github.com/dpb587/boshua/task"
	"github.com/dpb587/boshua/task/scheduler"
)

type Scheduler struct{}

func (s Scheduler) Schedule(tt task.Task) (scheduler.Task, error) {
	return NewTask(tt), nil
}
