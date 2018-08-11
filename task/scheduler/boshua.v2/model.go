package boshuaV2

import (
	schedulerpkg "github.com/dpb587/boshua/task/scheduler"
)

type mutationSchedule struct {
	Stemcell                   *statusResponse `json:"scheduleStemcellAnalysis"`
	Release                    *statusResponse `json:"scheduleReleaseAnalysis"`
	ReleaseCompilation         *statusResponse `json:"scheduleReleaseCompilation"`
	ReleaseCompilationAnalysis *statusResponse `json:"scheduleReleaseCompilationAnalysis"`
}

type statusResponse struct {
	Status schedulerpkg.Status `json:"status"`
}

func (m mutationSchedule) Status() schedulerpkg.Status {
	if m.Stemcell != nil {
		return m.Stemcell.Status
	} else if m.Release != nil {
		return m.Release.Status
	} else if m.ReleaseCompilation != nil {
		return m.ReleaseCompilation.Status
	} else if m.ReleaseCompilationAnalysis != nil {
		return m.ReleaseCompilationAnalysis.Status
	}

	panic("unexpected")
}
