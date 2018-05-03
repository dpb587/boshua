package compiledreleaseversion

import (
	"github.com/dpb587/boshua/api/v2/models/scheduler"
	"github.com/dpb587/metalink"
)

type GETCompilationInfoResponse struct {
	Data GETCompilationInfoResponseData `json:"data"`
}

type GETCompilationInfoResponseData struct {
	Artifact metalink.File `json:"artifact"`
}

type POSTCompilationQueueResponse struct {
	Data scheduler.TaskStatus `json:"data"`
}
