package scheduler

import (
	"time"

	"github.com/dpb587/boshua/task"
)

type Task interface {
	Status() (task.Status, error)
}

func WaitForTask(t Task, callback task.StatusChangeCallback) (task.Status, error) {
	currStatus := task.StatusUnknown

	if callback == nil {
		callback = func(_ task.Status) {}
	}

	for {
		status, err := t.Status()
		if status != currStatus {
			callback(status)

			currStatus = status
		}

		if status == task.StatusSucceeded || status == task.StatusFailed || err != nil {
			return status, err
		}

		time.Sleep(3 * time.Second)
	}
}
