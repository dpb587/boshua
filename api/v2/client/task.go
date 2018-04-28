package client

import "github.com/dpb587/boshua/api/v2/models/scheduler"

type TaskStatusWatcher func(scheduler.TaskStatus)
