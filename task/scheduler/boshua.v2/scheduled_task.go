package boshuaV2

import (
	schedulerpkg "github.com/dpb587/boshua/task/scheduler"
)

type scheduledTask struct {
	subject interface{}
	status  func() (schedulerpkg.Status, error)
}

var _ schedulerpkg.ScheduledTask = &scheduledTask{}

func newScheduledTask(status func() (schedulerpkg.Status, error), subject interface{}) *scheduledTask {
	return &scheduledTask{
		status:  status,
		subject: subject,
	}
}

func (t *scheduledTask) Subject() interface{} {
	return t.subject
}

func (t *scheduledTask) Status() (schedulerpkg.Status, error) {
	return t.status()
}
