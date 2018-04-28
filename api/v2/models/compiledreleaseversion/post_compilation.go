package compiledreleaseversion

import "github.com/dpb587/boshua/api/v2/models/scheduler"

type POSTCompilationResponse struct {
	Data scheduler.TaskStatus `json:"data"`
}
