package boshuaV2

import (
	"github.com/dpb587/boshua/task"
)

type Task struct {
	status func() (task.Status, error)
}

var _ task.ScheduledTask = &Task{}

func NewTask(status func() (task.Status, error)) *Task {
	return &Task{
		status: status,
	}
}

func (t *Task) Status() (task.Status, error) {
	return t.status()
}
