package analysis

import "github.com/dpb587/boshua/api/v2/models/scheduler"

type POSTAnalysisResponse struct {
	Data scheduler.TaskStatus `json:"data"`
}
