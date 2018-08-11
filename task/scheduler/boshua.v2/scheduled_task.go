package boshuaV2

import (
	schedulerpkg "github.com/dpb587/boshua/task/scheduler"
)

type scheduledTask struct {
	status func() (schedulerpkg.Status, error)
}

var _ schedulerpkg.ScheduledTask = &scheduledTask{}

func newScheduledTask(status func() (schedulerpkg.Status, error)) *scheduledTask {
	return &scheduledTask{
		status: status,
	}
}

func (t *scheduledTask) Subject() interface{} {
	panic("TODO")
}

func (t *scheduledTask) Status() (schedulerpkg.Status, error) {
	return t.status()
}
