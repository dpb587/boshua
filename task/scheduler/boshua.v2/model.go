package boshuaV2

import (
	"github.com/dpb587/boshua/task"
)

type mutationScheduleStemcellAnalysis struct {
	ScheduledTask statusResponse `json:"scheduleStemcellAnalysis"`
}

type statusResponse struct {
	Status task.Status `json:"status"`
}
