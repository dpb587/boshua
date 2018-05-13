package scheduler

import (
	"github.com/dpb587/boshua/task"
)

type Task interface {
	Status() (task.Status, error)
	Wait(func(task.Status)) (task.Status, error)
}
