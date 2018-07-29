package scheduler

import "github.com/dpb587/boshua/task"

type Scheduler interface {
	Schedule(t *task.Task) (Task, error)
}
